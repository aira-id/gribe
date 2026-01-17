package domain

// ============================================================================
// CLIENT EVENTS
// ============================================================================

// SessionUpdateClientEvent represents a session.update client event
type SessionUpdateClientEvent struct {
	BaseEvent
	Session *Session `json:"session"`
}

// InputAudioBufferAppendEvent represents input_audio_buffer.append event
type InputAudioBufferAppendEvent struct {
	BaseEvent
	Audio string `json:"audio"` // base64-encoded audio bytes
}

// InputAudioBufferCommitEvent represents input_audio_buffer.commit event
type InputAudioBufferCommitEvent struct {
	BaseEvent
}

// InputAudioBufferClearEvent represents input_audio_buffer.clear event
type InputAudioBufferClearEvent struct {
	BaseEvent
}

// ConversationItemCreateClientEvent represents conversation.item.create event
type ConversationItemCreateClientEvent struct {
	BaseEvent
	Item           *Item   `json:"item"`
	PreviousItemID *string `json:"previous_item_id,omitempty"` // null, "root", or item ID
}

// ConversationItemRetrieveEvent represents conversation.item.retrieve event
type ConversationItemRetrieveEvent struct {
	BaseEvent
	ItemID string `json:"item_id"`
}

// ConversationItemTruncateEvent represents conversation.item.truncate event
type ConversationItemTruncateEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"` // Usually 0
	AudioEndMs   int    `json:"audio_end_ms"`
}

// ConversationItemDeleteEvent represents conversation.item.delete event
type ConversationItemDeleteEvent struct {
	BaseEvent
	ItemID string `json:"item_id"`
}

// ResponseCreateClientEvent represents response.create event
type ResponseCreateClientEvent struct {
	BaseEvent
	Response *ResponseCreatePayload `json:"response,omitempty"`
}

// ResponseCreatePayload represents the response payload in response.create
type ResponseCreatePayload struct {
	Modalities       []string      `json:"modalities,omitempty"`
	Instructions     string        `json:"instructions,omitempty"`
	Tools            []Tool        `json:"tools,omitempty"`
	ToolChoice       string        `json:"tool_choice,omitempty"`
	Conversation     string        `json:"conversation,omitempty"` // "default" or "none"
	OutputModalities []string      `json:"output_modalities,omitempty"`
	Metadata         interface{}   `json:"metadata,omitempty"`
	Input            []interface{} `json:"input,omitempty"` // Items or item references
	MaxOutputTokens  interface{}   `json:"max_output_tokens,omitempty"`
}

// ResponseCancelEvent represents response.cancel event
type ResponseCancelEvent struct {
	BaseEvent
	ResponseID *string `json:"response_id,omitempty"`
}

// OutputAudioBufferClearEvent represents output_audio_buffer.clear event
type OutputAudioBufferClearEvent struct {
	BaseEvent
}

// ============================================================================
// SERVER EVENTS
// ============================================================================

// ErrorServerEvent represents an error event
type ErrorServerEvent struct {
	BaseEvent
	Error *ErrorDetail `json:"error"`
}

// SessionCreatedEvent represents session.created event
type SessionCreatedEvent struct {
	BaseEvent
	Session *Session `json:"session"`
}

// SessionUpdatedEvent represents session.updated event
type SessionUpdatedEvent struct {
	BaseEvent
	Session *Session `json:"session"`
}

// ConversationItemAddedEvent represents conversation.item.added event
type ConversationItemAddedEvent struct {
	BaseEvent
	Item           *Item   `json:"item"`
	PreviousItemID *string `json:"previous_item_id"`
}

// ConversationItemDoneEvent represents conversation.item.done event
type ConversationItemDoneEvent struct {
	BaseEvent
	Item           *Item   `json:"item"`
	PreviousItemID *string `json:"previous_item_id"`
}

// ConversationItemRetrievedEvent represents conversation.item.retrieved event
type ConversationItemRetrievedEvent struct {
	BaseEvent
	Item *Item `json:"item"`
}

// ConversationItemDeletedEvent represents conversation.item.deleted event
type ConversationItemDeletedEvent struct {
	BaseEvent
	ItemID string `json:"item_id"`
}

// ConversationItemTruncatedEvent represents conversation.item.truncated event
type ConversationItemTruncatedEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	AudioEndMs   int    `json:"audio_end_ms"`
}

// InputAudioBufferCommittedEvent represents input_audio_buffer.committed event
type InputAudioBufferCommittedEvent struct {
	BaseEvent
	ItemID         string  `json:"item_id"`
	PreviousItemID *string `json:"previous_item_id"`
}

// InputAudioBufferClearedEvent represents input_audio_buffer.cleared event
type InputAudioBufferClearedEvent struct {
	BaseEvent
}

// InputAudioBufferSpeechStartedEvent represents input_audio_buffer.speech_started event
type InputAudioBufferSpeechStartedEvent struct {
	BaseEvent
	AudioStartMs int    `json:"audio_start_ms"`
	ItemID       string `json:"item_id"`
}

// InputAudioBufferSpeechStoppedEvent represents input_audio_buffer.speech_stopped event
type InputAudioBufferSpeechStoppedEvent struct {
	BaseEvent
	AudioEndMs int    `json:"audio_end_ms"`
	ItemID     string `json:"item_id"`
}

// InputAudioBufferTimeoutTriggeredEvent represents input_audio_buffer.timeout_triggered event
type InputAudioBufferTimeoutTriggeredEvent struct {
	BaseEvent
	AudioStartMs int    `json:"audio_start_ms"`
	AudioEndMs   int    `json:"audio_end_ms"`
	ItemID       string `json:"item_id"`
}

// ResponseCreatedEvent represents response.created event
type ResponseCreatedEvent struct {
	BaseEvent
	Response *Response `json:"response"`
}

// ResponseDoneEvent represents response.done event
type ResponseDoneEvent struct {
	BaseEvent
	Response *Response `json:"response"`
}

// ResponseOutputItemAddedEvent represents response.output_item.added event
type ResponseOutputItemAddedEvent struct {
	BaseEvent
	ResponseID  string `json:"response_id"`
	OutputIndex int    `json:"output_index"`
	Item        *Item  `json:"item"`
}

// ResponseOutputItemDoneEvent represents response.output_item.done event
type ResponseOutputItemDoneEvent struct {
	BaseEvent
	ResponseID  string `json:"response_id"`
	OutputIndex int    `json:"output_index"`
	Item        *Item  `json:"item"`
}

// ResponseContentPartAddedEvent represents response.content_part.added event
type ResponseContentPartAddedEvent struct {
	BaseEvent
	ResponseID   string       `json:"response_id"`
	ItemID       string       `json:"item_id"`
	ContentIndex int          `json:"content_index"`
	OutputIndex  int          `json:"output_index"`
	Part         *ContentPart `json:"part"`
}

// ResponseContentPartDoneEvent represents response.content_part.done event
type ResponseContentPartDoneEvent struct {
	BaseEvent
	ResponseID   string       `json:"response_id"`
	ItemID       string       `json:"item_id"`
	ContentIndex int          `json:"content_index"`
	OutputIndex  int          `json:"output_index"`
	Part         *ContentPart `json:"part"`
}

// ResponseOutputTextDeltaEvent represents response.output_text.delta event
type ResponseOutputTextDeltaEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
	Delta        string `json:"delta"`
}

// ResponseOutputTextDoneEvent represents response.output_text.done event
type ResponseOutputTextDoneEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
	Text         string `json:"text"`
}

// ResponseOutputAudioTranscriptDeltaEvent represents response.output_audio_transcript.delta event
type ResponseOutputAudioTranscriptDeltaEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
	Delta        string `json:"delta"`
}

// ResponseOutputAudioTranscriptDoneEvent represents response.output_audio_transcript.done event
type ResponseOutputAudioTranscriptDoneEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
	Transcript   string `json:"transcript"`
}

// ResponseOutputAudioDeltaEvent represents response.output_audio.delta event
type ResponseOutputAudioDeltaEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
	Delta        string `json:"delta"` // base64-encoded audio
}

// ResponseOutputAudioDoneEvent represents response.output_audio.done event
type ResponseOutputAudioDoneEvent struct {
	BaseEvent
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	OutputIndex  int    `json:"output_index"`
}

// ConversationItemInputAudioTranscriptionCompletedEvent represents conversation.item.input_audio_transcription.completed event
type ConversationItemInputAudioTranscriptionCompletedEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
	Usage        *Usage `json:"usage"`
}

// ConversationItemInputAudioTranscriptionDeltaEvent represents conversation.item.input_audio_transcription.delta event
type ConversationItemInputAudioTranscriptionDeltaEvent struct {
	BaseEvent
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

// RateLimitsUpdatedEvent represents rate_limits.updated event
type RateLimitsUpdatedEvent struct {
	BaseEvent
	RateLimits []RateLimit `json:"rate_limits"`
}
