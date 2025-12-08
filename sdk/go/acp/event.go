package acp

import (
	"time"
)

type EventType string

const (
	EventTypeRunStarted   EventType = "run_started"
	EventTypeRunFinished  EventType = "run_finished"
	EventTypeRunError     EventType = "run_error"
	EventTypeBlockStart   EventType = "block_start"
	EventTypeBlockEnd     EventType = "block_end"
	EventTypeContentStart EventType = "content_start"
	EventTypeContentDelta EventType = "content_delta"
	EventTypeContentEnd   EventType = "content_end"
)

type Event interface {
	Type() EventType
	Timestamp() int64
}

type BaseEvent struct {
	EventType   EventType `json:"type"`
	TimestampMs int64     `json:"timestamp"`
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}

func (e BaseEvent) Timestamp() int64 {
	return e.TimestampMs
}

func NewBaseEvent(eventType EventType) BaseEvent {
	return BaseEvent{
		EventType:   eventType,
		TimestampMs: time.Now().UnixMilli(),
	}
}

// 生命周期事件
type RunStartedEvent struct {
	BaseEvent

	RunID string `json:"run_id"`
}

func NewRunStartedEvent(runID string) RunStartedEvent {
	return RunStartedEvent{
		BaseEvent: NewBaseEvent(EventTypeRunStarted),
		RunID:     runID,
	}
}

type RunFinishedEvent struct {
	BaseEvent

	RunID string `json:"run_id"`
}

func NewRunFinishedEvent(runID string) RunFinishedEvent {
	return RunFinishedEvent{
		BaseEvent: NewBaseEvent(EventTypeRunFinished),
		RunID:     runID,
	}
}

type RunErrorEvent struct {
	BaseEvent

	RunID string `json:"run_id"`
	Error string `json:"error"`
}

func NewRunErrorEvent(runID string, message string) RunErrorEvent {
	return RunErrorEvent{
		BaseEvent: NewBaseEvent(EventTypeRunError),
		RunID:     runID,
		Error:     message,
	}
}

// 区块事件
type BlockOption func(*BlockStartEvent)

type BlockStartEvent struct {
	BaseEvent

	BlockID       string         `json:"block_id"`
	IsParallel    bool           `json:"is_parallel,omitempty"`
	IsSubagent    bool           `json:"is_subagent,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	ParentBlockID string         `json:"parent_block_id,omitempty"`
}

func WithIsParallel() BlockOption {
	return func(e *BlockStartEvent) { e.IsParallel = true }
}

func WithIsSubagent() BlockOption {
	return func(e *BlockStartEvent) { e.IsSubagent = true }
}

func WithMetadata(metadata map[string]any) BlockOption {
	return func(e *BlockStartEvent) { e.Metadata = metadata }
}

func WithParentBlockID(parentBlockID string) BlockOption {
	return func(e *BlockStartEvent) { e.ParentBlockID = parentBlockID }
}

func NewBlockStartEvent(blockID string, opts ...BlockOption) BlockStartEvent {
	evt := BlockStartEvent{
		BaseEvent: NewBaseEvent(EventTypeBlockStart),
		BlockID:   blockID,
	}

	for _, o := range opts {
		o(&evt)
	}

	return evt
}

type BlockEndEvent struct {
	BaseEvent

	BlockID string `json:"block_id"`
	Usage   *Usage `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
}

func NewBlockEndEvent(blockID string, usage *Usage) BlockEndEvent {
	return BlockEndEvent{
		BaseEvent: NewBaseEvent(EventTypeBlockEnd),
		BlockID:   blockID,
		Usage:     usage,
	}
}

// 消息事件
type ContentStartEvent struct {
	BaseEvent

	ContentID      string `json:"content_id"`
	RelatedBlockID string `json:"related_block_id,omitempty"`
}

func NewContentStartEvent(contentID, blockID string) ContentStartEvent {
	return ContentStartEvent{
		BaseEvent:      NewBaseEvent(EventTypeContentStart),
		ContentID:      contentID,
		RelatedBlockID: blockID,
	}
}

type ContentEndEvent struct {
	BaseEvent

	ContentID string `json:"content_id"`
}

func NewContentEndEvent(contentID string) ContentEndEvent {
	return ContentEndEvent{
		BaseEvent: NewBaseEvent(EventTypeContentEnd),
		ContentID: contentID,
	}
}

type ContentDeltaEvent struct {
	BaseEvent

	ContentID string        `json:"content_id"`
	Content   StreamContent `json:"content"`
}

func NewContentDeltaEvent(contentID string, content StreamContent) ContentDeltaEvent {
	return ContentDeltaEvent{
		BaseEvent: NewBaseEvent(EventTypeContentDelta),
		ContentID: contentID,
		Content:   content,
	}
}
