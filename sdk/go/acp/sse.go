package acp

import (
	"fmt"
	"io"
	"net/http"

	"github.com/manucorporat/sse"
)

type SSEWriter struct {
	writer io.Writer
}

func NewSSEWriter(w io.Writer) *SSEWriter {
	return &SSEWriter{
		writer: w,
	}
}

func (w *SSEWriter) Send(evt Event) error {
	if err := sse.Encode(w.writer, sse.Event{
		Id:    fmt.Sprintf("%s_%d", evt.Type(), evt.Timestamp()),
		Event: string(evt.Type()),
		Data:  evt,
	}); err != nil {
		return err
	}

	if flusher, ok := w.writer.(flusher); ok {
		if err := flusher.Flush(); err != nil {
			return fmt.Errorf("SSE flush failed: %w", err)
		}
	}
	if flusher, ok := w.writer.(flusherWithoutError); ok {
		flusher.Flush()
	}
	return nil
}

type flusher interface {
	Flush() error
}

type flusherWithoutError = http.Flusher
