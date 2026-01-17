package sherpa

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aira-id/gribe/internal/domain"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

// Provider implements the ASRProvider interface using sherpa-onnx
type Provider struct {
	config          *domain.TranscriptionConfig
	recognizer      *sherpa.OnlineRecognizer
	mu              sync.Mutex
	isInitialized   bool
	supportedModels []string
	supportedLangs  []string
}

// New creates a new sherpa-onnx ASR provider
func New(config *domain.TranscriptionConfig) (*Provider, error) {
	if config == nil {
		config = &domain.TranscriptionConfig{
			Model:    "zipformer",
			Language: "id",
		}
	}

	provider := &Provider{
		config:          config,
		supportedModels: []string{"zipformer", "paraformer", "transducer"},
		supportedLangs:  []string{"id", "en", "zh", "de", "es", "fr", "ja", "ko", "ru"},
	}

	// Initialize the recognizer
	if err := provider.initializeRecognizer(); err != nil {
		return nil, fmt.Errorf("failed to initialize sherpa-onnx recognizer: %w", err)
	}

	return provider, nil
}

// initializeRecognizer initializes the sherpa-onnx recognizer
func (p *Provider) initializeRecognizer() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Initializing sherpa-onnx recognizer with model: %s (language: %s)",
		p.config.Model, p.config.Language)

	recognizerConfig := &sherpa.OnlineRecognizerConfig{}
	recognizerConfig.FeatConfig.SampleRate = 16000
	recognizerConfig.FeatConfig.FeatureDim = 80

	// Use the provided Indonesian model files (Transducer model)
	modelDir := "models/sherpa-onnx-streaming-zipformer2-id"
	recognizerConfig.ModelConfig.Transducer.Encoder = modelDir + "/encoder-iter-100000-avg-15-chunk-32-left-256.onnx"
	recognizerConfig.ModelConfig.Transducer.Decoder = modelDir + "/decoder-iter-100000-avg-15-chunk-32-left-256.onnx"
	recognizerConfig.ModelConfig.Transducer.Joiner = modelDir + "/joiner-iter-100000-avg-15-chunk-32-left-256.onnx"
	recognizerConfig.ModelConfig.Tokens = modelDir + "/tokens.txt"

	recognizerConfig.ModelConfig.NumThreads = 4
	recognizerConfig.ModelConfig.Provider = "cpu"
	recognizerConfig.ModelConfig.Debug = 0
	recognizerConfig.DecodingMethod = "greedy_search"
	recognizerConfig.MaxActivePaths = 4

	p.recognizer = sherpa.NewOnlineRecognizer(recognizerConfig)
	if p.recognizer == nil {
		err := fmt.Errorf("sherpa.NewOnlineRecognizer returned nil - check model paths and library compatibility")
		log.Printf("[ERROR] %v", err)
		return err
	}

	p.isInitialized = true
	log.Printf("Sherpa-onnx recognizer initialized successfully with Indonesian zipformer2 model")

	return nil
}

// Transcribe processes audio data and returns transcription results via a channel
func (p *Provider) Transcribe(ctx context.Context, audio []byte, config *domain.TranscriptionConfig) (<-chan domain.TranscriptionChunk, error) {
	resultChan := make(chan domain.TranscriptionChunk, 10)

	if !p.isInitialized {
		close(resultChan)
		return resultChan, fmt.Errorf("recognizer not initialized")
	}

	if len(audio) == 0 {
		close(resultChan)
		return resultChan, fmt.Errorf("audio data is empty")
	}

	go func() {
		defer close(resultChan)

		p.mu.Lock()
		defer p.mu.Unlock()

		stream := sherpa.NewOnlineStream(p.recognizer)
		if stream == nil {
			log.Printf("Error: failed to create OnlineStream")
			return
		}
		defer sherpa.DeleteOnlineStream(stream)

		// Convert bytes to float32 samples
		samples := bytesToFloat32(audio)

		// Add left padding (0.3 seconds of silence)
		leftPadding := make([]float32, 4800) // 16000 * 0.3
		stream.AcceptWaveform(16000, leftPadding)

		// Process the audio
		stream.AcceptWaveform(16000, samples)

		// Add right padding (0.6 seconds of silence)
		rightPadding := make([]float32, 9600) // 16000 * 0.6
		stream.AcceptWaveform(16000, rightPadding)

		// Input finished
		stream.InputFinished()

		// Decode
		for p.recognizer.IsReady(stream) {
			select {
			case <-ctx.Done():
				return
			default:
				p.recognizer.Decode(stream)
			}
		}

		// Get final result
		result := p.recognizer.GetResult(stream)

		// Send final result
		if result != nil && result.Text != "" {
			finalChunk := domain.TranscriptionChunk{
				Text:    result.Text,
				IsFinal: true,
				StartMs: 0,
				EndMs:   len(samples) * 1000 / 16000,
			}
			select {
			case <-ctx.Done():
				return
			case resultChan <- finalChunk:
			}
			log.Printf("Transcription completed: %s", result.Text)
		}
	}()

	return resultChan, nil
}

// TranscribeStream processes audio data in streaming mode
func (p *Provider) TranscribeStream(ctx context.Context, config *domain.TranscriptionConfig) (chan<- []byte, <-chan domain.TranscriptionChunk, error) {
	audioIn := make(chan []byte, 100)
	resultOut := make(chan domain.TranscriptionChunk, 10)

	if !p.isInitialized {
		close(audioIn)
		close(resultOut)
		return audioIn, resultOut, fmt.Errorf("recognizer not initialized")
	}

	go func() {
		defer close(resultOut)

		p.mu.Lock()
		stream := sherpa.NewOnlineStream(p.recognizer)
		p.mu.Unlock()

		if stream == nil {
			log.Printf("Error: failed to create OnlineStream")
			return
		}
		defer sherpa.DeleteOnlineStream(stream)

		var lastPartialResult string

		for {
			select {
			case <-ctx.Done():
				return

			case audio, ok := <-audioIn:
				if !ok {
					// Channel closed, finalize
					stream.InputFinished()

					p.mu.Lock()
					// Finalize decoding
					for p.recognizer.IsReady(stream) {
						p.recognizer.Decode(stream)
					}
					result := p.recognizer.GetResult(stream)
					p.mu.Unlock()

					// Send final result
					if result != nil && result.Text != "" && result.Text != lastPartialResult {
						chunk := domain.TranscriptionChunk{
							Text:    result.Text[len(lastPartialResult):],
							IsFinal: true,
						}
						select {
						case <-ctx.Done():
							return
						case resultOut <- chunk:
						}
					} else {
						// Send empty final chunk if no new text
						chunk := domain.TranscriptionChunk{
							Text:    "",
							IsFinal: true,
						}
						resultOut <- chunk
					}

					return
				}

				// Convert bytes to float32 samples
				samples := bytesToFloat32(audio)

				p.mu.Lock()
				// Accept waveform
				stream.AcceptWaveform(16000, samples)

				// Decode if ready
				for p.recognizer.IsReady(stream) {
					p.recognizer.Decode(stream)
				}

				// Get current result
				result := p.recognizer.GetResult(stream)
				p.mu.Unlock()

				// Send delta event if result changed
				if result != nil && result.Text != "" && result.Text != lastPartialResult {
					delta := result.Text[len(lastPartialResult):]
					if delta != "" {
						chunk := domain.TranscriptionChunk{
							Text:    delta,
							IsFinal: false,
						}
						select {
						case <-ctx.Done():
							return
						case resultOut <- chunk:
						}
						lastPartialResult = result.Text
					}
				}
			}
		}
	}()

	return audioIn, resultOut, nil
}

// GetSupportedModels returns list of supported ASR models
func (p *Provider) GetSupportedModels() []string {
	return p.supportedModels
}

// GetSupportedLanguages returns list of supported language codes
func (p *Provider) GetSupportedLanguages() []string {
	return p.supportedLangs
}

// Close releases any resources held by the provider
func (p *Provider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.recognizer != nil {
		sherpa.DeleteOnlineRecognizer(p.recognizer)
		p.recognizer = nil
	}

	p.isInitialized = false
	log.Printf("Sherpa-onnx provider closed")
	return nil
}

// bytesToFloat32 converts byte array (PCM 16-bit little-endian) to float32 array
func bytesToFloat32(data []byte) []float32 {
	numSamples := len(data) / 2
	samples := make([]float32, numSamples)

	for i := 0; i < numSamples; i++ {
		// Read 16-bit signed integer in little-endian
		b1 := int16(data[i*2])
		b2 := int16(data[i*2+1])
		sample := (b2 << 8) | (b1 & 0xFF)

		// Convert to float32 in range [-1, 1)
		samples[i] = float32(sample) / 32768.0
	}

	return samples
}
