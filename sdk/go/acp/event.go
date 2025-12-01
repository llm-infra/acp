package acp

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeRunStarted   EventType = "run_started"
	EventTypeRunFinished  EventType = "run_finished"
	EventTypeRunInterrupt EventType = "run_interrupt"
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

	ParentRunID string `json:"parent_run_id"`
	RunID       string `json:"run_id"`
}

func NewRunStartedEvent(parentRunID string) RunStartedEvent {
	return RunStartedEvent{
		BaseEvent:   NewBaseEvent(EventTypeRunStarted),
		ParentRunID: parentRunID,
		RunID:       uuid.NewString(),
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

type RunInterruptEvent struct {
	BaseEvent

	RunID string `json:"run_id"`
}

func NewRunInterruptEvent(runID string) RunInterruptEvent {
	return RunInterruptEvent{
		BaseEvent: NewBaseEvent(EventTypeRunInterrupt),
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
type BlockStartEvent struct {
	BaseEvent

	ParentBlockID string `json:"parent_block_id"`
	BlockID       string `json:"block_id"`
}

func NewBlockStartEvent(parentBlockID string) BlockStartEvent {
	return BlockStartEvent{
		BaseEvent:     NewBaseEvent(EventTypeBlockStart),
		ParentBlockID: parentBlockID,
		BlockID:       uuid.NewString(),
	}
}

type BlockEndEvent struct {
	BaseEvent

	BlockID string `json:"block_id"`
}

func NewBlockEndEvent(blockID string) BlockEndEvent {
	return BlockEndEvent{
		BaseEvent: NewBaseEvent(EventTypeBlockEnd),
		BlockID:   blockID,
	}
}

// 消息事件
type ContentStartEvent struct {
	BaseEvent

	ParentContentID string `json:"parent_content_id"`
	ContentID       string `json:"content_id"`
	Index           int    `json:"index"`
}

func NewContentStartEvent(parentContentID string, index int) ContentStartEvent {
	return ContentStartEvent{
		BaseEvent:       NewBaseEvent(EventTypeContentStart),
		ParentContentID: parentContentID,
		ContentID:       uuid.NewString(),
		Index:           index,
	}
}

type ContentDeltaEvent struct {
	BaseEvent
	Index int    `json:"index"`
	Delta string `json:"delta"`
}

func NewContentDeltaEvent(index int, delta string) ContentDeltaEvent {
	return ContentDeltaEvent{
		BaseEvent: NewBaseEvent(EventTypeContentDelta),
		Index:     index,
		Delta:     delta,
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
