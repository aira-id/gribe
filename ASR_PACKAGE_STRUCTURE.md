# ASR Providers - Package Structure Reorganization

## New Directory Structure

```
internal/
├── pkg/                           # Provider implementations
│   ├── sherpa/
│   │   └── provider.go           # Sherpa-onnx provider
│   ├── whisper/
│   │   └── provider.go           # Whisper.cpp provider (placeholder)
│   └── mock/
│       └── provider.go           # Mock provider for testing
│
└── usecase/
    ├── asr_factory.go            # Factory pattern for provider creation
    ├── asr_examples.go           # Usage examples
    ├── session_usecase.go        # Main session handler
    └── ... other files
```

## Benefits of This Structure

✅ **Clean Separation of Concerns**
- Each provider is isolated in its own package
- Easy to find and update specific providers

✅ **Improved Discoverability**
- Organized by provider type instead of usecase
- Mirrors standard Go project structure

✅ **Better Modularity**
- Each provider can have additional files if needed (e.g., `sherpa/config.go`, `sherpa/utils.go`)
- Easy to add provider-specific documentation

✅ **Scalability**
- Adding new providers doesn't clutter the usecase directory
- Clear pattern for team collaboration

## File Locations

| Component | Location |
|-----------|----------|
| Sherpa-onnx | `internal/pkg/sherpa/provider.go` |
| Whisper.cpp | `internal/pkg/whisper/provider.go` |
| Mock | `internal/pkg/mock/provider.go` |
| Factory | `internal/usecase/asr_factory.go` |
| Session Integration | `internal/usecase/session_usecase.go` |
| Examples | `internal/usecase/asr_examples.go` |

## Usage (Unchanged from User Perspective)

```go
// Creating providers still works the same way
config := &ASRProviderConfig{
    Type: ProviderSherpaOnnx,
    TranscriptionConfig: &domain.TranscriptionConfig{
        Model:    "zipformer",
        Language: "en",
    },
}

usecase, err := NewSessionUsecaseWithASRProvider(config)
```

## Import Changes (Internal)

**Before:**
```go
import (
    "github.com/aira-id/gribe/internal/usecase"
)
```

**After:**
```go
import (
    "github.com/aira-id/gribe/internal/pkg/sherpa"
    "github.com/aira-id/gribe/internal/pkg/whisper"
    "github.com/aira-id/gribe/internal/pkg/mock"
    "github.com/aira-id/gribe/internal/usecase"
)
```

The factory handles all imports automatically for users.

## Future Provider Structure

When adding a new provider (e.g., Google Cloud Speech-to-Text):

```
internal/pkg/google/
├── provider.go           # Main provider implementation
├── config.go            # Google-specific configuration
├── auth.go              # Authentication logic
└── utils.go             # Helper functions
```

## Migration Guide

For existing code:
- No changes needed! The factory pattern abstracts all provider changes
- Internal imports are handled by the factory
- Users continue using the same API

## Adding a New Provider in New Structure

```go
// Step 1: Create new package
// mkdir -p internal/pkg/newprovider

// Step 2: Create provider.go
// package newprovider
// type Provider struct { ... }
// func New() (*Provider, error) { ... }

// Step 3: Update factory (internal/usecase/asr_factory.go)
import "github.com/aira-id/gribe/internal/pkg/newprovider"

// Add constant
const ProviderNew ASRProviderType = "new-provider"

// Register factory function
factory.Register(ProviderNew, createNewProvider)

// Step 4: Use it
config := &ASRProviderConfig{Type: ProviderNew}
usecase, err := NewSessionUsecaseWithASRProvider(config)
```

## Testing Benefits

Each provider can have its own test file:
```
internal/pkg/sherpa/provider.go
internal/pkg/sherpa/provider_test.go  (future)
internal/pkg/whisper/provider.go
internal/pkg/whisper/provider_test.go (future)
internal/pkg/mock/provider.go
internal/pkg/mock/provider_test.go
```

## Summary

The reorganization:
- ✅ Moves providers to dedicated packages
- ✅ Keeps the same factory pattern
- ✅ Maintains backwards compatibility
- ✅ Improves code organization
- ✅ Makes it easier to add new providers
- ✅ All 16 tests still pass
- ✅ 0 compilation errors
