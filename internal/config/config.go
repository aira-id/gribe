package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Auth   AuthConfig
	Audio  AudioConfig
	Rate   RateLimitConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port           string
	AllowedOrigins []string // Empty means allow all (wildcard)
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeys []string // List of valid API keys, empty means no auth required
}

// AudioConfig holds audio processing limits
type AudioConfig struct {
	Provider             string        // ASR Provider type (e.g., "sherpa-onnx", "mock")
	MaxBufferSize        int           // Maximum audio buffer size in bytes (default 15MB)
	TranscriptionTimeout time.Duration // Timeout for transcription calls (default 30s)
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxConnectionsPerIP int           // Max concurrent connections per IP
	RequestsPerSecond   int           // Max requests per second per connection
	BurstSize           int           // Burst allowance for rate limiting
	CleanupInterval     time.Duration // How often to clean up old entries
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           getEnv("GRIBE_PORT", "8080"),
			AllowedOrigins: getEnvSlice("GRIBE_ALLOWED_ORIGINS", nil), // nil = wildcard
		},
		Auth: AuthConfig{
			APIKeys: getEnvSlice("GRIBE_API_KEYS", nil), // nil = no auth required
		},
		Audio: AudioConfig{
			Provider:             getEnv("GRIBE_ASR_PROVIDER", "sherpa-onnx"),
			MaxBufferSize:        getEnvInt("GRIBE_MAX_AUDIO_BUFFER_SIZE", 15*1024*1024), // 15MB default
			TranscriptionTimeout: time.Duration(getEnvInt("GRIBE_TRANSCRIPTION_TIMEOUT_SECONDS", 30)) * time.Second,
		},
		Rate: RateLimitConfig{
			MaxConnectionsPerIP: getEnvInt("GRIBE_MAX_CONNECTIONS_PER_IP", 10),
			RequestsPerSecond:   getEnvInt("GRIBE_REQUESTS_PER_SECOND", 100),
			BurstSize:           getEnvInt("GRIBE_RATE_BURST_SIZE", 50),
			CleanupInterval:     time.Duration(getEnvInt("GRIBE_RATE_CLEANUP_SECONDS", 60)) * time.Second,
		},
	}
}

// IsOriginAllowed checks if the given origin is allowed
func (c *Config) IsOriginAllowed(origin string) bool {
	// If no origins configured, allow all (wildcard)
	if len(c.Server.AllowedOrigins) == 0 {
		return true
	}

	// Check if origin matches any allowed origin
	for _, allowed := range c.Server.AllowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
	}
	return false
}

// IsAPIKeyValid checks if the given API key is valid
func (c *Config) IsAPIKeyValid(apiKey string) bool {
	// If no API keys configured, allow all (no auth required)
	if len(c.Auth.APIKeys) == 0 {
		return true
	}

	// Check if key matches any configured key
	for _, validKey := range c.Auth.APIKeys {
		if validKey == apiKey {
			return true
		}
	}
	return false
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Split by comma and trim whitespace
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}
