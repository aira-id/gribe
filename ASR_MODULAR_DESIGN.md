# Modular ASR Provider System

## Overview

The ASR (Automatic Speech Recognition) system is now fully modular and extensible. You can easily switch between different ASR implementations (sherpa-onnx, whisper.cpp, etc.) without changing your application logic.

## Architecture

### Core Components

1. **ASRProvider Interface** (`internal/domain/asr.go`)
   - Defines the contract all ASR implementations must follow
   - Methods: `Transcribe()`, `TranscribeStream()`, `GetSupportedModels()`, `GetSupportedLanguages()`, `Close()`

2. **ASRProviderFactory** (`internal/usecase/asr_factory.go`)
   - Central factory for creating provider instances
   - Provider registry for registering custom implementations
   - Configuration-driven provider selection

3. **Provider Implementations**
   - **SherpaOnnxASRProvider** (`internal/usecase/asr_sherpa_onnx.go`) - Current default
   - **WhisperCppASRProvider** (`internal/usecase/asr_whisper_cpp.go`) - Placeholder for future
   - **MockASRProvider** (`internal/usecase/asr_mock.go`) - Testing

## Usage

### Option 1: Using the Factory (Recommended)

```go
// Create factory
factory := NewASRProviderFactory()

// Create provider with configuration
config := &ASRProviderConfig{
    Type: ProviderSherpaOnnx,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    "zipformer",
        Language: "en",
    },
    ProviderSpecificConfig: map[string]interface{}{
        "numThreads": 4,
        "provider":   "cpu",
    },
}

provider, err := factory.Create(config)
if err != nil {
    log.Fatal(err)
}
defer provider.Close()

// Use the provider
resultsChan, err := provider.Transcribe(ctx, audioBytes, config.TranscriptionConfig)
for chunk := range resultsChan {
    fmt.Printf("Text: %s (Final: %v)\n", chunk.Text, chunk.IsFinal)
}
```

### Option 2: With SessionUsecase

```go
// Using factory pattern
asrConfig := &ASRProviderConfig{
    Type: ProviderSherpaOnnx,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    "zipformer",
        Language: "en",
    },
}

usecase, err := NewSessionUsecaseWithASRProvider(asrConfig)
if err != nil {
    log.Fatal(err)
}
```

### Option 3: Direct Provider Instantiation

```go
// For sherpa-onnx
config := &domain.TranscriptionConfig{
    Model:    "zipformer",
    Language: "en",
}
provider, err := NewSherpaOnnxASRProvider(config)
if err != nil {
    log.Fatal(err)
}
```

## Switching Providers

### Switch from Sherpa-onnx to Mock (for testing)

```go
asrConfig := &ASRProviderConfig{
    Type: ProviderMock,
}

usecase, err := NewSessionUsecaseWithASRProvider(asrConfig)
```

### Switch from Sherpa-onnx to Whisper.cpp (future)

```go
asrConfig := &ASRProviderConfig{
    Type: ProviderWhisperCpp,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    "base",
        Language: "en",
    },
    ProviderSpecificConfig: map[string]interface{}{
        "modelPath": "/path/to/ggml-base.bin",
    },
}

usecase, err := NewSessionUsecaseWithASRProvider(asrConfig)
```

## Adding a New Provider

### Step 1: Create Provider Type

Add the provider type to `asr_factory.go`:

```go
const (
    ProviderMyProvider ASRProviderType = "my-provider"
)
```

### Step 2: Create Provider Implementation

Create a new file `internal/usecase/asr_my_provider.go`:

```go
package usecase

import (
    "context"
    "github.com/aira-id/gribe/internal/domain"
)

type MyProvider struct {
    // Your fields here
}

func NewMyProvider(config *domain.TranscriptionConfig) (*MyProvider, error) {
    // Initialize your provider
    return &MyProvider{...}, nil
}

func (p *MyProvider) Transcribe(ctx context.Context, audio []byte, config *domain.TranscriptionConfig) (<-chan domain.TranscriptionChunk, error) {
    resultChan := make(chan domain.TranscriptionChunk)
    go func() {
        defer close(resultChan)
        // Implement transcription logic
    }()
    return resultChan, nil
}

func (p *MyProvider) TranscribeStream(ctx context.Context, config *domain.TranscriptionConfig) (chan<- []byte, <-chan domain.TranscriptionChunk, error) {
    audioIn := make(chan []byte, 100)
    resultOut := make(chan domain.TranscriptionChunk, 10)
    go func() {
        defer close(resultOut)
        // Implement streaming logic
    }()
    return audioIn, resultOut, nil
}

func (p *MyProvider) GetSupportedModels() []string {
    return []string{"model1", "model2"}
}

func (p *MyProvider) GetSupportedLanguages() []string {
    return []string{"en", "es", "fr"}
}

func (p *MyProvider) Close() error {
    // Cleanup
    return nil
}
```

### Step 3: Register Provider

Add a factory function in `asr_factory.go`:

```go
func createMyProvider(config *ASRProviderConfig) (domain.ASRProvider, error) {
    transcConfig := config.TranscriptionConfig
    if transcConfig == nil {
        transcConfig = &domain.TranscriptionConfig{
            Model:    "default",
            Language: "en",
        }
    }
    return NewMyProvider(transcConfig)
}
```

### Step 4: Register in Factory

In `NewASRProviderFactory()`:

```go
factory.Register(ProviderMyProvider, createMyProvider)
```

### Step 5: Use in Application

```go
asrConfig := &ASRProviderConfig{
    Type: ProviderMyProvider,
}
usecase, err := NewSessionUsecaseWithASRProvider(asrConfig)
```

## Configuration-Driven Provider Selection

You can configure the provider via environment variables or config file:

```go
// In config/config.go
type AudioConfig struct {
    Provider                 string                            `yaml:"provider"` // e.g., "sherpa-onnx"
    Model                    string                            `yaml:"model"`
    Language                 string                            `yaml:"language"`
    ProviderSpecificSettings map[string]interface{} `yaml:"provider_settings"`
}

// Usage
asrConfig := &ASRProviderConfig{
    Type: ASRProviderType(cfg.Audio.Provider),
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    cfg.Audio.Model,
        Language: cfg.Audio.Language,
    },
    ProviderSpecificConfig: cfg.Audio.ProviderSpecificSettings,
}
usecase, err := NewSessionUsecaseWithASRProvider(asrConfig)
```

## Supported Providers

| Provider | Status | Notes |
|----------|--------|-------|
| sherpa-onnx | ✅ Active | Default, production-ready for zipformer |
| whisper.cpp | 🔄 Placeholder | Ready for implementation |
| Mock | ✅ Active | For testing and development |

## Future Enhancements

1. **Multiple Provider Fallback**: If one provider fails, automatically try another
2. **Provider Health Checks**: Monitor provider availability and latency
3. **A/B Testing**: Route requests to different providers for comparison
4. **Dynamic Provider Switching**: Change providers without restarting
5. **Provider Pool**: Manage multiple instances of the same provider for load balancing

## Example Configuration (YAML)

```yaml
audio:
  provider: "sherpa-onnx"
  model: "zipformer"
  language: "en"
  provider_settings:
    num_threads: 4
    onnx_provider: "cpu"
    model_path: "./models/sherpa-onnx"

# or for whisper.cpp (future)
#  provider: "whisper-cpp"
#  model: "base"
#  provider_settings:
#    model_path: "./models/ggml-base.bin"
#    n_threads: 4
```

## Testing

```go
func TestASRProviderFactory(t *testing.T) {
    factory := NewASRProviderFactory()
    
    // Test creating different providers
    configs := []ASRProviderType{
        ProviderMock,
        ProviderSherpaOnnx,
    }
    
    for _, providerType := range configs {
        config := &ASRProviderConfig{
            Type: providerType,
        }
        
        provider, err := factory.Create(config)
        if err != nil {
            t.Fatalf("Failed to create %s provider: %v", providerType, err)
        }
        defer provider.Close()
        
        // Test provider methods
        models := provider.GetSupportedModels()
        if len(models) == 0 {
            t.Errorf("%s provider returned no models", providerType)
        }
    }
}
```

## Benefits

✅ **Modularity**: Each provider is independent  
✅ **Extensibility**: Easy to add new providers without modifying core code  
✅ **Testability**: Can swap mock providers for testing  
✅ **Flexibility**: Switch providers via configuration  
✅ **Maintainability**: Clean separation of concerns  
✅ **Future-Proof**: Ready for new ASR technologies
