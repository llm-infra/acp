package acp

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRunEvent     = errors.New("run event structure error")
	ErrBlockEvent   = errors.New("block event structure error")
	ErrContentEvent = errors.New("content event structure error")
)

type Creator struct {
	*Message

	hasStarted   bool
	hasFinished  bool
	writer       *SSEWriter
	mux          sync.Mutex
	contentMap   map[string]Content // content_id -> content
	contentIDMap map[string]string  // content_id -> block_id
}

func NewCreator(writer *SSEWriter) *Creator {
	return &Creator{
		Message: &Message{
			ID:        uuid.NewString(),
			Role:      RoleAssistant,
			Blocks:    make([]Block, 0),
			CreatedAt: time.Now().UnixMicro(),
			UpdatedAt: time.Now().UnixMicro(),
		},
		writer:       writer,
		mux:          sync.Mutex{},
		contentMap:   make(map[string]Content),
		contentIDMap: make(map[string]string),
	}
}

func (m *Creator) AddEvent(e Event) (err error) {
	switch e.Type() {
	case EventTypeRunStarted, EventTypeRunFinished, EventTypeRunError:
		if e.Type() == EventTypeRunStarted {
			if m.hasStarted {
				return nil
			}
			m.hasStarted = true
		} else {
			if m.hasFinished {
				return errors.New("event stream already done")
			}
			m.hasFinished = true
		}

		err = m.processRunEvent(e)
	case EventTypeBlockStart, EventTypeBlockEnd:
		err = m.processBlockEvent(e)
	case EventTypeContentStart, EventTypeContentDelta, EventTypeContentEnd:
		err = m.processContentEvent(e)
	default:
		return fmt.Errorf("unsupport event: %s", e.Type())
	}
	if err != nil {
		return err
	}

	if m.writer != nil {
		if err := m.writer.Send(e); err != nil {
			return err
		}
	}

	m.UpdatedAt = time.Now().UnixMicro()
	return nil
}

func (m *Creator) processRunEvent(e Event) error {
	switch e.(type) {
	case RunStartedEvent:
		evt, ok := e.(RunStartedEvent)
		if !ok {
			return ErrRunEvent
		}

		m.ID = evt.RunID
		return nil

	case RunFinishedEvent:
		return nil

	case RunErrorEvent:
		evt, ok := e.(RunErrorEvent)
		if !ok {
			return ErrRunEvent
		}

		m.Errors = evt.Error
		return nil

	default:
		return ErrRunEvent
	}
}

func (m *Creator) processBlockEvent(e Event) error {
	switch e.(type) {
	case BlockStartEvent:
		evt, ok := e.(BlockStartEvent)
		if !ok {
			return ErrBlockEvent
		}

		m.Blocks = append(m.Blocks, Block{
			ID:            evt.BlockID,
			Contents:      make([]Content, 0),
			IsParallel:    evt.IsParallel,
			IsSubagent:    evt.IsSubagent,
			Metadata:      evt.Metadata,
			ParentBlockID: evt.ParentBlockID,
		})
		return nil

	case BlockEndEvent:
		evt, ok := e.(BlockEndEvent)
		if !ok {
			return ErrBlockEvent
		}

		for i := len(m.Blocks) - 1; i >= 0; i-- {
			if m.Blocks[i].ID == evt.BlockID {
				m.Blocks[i].Usage = evt.Usage
				break
			}
		}
		return nil

	default:
		return ErrBlockEvent
	}
}

func (m *Creator) processContentEvent(e Event) error {
	switch e.(type) {
	case ContentStartEvent:
		evt, ok := e.(ContentStartEvent)
		if !ok {
			return ErrContentEvent
		}

		for i := len(m.Blocks) - 1; i >= 0; i-- {
			if m.Blocks[i].ID == evt.RelatedBlockID {
				m.mux.Lock()
				m.contentMap[evt.ContentID] = nil
				m.contentIDMap[evt.ContentID] = evt.RelatedBlockID
				m.mux.Unlock()
				break
			}
		}
		return nil

	case ContentDeltaEvent:
		evt, ok := e.(ContentDeltaEvent)
		if !ok {
			return ErrContentEvent
		}

		return m.processContent(evt.ContentID, evt.Content)

	case ContentEndEvent:
		evt, ok := e.(ContentEndEvent)
		if !ok {
			return ErrContentEvent
		}

		m.mux.Lock()
		content, ok1 := m.contentMap[evt.ContentID]
		blockID, ok2 := m.contentIDMap[evt.ContentID]
		if ok1 && ok2 {
			for i := len(m.Blocks) - 1; i >= 0; i-- {
				if m.Blocks[i].ID == blockID {
					m.Blocks[i].Contents = append(m.Blocks[i].Contents, content)
					break
				}
			}

			delete(m.contentIDMap, evt.ContentID)
			delete(m.contentMap, evt.ContentID)
		}
		m.mux.Unlock()
		return nil

	default:
		return ErrContentEvent
	}
}

func (m *Creator) processContent(id string, sc StreamContent) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	content, ok := m.contentMap[id]
	if !ok {
		return nil
	}

	switch sc.SType() {
	case ContentTypeText:
		evt, ok := sc.(StreamTextContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewTextContent(id, evt.Delta)
		} else {
			content.(*TextContent).Append(evt.Delta)
		}

	case ContentTypeThinking:
		evt, ok := sc.(StreamThinkingContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewThinkingContent(id, evt.Delta)
		} else {
			content.(*ThinkingContent).Append(evt.Delta)
		}

	case ContentTypeToolCall:
		evt, ok := sc.(StreamToolCallContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewToolCallContent(evt.ToolName)
		}

	case ContentTypeToolArgs:
		evt, ok := sc.(StreamToolArgsContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*ToolCallContent).ToolArgs += evt.Delta

	case ContentTypeToolResult:
		evt, ok := sc.(StreamToolResultContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*ToolCallContent).ToolResult += evt.Delta

	case ContentTypeMcpCall:
		evt, ok := sc.(StreamMCPCallContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewMCPContent(evt.McpName, evt.ToolName)
		}

	case ContentTypeMcpArgs:
		evt, ok := sc.(StreamMCPArgsContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*MCPContent).ToolArgs += evt.Delta

	case ContentTypeMcpResult:
		evt, ok := sc.(StreamMCPResultContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*MCPContent).ToolResult += evt.Delta

	case ContentTypeCommandExecution:
		evt, ok := sc.(StreamCommandContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewCommandContent(evt.Command)
		}

	case ContentTypeCommandExecutionResult:
		evt, ok := sc.(StreamCommandResultContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*CommandContent).Result += evt.Delta

	case ContentTypeCodeExecution:
		evt, ok := sc.(StreamCodeExecutionContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewCodeExecutionContent(evt.Lang, evt.Delta)
		} else {
			content.(*CodeExecutionContent).Code += evt.Delta
		}

	case ContentTypeCodeExecutionResult:
		evt, ok := sc.(StreamCodeExecutionResultContent)
		if !ok || content == nil {
			return ErrContentEvent
		}
		content.(*CodeExecutionContent).Result += evt.Delta
	}

	m.contentMap[id] = content
	return nil
}
