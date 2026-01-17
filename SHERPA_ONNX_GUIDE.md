// File: internal/usecase/asr_sherpa_onnx.go
// 
// SherpaOnnxASRProvider Implementation Guide
// ==========================================
//
// This file implements the ASRProvider interface using sherpa-onnx for automatic speech recognition.
// It provides both batch and streaming transcription capabilities.
//
// Key Features:
// - Implements domain.ASRProvider interface from internal/domain/asr.go
// - Supports both Transcribe() (batch) and TranscribeStream() (streaming) methods
// - Uses OnlineRecognizer from sherpa-onnx with zipformer model as default
// - PCM 16-bit audio format support (16kHz sample rate)
// - Thread-safe operations with sync.Mutex
// - Mock implementation ready (TODO markers for actual sherpa-onnx library integration)
//
// Usage:
// ------
//
// 1. Create a new provider with transcription config:
//    config := &domain.TranscriptionConfig{
//        Model:    "zipformer",
//        Language: "en",
//    }
//    provider, err := NewSherpaOnnxASRProvider(config)
//    if err != nil {
//        // handle error
//    }
//    defer provider.Close()
//
// 2. Use in SessionUsecase:
//    usecase, err := NewSessionUsecaseWithSherpaOnnx(config)
//    if err != nil {
//        // handle error
//    }
//
// 3. Or use as a standalone ASRProvider:
//    resultsChan, err := provider.Transcribe(ctx, audioBytes, config)
//    for chunk := range resultsChan {
//        if chunk.IsFinal {
//            fmt.Println("Final:", chunk.Text)
//        } else {
//            fmt.Println("Partial:", chunk.Text)
//        }
//    }
//
// Configuration:
// ---------------
// The provider uses these default settings:
// - Model Type: zipformer
// - Sample Rate: 16000 Hz
// - Provider: CPU (can be changed to CUDA when library is available)
// - Threads: 4
// - Decoding Method: greedy_search
//
// Integration with Session Handlers:
// -----------------------------------
// The transcribeAudio() function in session_usecase.go already uses the ASRProvider interface:
// 1. Calls provider.Transcribe() with audio bytes and config
// 2. Receives stream of TranscriptionChunk events
// 3. Sends conversation.item.input_audio_transcription.delta events for partial results
// 4. Sends conversation.item.input_audio_transcription.completed event for final result
//
// Actual Sherpa-onnx Integration:
// --------------------------------
// The implementation includes TODO markers where actual sherpa-onnx library calls should go.
// To complete the integration:
//
// 1. Add the sherpa-onnx Go bindings to go.mod:
//    go get github.com/k2-fsa/sherpa-onnx-go
//
// 2. In initializeRecognizer(), uncomment and use:
//    recognizerConfig := &sherpa.OnlineRecognizerConfig{}
//    recognizerConfig.FeatConfig = sherpa.FeatureConfig{
//        SampleRate: 16000,
//        FeatureDim: 80,
//    }
//    recognizerConfig.ModelConfig.Zipformer2Ctc.Model = p.config.Model
//    p.recognizer = sherpa.NewOnlineRecognizer(recognizerConfig)
//    p.stream = sherpa.NewOnlineStream(p.recognizer)
//
// 3. In Transcribe(), uncomment the audio processing logic:
//    p.stream.AcceptWaveform(16000, samples)
//    for p.recognizer.IsReady(p.stream) {
//        p.recognizer.Decode(p.stream)
//    }
//    result := p.recognizer.GetResult(p.stream)
//
// 4. In TranscribeStream(), uncomment the streaming loop logic
//
// 5. In Close(), uncomment the cleanup:
//    sherpa.DeleteOnlineStream(p.stream)
//    sherpa.DeleteOnlineRecognizer(p.recognizer)
//
// Audio Format Conversion:
// -------------------------
// Helper function bytesToFloat32() converts PCM 16-bit little-endian byte arrays to float32 samples
// in the range [-1, 1) as expected by sherpa-onnx.
//
// Error Handling:
// ----------------
// - Returns fmt.Errorf if recognizer not initialized
// - Returns fmt.Errorf if audio data is empty
// - Sends error events through the result channel on context timeout
//
// Thread Safety:
// ---------------
// All operations are protected by sync.Mutex to ensure thread-safe access to:
// - recognizer state
// - active stream
// - results processing
//
// Testing:
// ---------
// Currently, the implementation includes mock transcription results for testing.
// Once sherpa-onnx library is integrated, all TODO sections should be uncommented
// for production use.
