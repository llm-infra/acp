package acp

import "sync"

const (
	NoneRunID     = "00000000-0000-4000-8000-000000000000"
	NoneBlockID   = "00000000-0000-4000-8000-100000000000"
	NoneContentID = "00000000-0000-4000-8000-200000000000"
)

type EventGenerator struct {
	mutex    sync.Mutex
	index    int
	indexMap map[string]int

	sseWriter   *SSEWriter
	parentRunID string
}

func NewEventGenerator(sseWriter *SSEWriter, parentRunID string) *EventGenerator {
	return &EventGenerator{
		mutex:    sync.Mutex{},
		index:    0,
		indexMap: make(map[string]int),

		sseWriter:   sseWriter,
		parentRunID: parentRunID,
	}
}

func (e *EventGenerator) SendRunStartedEvent() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(e.parentRunID) == 0 {
		e.parentRunID = NoneRunID
	}

	return e.sseWriter.Send(NewRunStartedEvent(e.parentRunID))
}

func (e *EventGenerator) SendRunFinishedEvent(runID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewRunFinishedEvent(runID))
}

func (e *EventGenerator) SendRunInterruptEvent(runID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewRunInterruptEvent(runID))
}

func (e *EventGenerator) SendRunErrorEvent(runID string, message string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewRunErrorEvent(runID, message))
}

func (e *EventGenerator) SendBlockStartEvent(parentBlockID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(parentBlockID) == 0 {
		parentBlockID = NoneBlockID
	}

	return e.sseWriter.Send(NewBlockStartEvent(parentBlockID))
}

func (e *EventGenerator) SendBlockEndEvent(blockID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewBlockEndEvent(blockID))
}

func (e *EventGenerator) SendContentStartEvent(parentContentID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(parentContentID) == 0 {
		parentContentID = NoneContentID
	}

	evt := NewContentStartEvent(parentContentID, e.index)
	e.index++
	e.indexMap[evt.ContentID] = e.index
	return e.sseWriter.Send(evt)
}

func (e *EventGenerator) NewContentDeltaEvent(contentID string, delta string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewContentDeltaEvent(e.indexMap[contentID], delta))
}

func (e *EventGenerator) SendContentEndEvent(contentID string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.sseWriter.Send(NewContentEndEvent(contentID))
}
