package usecase

import (
	"fmt"

	"github.com/aira-id/gribe/internal/domain"
	"github.com/aira-id/gribe/internal/pkg/mock"
	"github.com/aira-id/gribe/internal/pkg/sherpa"
	"github.com/aira-id/gribe/internal/pkg/whisper"
)

// Ensure domain is used (for ASRProvider interface)
var _ domain.ASRProvider = (*sherpa.Provider)(nil)

// ASRProviderType represents the type of ASR provider
type ASRProviderType string

const (
	// ProviderSherpaOnnx uses sherpa-onnx for speech recognition
	ProviderSherpaOnnx ASRProviderType = "sherpa-onnx"

	// ProviderWhisperCpp uses whisper.cpp for speech recognition
	ProviderWhisperCpp ASRProviderType = "whisper-cpp"

	// ProviderMock uses a mock provider for testing
	ProviderMock ASRProviderType = "mock"
)

// ASRProviderConfig contains configuration for creating an ASR provider
type ASRProviderConfig struct {
	// Type specifies which ASR provider to use
	Type ASRProviderType

	// TranscriptionConfig contains transcription-specific settings
	TranscriptionConfig *domain.TranscriptionConfig

	// ProviderSpecificConfig holds provider-specific settings
	// For sherpa-onnx: model path, num threads, ONNX provider (cpu/cuda)
	// For whisper.cpp: model path, num threads
	ProviderSpecificConfig map[string]interface{}
}

// ASRProviderFactory creates ASR provider instances based on configuration
type ASRProviderFactory struct {
	providers map[ASRProviderType]func(*ASRProviderConfig) (domain.ASRProvider, error)
}

// NewASRProviderFactory creates a new ASR provider factory with default providers
func NewASRProviderFactory() *ASRProviderFactory {
	factory := &ASRProviderFactory{
		providers: make(map[ASRProviderType]func(*ASRProviderConfig) (domain.ASRProvider, error)),
	}

	// Register built-in providers
	factory.Register(ProviderSherpaOnnx, createSherpaOnnxProvider)
	factory.Register(ProviderWhisperCpp, createWhisperCppProvider)
	factory.Register(ProviderMock, createMockProvider)

	return factory
}

// Register registers a custom provider factory function
func (f *ASRProviderFactory) Register(providerType ASRProviderType, creator func(*ASRProviderConfig) (domain.ASRProvider, error)) {
	f.providers[providerType] = creator
}

// Create creates an ASR provider based on the given configuration
func (f *ASRProviderFactory) Create(config *ASRProviderConfig) (domain.ASRProvider, error) {
	if config == nil {
		config = &ASRProviderConfig{
			Type: ProviderSherpaOnnx,
		}
	}

	creator, exists := f.providers[config.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported ASR provider type: %s", config.Type)
	}

	return creator(config)
}

// GetSupportedProviders returns list of supported provider types
func (f *ASRProviderFactory) GetSupportedProviders() []ASRProviderType {
	providers := make([]ASRProviderType, 0, len(f.providers))
	for providerType := range f.providers {
		providers = append(providers, providerType)
	}
	return providers
}

// Provider creation functions

// createSherpaOnnxProvider creates a sherpa-onnx ASR provider
func createSherpaOnnxProvider(config *ASRProviderConfig) (domain.ASRProvider, error) {
	if config.ProviderSpecificConfig == nil {
		return nil, fmt.Errorf("sherpa-onnx provider requires ProviderSpecificConfig")
	}

	// Extract sherpa config from ProviderSpecificConfig
	sherpaConfig := &sherpa.Config{}

	if v, ok := config.ProviderSpecificConfig["provider"].(string); ok {
		sherpaConfig.Provider = v
	}
	if v, ok := config.ProviderSpecificConfig["num_threads"].(int); ok {
		sherpaConfig.NumThreads = v
	}
	if v, ok := config.ProviderSpecificConfig["models_dir"].(string); ok {
		sherpaConfig.ModelsDir = v
	}
	if v, ok := config.ProviderSpecificConfig["model_name"].(string); ok {
		sherpaConfig.ModelName = v
	}
	if v, ok := config.ProviderSpecificConfig["encoder"].(string); ok {
		sherpaConfig.Encoder = v
	}
	if v, ok := config.ProviderSpecificConfig["decoder"].(string); ok {
		sherpaConfig.Decoder = v
	}
	if v, ok := config.ProviderSpecificConfig["joiner"].(string); ok {
		sherpaConfig.Joiner = v
	}
	if v, ok := config.ProviderSpecificConfig["tokens"].(string); ok {
		sherpaConfig.Tokens = v
	}
	if v, ok := config.ProviderSpecificConfig["languages"].([]string); ok {
		sherpaConfig.Languages = v
	}
	if v, ok := config.ProviderSpecificConfig["language"].(string); ok {
		sherpaConfig.Language = v
	}

	return sherpa.New(sherpaConfig)
}

// createWhisperCppProvider creates a whisper.cpp ASR provider
// This is a placeholder for future implementation
func createWhisperCppProvider(config *ASRProviderConfig) (domain.ASRProvider, error) {
	transcConfig := config.TranscriptionConfig
	if transcConfig == nil {
		transcConfig = &domain.TranscriptionConfig{
			Model:    "base",
			Language: "en",
		}
	}

	modelPath := ""
	if config.ProviderSpecificConfig != nil {
		if path, ok := config.ProviderSpecificConfig["modelPath"].(string); ok {
			modelPath = path
		}
	}

	return whisper.New(transcConfig, modelPath)
}

// createMockProvider creates a mock ASR provider for testing
func createMockProvider(config *ASRProviderConfig) (domain.ASRProvider, error) {
	return mock.New(), nil
}
