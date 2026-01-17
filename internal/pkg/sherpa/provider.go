package sherpa

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aira-id/gribe/internal/domain"
)

// OnlineRecognizer is a wrapper for sherpa-onnx OnlineRecognizer
// This would be the actual sherpa-onnx type when the library is available
type OnlineRecognizer interface {
	IsReady(stream interface{}) bool
	Decode(stream interface{})
	GetResult(stream interface{}) *Result
}

// OnlineStream is a wrapper for sherpa-onnx OnlineStream
type OnlineStream interface {
	AcceptWaveform(sampleRate int, samples []float32)
}

// Result represents the result from sherpa-onnx
type Result struct {
	Text string
}

// Provider implements the ASRProvider interface using sherpa-onnx
type Provider struct {
	config          *domain.TranscriptionConfig
	recognizer      OnlineRecognizer
	stream          OnlineStream
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
			Language: "en",
		}
	}

	provider := &Provider{
		config:          config,
		supportedModels: []string{"zipformer", "paraformer", "transducer"},
		supportedLangs:  []string{"en", "zh", "de", "es", "fr", "ja", "ko", "ru"},
	}

	// Initialize the recognizer
	if err := provider.initializeRecognizer(); err != nil {
		return nil, fmt.Errorf("failed to initialize sherpa-onnx recognizer: %w", err)
	}

	return provider, nil
}

// initializeRecognizer initializes the sherpa-onnx recognizer with OnlineRecognizer
func (p *Provider) initializeRecognizer() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Initializing sherpa-onnx recognizer with model: %s (language: %s)",
		p.config.Model, p.config.Language)

	// TODO: Uncomment and use when sherpa-onnx Go bindings are available
	// This is a placeholder for the actual initialization code
	/*
		recognizerConfig := &sherpa.OnlineRecognizerConfig{}
		recognizerConfig.FeatConfig = sherpa.FeatureConfig{
			SampleRate: 16000,
			FeatureDim: 80,
		}

		recognizerConfig.ModelConfig.Zipformer2Ctc.Model = p.config.Model
		recognizerConfig.ModelConfig.NumThreads = 4
		recognizerConfig.ModelConfig.Provider = "cpu"
		recognizerConfig.ModelConfig.Debug = 0
		recognizerConfig.DecodingMethod = "greedy_search"
		recognizerConfig.MaxActivePaths = 4

		var err error
		p.recognizer, err = sherpa.NewOnlineRecognizer(recognizerConfig)
		if err != nil {
			return fmt.Errorf("failed to create OnlineRecognizer: %w", err)
		}

		p.stream, err = sherpa.NewOnlineStream(p.recognizer)
		if err != nil {
			return fmt.Errorf("failed to create OnlineStream: %w", err)
		}
	*/

	p.isInitialized = true
	log.Printf("Sherpa-onnx recognizer initialized successfully (mock mode)")

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

		// Convert bytes to float32 samples (assuming PCM 16-bit)
		samples := bytesToFloat32(audio)

		// TODO: Uncomment when sherpa-onnx library is available
		/*
			// Add left padding (0.3 seconds of silence)
			leftPadding := make([]float32, int(float32(16000)*0.3))
			p.stream.AcceptWaveform(16000, leftPadding)

			// Process the audio
			p.stream.AcceptWaveform(16000, samples)

			// Add right padding (0.6 seconds of silence)
			rightPadding := make([]float32, int(float32(16000)*0.6))
			p.stream.AcceptWaveform(16000, rightPadding)

			// Decode
			for p.recognizer.IsReady(p.stream) {
				select {
				case <-ctx.Done():
					return
				default:
					p.recognizer.Decode(p.stream)
				}
			}

			// Get results
			result := p.recognizer.GetResult(p.stream)

			// Send partial result as a chunk
			if result.Text != "" {
				chunk := domain.TranscriptionChunk{
					Text:    result.Text,
					IsFinal: false,
					StartMs: 0,
					EndMs:   len(samples) * 1000 / 16000,
				}
				select {
				case <-ctx.Done():
					return
				case resultChan <- chunk:
				}
			}

			// Send final result
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
		*/

		// For now, send a mock transcription result
		log.Printf("Transcribing %d bytes of audio", len(audio))
		mockText := "Sample transcription from sherpa-onnx"
		chunk := domain.TranscriptionChunk{
			Text:    mockText,
			IsFinal: true,
			StartMs: 0,
			EndMs:   len(samples) * 1000 / 16000,
		}
		select {
		case <-ctx.Done():
			return
		case resultChan <- chunk:
		}
	}()

	return resultChan, nil
}

// TranscribeStream processes audio data in streaming mode
// Audio chunks are sent to the input channel, results come from output channel
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

		// TODO: Uncomment when sherpa-onnx library is available
		/*
			p.mu.Lock()
			// Add left padding at the start
			leftPadding := make([]float32, int(float32(16000)*0.3))
			p.stream.AcceptWaveform(16000, leftPadding)
			p.mu.Unlock()
		*/

		var lastPartialResult string

		for {
			select {
			case <-ctx.Done():
				return

			case audio, ok := <-audioIn:
				if !ok {
					// Channel closed, finalize
					p.mu.Lock()

					// TODO: Uncomment when sherpa-onnx library is available
					/*
						// Add right padding at the end
						rightPadding := make([]float32, int(float32(16000)*0.6))
						p.stream.AcceptWaveform(16000, rightPadding)

						// Finalize decoding
						for p.recognizer.IsReady(p.stream) {
							p.recognizer.Decode(p.stream)
						}

						result := p.recognizer.GetResult(p.stream)
					*/

					// Mock result
					result := &Result{
						Text: "Transcription complete",
					}

					p.mu.Unlock()

					// Send final result
					if result.Text != lastPartialResult {
						chunk := domain.TranscriptionChunk{
							Text:    result.Text,
							IsFinal: true,
						}
						select {
						case <-ctx.Done():
							return
						case resultOut <- chunk:
						}
						lastPartialResult = result.Text
					}

					return
				}

				p.mu.Lock()

				// Convert bytes to float32 samples
				samples := bytesToFloat32(audio)

				// TODO: Uncomment when sherpa-onnx library is available
				/*
					// Accept waveform
					p.stream.AcceptWaveform(16000, samples)

					// Decode if ready
					for p.recognizer.IsReady(p.stream) {
						p.recognizer.Decode(p.stream)
					}

					// Get current result
					result := p.recognizer.GetResult(p.stream)
				*/

				// Mock result based on audio data
				mockTranscript := fmt.Sprintf("Received %d samples", len(samples))
				result := &Result{
					Text: mockTranscript,
				}

				p.mu.Unlock()

				// Send delta event if result changed
				if result.Text != lastPartialResult {
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
					}
					lastPartialResult = result.Text
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

	// TODO: Uncomment when sherpa-onnx library is available
	/*
		if p.stream != nil {
			sherpa.DeleteOnlineStream(p.stream)
			p.stream = nil
		}

		if p.recognizer != nil {
			sherpa.DeleteOnlineRecognizer(p.recognizer)
			p.recognizer = nil
		}
	*/

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
