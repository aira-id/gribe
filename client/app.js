/**
 * Gribe Realtime STT Web Client
 * Compatible with OpenAI Realtime API
 */

class GribeClient {
    constructor() {
        this.ws = null;
        this.audioContext = null;
        this.processor = null;
        this.stream = null;
        this.isConnected = false;
        this.isRecording = false;

        // UI Elements
        this.wsUrlInput = document.getElementById('ws-url');
        this.apiKeyInput = document.getElementById('api-key');
        this.connectBtn = document.getElementById('connect-btn');
        this.disconnectBtn = document.getElementById('disconnect-btn');
        this.micBtn = document.getElementById('mic-btn');
        this.statusDot = document.getElementById('status-dot');
        this.statusText = document.getElementById('status-text');
        this.sessionInfo = document.getElementById('session-info');
        this.transcriptionLog = document.getElementById('transcription-log');
        this.currentTranscription = document.getElementById('current-transcription');
        this.modelSelect = document.getElementById('model-select');
        this.languageSelect = document.getElementById('language-select');

        this.setupEventListeners();
    }

    setupEventListeners() {
        this.connectBtn.onclick = () => this.connect();
        this.disconnectBtn.onclick = () => this.disconnect();
        this.micBtn.onclick = () => this.toggleMicrophone();

        this.modelSelect.onchange = () => this.updateSession();
        this.languageSelect.onchange = () => this.updateSession();
    }

    log(message, type = 'status') {
        const item = document.createElement('div');
        item.className = `log-item ${type}`;

        const timestamp = document.createElement('span');
        timestamp.className = 'timestamp';
        timestamp.textContent = new Date().toLocaleTimeString();

        const content = document.createElement('div');
        content.textContent = message;

        item.appendChild(timestamp);
        item.appendChild(content);

        this.transcriptionLog.appendChild(item);
        this.transcriptionLog.scrollTop = this.transcriptionLog.scrollHeight;
    }

    updateStatus(status) {
        this.statusDot.className = 'status-dot';
        switch (status) {
            case 'disconnected':
                this.statusDot.classList.add('status-disconnected');
                this.statusText.textContent = 'Disconnected';
                this.connectBtn.disabled = false;
                this.disconnectBtn.disabled = true;
                this.micBtn.disabled = true;
                this.isConnected = false;
                break;
            case 'connecting':
                this.statusDot.classList.add('status-connecting');
                this.statusText.textContent = 'Connecting...';
                this.connectBtn.disabled = true;
                break;
            case 'connected':
                this.statusDot.classList.add('status-connected');
                this.statusText.textContent = 'Connected';
                this.connectBtn.disabled = true;
                this.disconnectBtn.disabled = false;
                this.micBtn.disabled = false;
                this.isConnected = true;
                break;
        }
    }

    async connect() {
        const url = this.wsUrlInput.value;
        const apiKey = this.apiKeyInput.value;

        this.updateStatus('connecting');
        this.log(`Connecting to ${url}...`);

        try {
            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                this.updateStatus('connected');
                this.log('WebSocket Connected', 'status');

                // Send initial session configuration
                this.updateSession();
            };

            this.ws.onclose = () => {
                this.updateStatus('disconnected');
                this.log('WebSocket Disconnected', 'status');
                this.stopMicrophone(); // Ensure mic stops if WS closes
            };

            this.ws.onerror = (error) => {
                this.log('WebSocket Error', 'error');
                console.error(error);
            };

            this.ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                this.handleServerEvent(data);
            };

        } catch (error) {
            this.log(`Connection failed: ${error.message}`, 'error');
            this.updateStatus('disconnected');
        }
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }

    updateSession() {
        if (!this.isConnected) return;

        const model = this.modelSelect.value;
        const language = this.languageSelect.value;

        this.log(`Updating session: model=${model}, language=${language}`, 'status');

        this.sendEvent({
            type: 'session.update',
            session: {
                audio: {
                    input: {
                        transcription: {
                            model: model,
                            language: language
                        }
                    }
                }
            }
        });
    }

    sendEvent(event) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(event));
        }
    }

    handleServerEvent(event) {
        console.log('Server Event:', event.type, event);

        switch (event.type) {
            case 'session.created':
                this.sessionInfo.textContent = `Session: ${event.session.id}`;
                this.log(`Session created: ${event.session.id}`, 'status');
                break;

            case 'error':
                this.log(`Server Error: ${event.error.message}`, 'error');
                break;

            case 'conversation.item.input_audio_transcription.delta':
                this.updateCurrentTranscription(event.delta, true);
                break;

            case 'conversation.item.input_audio_transcription.completed':
                this.log(event.transcript, 'transcript');
                this.clearCurrentTranscription();
                break;

            case 'response.audio_transcript.delta':
                // For STT service, this might not be relevant unless it's full conversational
                this.updateCurrentTranscription(event.delta, true);
                break;

            case 'response.done':
                // Reset or final touch
                break;
        }
    }

    updateCurrentTranscription(text, isDelta = false) {
        if (isDelta) {
            if (this.currentTranscription.querySelector('.delta-text')) {
                const deltaSpan = this.currentTranscription.querySelector('.delta-text');
                deltaSpan.textContent += text;
            } else {
                this.currentTranscription.innerHTML = `<span class="delta-text">${text}</span>`;
            }
        } else {
            this.currentTranscription.textContent = text;
        }
    }

    clearCurrentTranscription() {
        this.currentTranscription.innerHTML = '<span class="text-secondary">Waiting for speech...</span>';
    }

    async toggleMicrophone() {
        if (this.isRecording) {
            this.stopMicrophone();
        } else {
            await this.startMicrophone();
        }
    }

    async startMicrophone() {
        try {
            this.stream = await navigator.mediaDevices.getUserMedia({
                audio: {
                    channelCount: 1,
                    sampleRate: 16000,
                    echoCancellation: true,
                    noiseSuppression: true
                }
            });

            this.audioContext = new (window.AudioContext || window.webkitAudioContext)({
                sampleRate: 16000
            });

            const source = this.audioContext.createMediaStreamSource(this.stream);

            // ScriptProcessorNode is deprecated but widely supported for this use case
            // Buffer size 4096 is standard
            this.processor = this.audioContext.createScriptProcessor(4096, 1, 1);

            source.connect(this.processor);
            this.processor.connect(this.audioContext.destination);

            this.processor.onaudioprocess = (e) => {
                if (!this.isRecording) return;

                const inputData = e.inputBuffer.getChannelData(0);
                const pcmData = this.float32ToInt16(inputData);
                const base64Audio = this.arrayBufferToBase64(pcmData.buffer);

                this.sendEvent({
                    type: 'input_audio_buffer.append',
                    audio: base64Audio
                });
            };

            this.isRecording = true;
            this.micBtn.classList.remove('btn-success');
            this.micBtn.classList.add('btn-danger');
            document.getElementById('mic-text').textContent = 'Stop Microphone';
            document.getElementById('mic-icon').textContent = '⏹️';
            this.log('Microphone started', 'status');

        } catch (error) {
            this.log(`Microphone error: ${error.message}`, 'error');
            console.error(error);
        }
    }

    stopMicrophone() {
        if (!this.isRecording) return;

        this.isRecording = false;

        // Commit the buffer
        this.sendEvent({ type: 'input_audio_buffer.commit' });

        if (this.processor) {
            this.processor.disconnect();
            this.processor = null;
        }

        if (this.stream) {
            this.stream.getTracks().forEach(track => track.stop());
            this.stream = null;
        }

        if (this.audioContext) {
            this.audioContext.close();
            this.audioContext = null;
        }

        this.micBtn.classList.remove('btn-danger');
        this.micBtn.classList.add('btn-success');
        document.getElementById('mic-text').textContent = 'Start Microphone';
        document.getElementById('mic-icon').textContent = '🎤';
        this.log('Microphone stopped', 'status');
    }

    // Utility: Convert Float32Array to Int16Array (PCM)
    float32ToInt16(buffer) {
        let l = buffer.length;
        let buf = new Int16Array(l);
        while (l--) {
            let s = Math.max(-1, Math.min(1, buffer[l]));
            buf[l] = s < 0 ? s * 0x8000 : s * 0x7FFF;
        }
        return buf;
    }

    // Utility: ArrayBuffer to Base64
    arrayBufferToBase64(buffer) {
        let binary = '';
        const bytes = new Uint8Array(buffer);
        const len = bytes.byteLength;
        for (let i = 0; i < len; i++) {
            binary += String.fromCharCode(bytes[i]);
        }
        return window.btoa(binary);
    }
}

// Initialize client
window.onload = () => {
    window.gribe = new GribeClient();
};
