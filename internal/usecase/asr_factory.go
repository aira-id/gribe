package usecase

import (
	"github.com/aira-id/gribe/internal/domain"
	"github.com/aira-id/gribe/internal/pkg/mock"
	"github.com/aira-id/gribe/internal/pkg/sherpa"
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

// NewMockProvider creates a mock ASR provider for testing
func NewMockProvider() domain.ASRProvider {
	return mock.New()
}
