# Gribe (Go Realtime STT API)

Gribe is an open-source speech-to-text API compatible with OpenAI's Realtime API, built using Golang and WebSockets. It supports multiple ASR providers, including `sherpa-onnx` and `mock` for testing.

## Features
- **OpenAI Compatible**: Implements the OpenAI Realtime API protocol.
- **Modular ASR**: Support for different ASR backends (currently Sherpa-onnx Zipformer2).
- **Indonesian Language Support**: Built-in configuration for Indonesian Zipformer2 models.
- **Graceful Shutdown**: Handles SIGINT/SIGTERM for safe cleanup.
- **Configurable**: Full control via `config.yaml` and environment variables.

## Getting Started

### Prerequisites
- Go 1.21+
- ONNX Runtime libraries (for Sherpa-onnx)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/aira-id/gribe.git
   cd gribe
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Setup ASR models (see [Model Setup](#model-setup) below).

## Model Setup

Gribe requires pre-trained ONNX models to perform speech-to-text. By default, it looks for models in the `models/` directory.

### 1. Download Models
Download the streaming Zipformer models from Hugging Face:

- **Indonesian (ID)**: [sherpa-onnx-streaming-zipformer2-id](https://huggingface.co/spacewave/sherpa-onnx-streaming-zipformer2-id)
- **English (EN)**: [sherpa-onnx-streaming-zipformer-en-2023-06-26](https://huggingface.co/csukuangfj/sherpa-onnx-streaming-zipformer-en-2023-06-26)

### 2. Directory Structure
Place the downloaded folders under the `models/` directory. Your structure should look like this:

```text
gribe/
├── models/
│   ├── sherpa-onnx-streaming-zipformer2-id/
│   │   ├── encoder-iter-100000-avg-15-chunk-32-left-256.onnx
│   │   ├── decoder-iter-100000-avg-15-chunk-32-left-256.onnx
│   │   ├── joiner-iter-100000-avg-15-chunk-32-left-256.onnx
│   │   └── tokens.txt
│   └── sherpa-onnx-streaming-zipformer-en-2023-06-26/
│       ├── encoder-epoch-99-avg-1-chunk-16-left-128.onnx
│       ├── decoder-epoch-99-avg-1-chunk-16-left-128.onnx
│       ├── joiner-epoch-99-avg-1-chunk-16-left-128.onnx
│       └── tokens.txt
└── config.yaml
```

### 3. Verify `config.yaml`
Ensure the file names in `config.yaml` match the files you downloaded.

### Running the Server
```bash
go run main.go
```
The server will start on port `8080` (default).

## Configuration

Gribe uses `config.yaml` for main configuration. Environment variables can also be used for most settings.

### `config.yaml` Structure

```yaml
server:
  port: "8080"
  allowed_origins: [] # List of allowed CORS origins, empty for all

auth:
  api_keys: [] # List of valid API keys for authentication

audio:
  max_audio_buffer_size: 15728640 # Max PCM audio buffer (default 15MB)
  transcription_timeout: "30s"

rate:
  max_connections_per_ip: 10
  requests_per_second: 100
  burst_size: 50
  cleanup_interval: "1m"

asr:
  provider: "cpu" # Default compute provider (cpu or gpu)
  num_threads: 4
  models_dir: "./models"
  default_model: "sherpa-onnx-streaming-zipformer2-id"
  models:
    sherpa-onnx-streaming-zipformer2-id:
      provider: "sherpa-onnx" # Provider for this specific model
      encoder: "encoder-iter-..."
      decoder: "decoder-iter-..."
      joiner: "joiner-iter-..."
      tokens: "tokens.txt"
      languages: ["id", "en"]
```

### Environment Variables
- `GRIBE_PORT`: Server port
- `GRIBE_ALLOWED_ORIGINS`: Comma-separated list of origins
- `GRIBE_API_KEYS`: Comma-separated list of API keys
- `GRIBE_MAX_AUDIO_BUFFER_SIZE`: Buffer size in bytes

## API Usage

### WebSocket Endpoint
`ws://localhost:8080/v1/realtime`

### Client Events
Follows OpenAI Realtime client events:
- `session.update`
- `input_audio_buffer.append`
- `input_audio_buffer.commit`
- `input_audio_buffer.clear`

### Server Events
Follows OpenAI Realtime server events:
- `session.created`
- `session.updated`
- `input_audio_buffer.committed`
- `conversation.item.created`
- `conversation.item.input_audio_transcription.delta`
- `conversation.item.input_audio_transcription.completed`

## Documentation
- [Modular ASR Design](ASR_MODULAR_DESIGN.md)
- [Sherpa-onnx Guide](SHERPA_ONNX_GUIDE.md)
- [ASR Package Structure](ASR_PACKAGE_STRUCTURE.md)
