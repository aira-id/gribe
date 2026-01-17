package usecase

import (
	"context"
	"log"

	"github.com/aira-id/gribe/internal/domain"
)

// Example functions demonstrating modular ASR provider usage

// ExampleUsingFactory demonstrates using the ASR provider factory
func ExampleUsingFactory() {
	factory := NewASRProviderFactory()

	// Create a sherpa-onnx provider
	config := &ASRProviderConfig{
		Type: ProviderSherpaOnnx,
		TranscriptionConfig: &domain.TranscriptionConfig{
			Model:    "zipformer",
			Language: "en",
		},
	}

	provider, err := factory.Create(config)
	if err != nil {
		log.Fatal("Failed to create provider:", err)
	}
	defer provider.Close()

	log.Printf("Provider created: %T", provider)
	log.Printf("Supported models: %v", provider.GetSupportedModels())
	log.Printf("Supported languages: %v", provider.GetSupportedLanguages())
}

// ExampleSwitchingProviders demonstrates switching between different providers
func ExampleSwitchingProviders() {
	factory := NewASRProviderFactory()

	// Get list of available providers
	available := factory.GetSupportedProviders()
	log.Printf("Available providers: %v", available)

	// Try creating each provider
	for _, providerType := range available {
		config := &ASRProviderConfig{
			Type: providerType,
		}

		provider, err := factory.Create(config)
		if err != nil {
			log.Printf("Failed to create %s provider: %v", providerType, err)
			continue
		}

		log.Printf("Successfully created %s provider", providerType)
		models := provider.GetSupportedModels()
		log.Printf("  Models: %d available", len(models))

		provider.Close()
	}
}

// ExampleConfigDriven demonstrates configuration-driven provider selection
func ExampleConfigDriven(providerName string) {
	factory := NewASRProviderFactory()

	config := &ASRProviderConfig{
		Type: ASRProviderType(providerName), // e.g., "sherpa-onnx", "whisper-cpp", "mock"
		TranscriptionConfig: &domain.TranscriptionConfig{
			Model:    "zipformer",
			Language: "en",
		},
	}

	provider, err := factory.Create(config)
	if err != nil {
		log.Fatalf("Failed to create %s provider: %v", providerName, err)
	}
	defer provider.Close()

	log.Printf("Using %s provider with model %s", providerName, config.TranscriptionConfig.Model)
}

// ExampleCustomProvider demonstrates registering a custom provider
func ExampleCustomProvider() {
	factory := NewASRProviderFactory()

	// Define a custom provider type
	customType := ASRProviderType("custom-asr")

	// Register a custom creator function
	factory.Register(customType, func(config *ASRProviderConfig) (domain.ASRProvider, error) {
		// Return any provider instance, in this case a mock
		return NewMockASRProvider(), nil
	})

	// Use the custom provider
	config := &ASRProviderConfig{
		Type: customType,
	}

	provider, err := factory.Create(config)
	if err != nil {
		log.Fatal("Failed to create custom provider:", err)
	}
	defer provider.Close()

	log.Printf("Custom provider registered and created successfully")
}

// ExampleTranscribeWithDifferentProviders demonstrates transcribing with different providers
func ExampleTranscribeWithDifferentProviders(audioBytes []byte) {
	factory := NewASRProviderFactory()
	ctx := context.Background()

	// Example audio bytes (in real usage, this would be actual audio data)
	if len(audioBytes) == 0 {
		log.Println("No audio data provided")
		return
	}

	transcConfig := &domain.TranscriptionConfig{
		Model:    "zipformer",
		Language: "en",
	}

	// Try transcribing with mock provider
	mockConfig := &ASRProviderConfig{
		Type:                   ProviderMock,
		TranscriptionConfig:    transcConfig,
		ProviderSpecificConfig: make(map[string]interface{}),
	}

	provider, err := factory.Create(mockConfig)
	if err != nil {
		log.Fatal("Failed to create mock provider:", err)
	}

	log.Println("Transcribing with mock provider...")
	results, err := provider.Transcribe(ctx, audioBytes, transcConfig)
	if err != nil {
		log.Fatal("Transcription error:", err)
	}

	for chunk := range results {
		log.Printf("  [%v] %s", chunk.IsFinal, chunk.Text)
	}

	provider.Close()

	// Try transcribing with sherpa-onnx provider
	sherpaConfig := &ASRProviderConfig{
		Type:                   ProviderSherpaOnnx,
		TranscriptionConfig:    transcConfig,
		ProviderSpecificConfig: make(map[string]interface{}),
	}

	provider, err = factory.Create(sherpaConfig)
	if err != nil {
		log.Printf("Failed to create sherpa-onnx provider: %v (may need library)", err)
		return
	}

	log.Println("Transcribing with sherpa-onnx provider...")
	results, err = provider.Transcribe(ctx, audioBytes, transcConfig)
	if err != nil {
		log.Fatal("Transcription error:", err)
	}

	for chunk := range results {
		log.Printf("  [%v] %s", chunk.IsFinal, chunk.Text)
	}

	provider.Close()
}

// ExampleSessionWithFactory demonstrates creating SessionUsecase with factory
func ExampleSessionWithFactory() {
	// Using mock provider
	mockConfig := &ASRProviderConfig{
		Type: ProviderMock,
	}

	usecase, err := NewSessionUsecaseWithASRProvider(mockConfig)
	if err != nil {
		log.Fatal("Failed to create session usecase:", err)
	}

	log.Printf("SessionUsecase created with mock provider: %T", usecase.asrProvider)

	// Using sherpa-onnx provider
	sherpaConfig := &ASRProviderConfig{
		Type: ProviderSherpaOnnx,
		TranscriptionConfig: &domain.TranscriptionConfig{
			Model:    "zipformer",
			Language: "en",
		},
	}

	usecase, err = NewSessionUsecaseWithASRProvider(sherpaConfig)
	if err != nil {
		log.Printf("Failed to create session with sherpa-onnx: %v", err)
		return
	}

	log.Printf("SessionUsecase created with sherpa-onnx provider: %T", usecase.asrProvider)
}

// ExampleProviderCapabilities demonstrates checking provider capabilities
func ExampleProviderCapabilities() {
	factory := NewASRProviderFactory()

	providers := factory.GetSupportedProviders()

	for _, providerType := range providers {
		config := &ASRProviderConfig{
			Type: providerType,
		}

		provider, err := factory.Create(config)
		if err != nil {
			log.Printf("⚠️  %s: Failed to initialize (%v)", providerType, err)
			continue
		}

		models := provider.GetSupportedModels()
		languages := provider.GetSupportedLanguages()

		log.Printf("✓ %s", providerType)
		log.Printf("  Models: %d", len(models))
		if len(models) > 0 {
			log.Printf("    Examples: %v", models[:minInt(3, len(models))])
		}
		log.Printf("  Languages: %d", len(languages))
		if len(languages) > 0 {
			log.Printf("    Examples: %v", languages[:minInt(3, len(languages))])
		}

		provider.Close()
	}
}

// Helper function
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
