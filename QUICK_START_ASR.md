# Quick Start - Modular ASR Providers

## Current Setup (Sherpa-onnx)
```go
config := &ASRProviderConfig{
    Type: ProviderSherpaOnnx,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    "zipformer",
        Language: "en",
    },
}
usecase, err := NewSessionUsecaseWithASRProvider(config)
```

## For Testing (Mock)
```go
config := &ASRProviderConfig{
    Type: ProviderMock,
}
usecase, err := NewSessionUsecaseWithASRProvider(config)
```

## Future (Whisper.cpp)
```go
config := &ASRProviderConfig{
    Type: ProviderWhisperCpp,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model: "base",
    },
    ProviderSpecificConfig: map[string]interface{}{
        "modelPath": "/path/to/ggml-base.bin",
    },
}
usecase, err := NewSessionUsecaseWithASRProvider(config)
```

## Available Providers
- `ProviderSherpaOnnx` - Current default
- `ProviderWhisperCpp` - Placeholder for future
- `ProviderMock` - For testing

## Configuration-Driven
```yaml
audio:
  provider: "sherpa-onnx"    # Switch provider here
  model: "zipformer"
  language: "en"
```

## See Also
- `ASR_MODULAR_DESIGN.md` - Full architecture guide
- `asr_examples.go` - Code examples
- `SHERPA_ONNX_GUIDE.md` - Provider details
