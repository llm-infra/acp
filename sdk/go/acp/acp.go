package acp

import (
	"encoding/base64"
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
	switch evt := e.(type) {
	case RunStartedEvent:
		m.ID = evt.RunID
		return nil

	case RunFinishedEvent:
		return nil

	case RunErrorEvent:
		m.Errors = evt.Error
		return nil

	default:
		return ErrRunEvent
	}
}

func (m *Creator) processBlockEvent(e Event) error {
	switch evt := e.(type) {
	case BlockStartEvent:
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
	switch evt := e.(type) {
	case ContentStartEvent:
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
		return m.processContent(evt.ContentID, evt.Content)

	case ContentEndEvent:
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
		if evt.Error != nil {
			content.(*ToolCallContent).Error = evt.Error
		} else {
			content.(*ToolCallContent).ToolResult += evt.Delta
		}

	case ContentTypeFile:
		evt, ok := sc.(StreamFileContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewFileContent(evt.MimeType, evt.FileID)
		}

	case ContentTypeData:
		evt, ok := sc.(StreamDataContent)
		if !ok {
			return ErrContentEvent
		}

		// base64 decode
		decoded, err := base64.StdEncoding.DecodeString(evt.Delta)
		if err != nil {
			return err
		}

		if content == nil {
			content = NewDataContent(evt.MimeType, decoded)
		} else {
			content.(*DataContent).Origin = append(content.(*DataContent).Origin, decoded...)
			content.(*DataContent).Data = base64.StdEncoding.EncodeToString(content.(*DataContent).Origin)
		}

	case ContentTypeArtifact:
		evt, ok := sc.(StreamArtifactContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewArtifactContent(evt.MimeType, evt.FileID)
		}

	case ContentTypeVariable:
		evt, ok := sc.(StreamVariableContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewVariableContent(evt.Delta)
		} else {
			for k, v := range evt.Delta {
				content.(*VariableContent).Variables[k] = v
			}
		}

	case ContentTypeInteraction:
		evt, ok := sc.(StreamInteractionContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewInteractionContent(evt.InteractionID, evt.A2UIVersion, evt.A2UIMessage)
		} else {
			interaction, ok := content.(*InteractionContent)
			if !ok {
				return ErrContentEvent
			}
			if interaction.InteractionID == "" && evt.InteractionID != "" {
				interaction.InteractionID = evt.InteractionID
			}
			if interaction.A2UIVersion == "" && evt.A2UIVersion != "" {
				interaction.A2UIVersion = evt.A2UIVersion
			}
			interaction.AddMessage(evt.A2UIMessage)
		}

	case ContentTypeCustom:
		evt, ok := sc.(StreamCustomContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewCustomContent(evt.Raw)
		}

	case ContentTypeMcpCall:
		evt, ok := sc.(StreamMCPCallContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewMCPContent(evt.Server, evt.ToolName)
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
		if evt.Error != nil {
			content.(*MCPContent).Error = evt.Error
		} else {
			content.(*MCPContent).ToolResult += evt.Delta
		}

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
		if evt.Error != nil {
			content.(*CommandContent).Error = evt.Error
		} else {
			content.(*CommandContent).Result += evt.Delta
		}

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
		if evt.Error != nil {
			content.(*CodeExecutionContent).Error = evt.Error
		} else {
			content.(*CodeExecutionContent).Result += evt.Delta
		}

	case ContentTypeWebSearch:
		evt, ok := sc.(StreamWebSearchContent)
		if !ok {
			return ErrContentEvent
		}

		if content == nil {
			content = NewWebSearchContent(evt.Delta)
		}

	case ContentTypeWebSearchResult:
		evt, ok := sc.(StreamWebSearchResultContent)
		if !ok {
			return ErrContentEvent
		}

		if evt.Error != nil {
			content.(*WebSearchContent).Error = evt.Error
		} else {
			content.(*WebSearchContent).Answer = evt.Answer
			content.(*WebSearchContent).Results = evt.Results
		}
	}

	m.contentMap[id] = content
	return nil
}
