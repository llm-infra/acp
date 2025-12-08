package acp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Content Type
const (
	ContentTypeText        = "text"
	ContentTypeThinking    = "thinking"
	ContentTypeToolCall    = "tool_call"
	ContentTypeToolArgs    = "tool_args"
	ContentTypeToolResult  = "tool_result"
	ContentTypeFile        = "file"
	ContentTypeData        = "data"
	ContentTypeArtifact    = "artifact"
	ContentTypeVariable    = "variable"
	ContentTypeInteraction = "interaction"
	ContentTypePatch       = "patch"
	ContentTypeCustom      = "custom"

	ContentTypeMcpCall                = "mcp_call"
	ContentTypeMcpArgs                = "mcp_args"
	ContentTypeMcpResult              = "mcp_result"
	ContentTypeCommandExecution       = "command_execution"
	ContentTypeCommandExecutionResult = "command_execution_result"
	ContentTypeCodeExecution          = "code_execution"
	ContentTypeCodeExecutionResult    = "code_execution_result"
	ContentTypeWebSearch              = "web_search"
	ContentTypeWebSearchResult        = "web_search_result"
	ContentTypeTodoList               = "todo_list"
)

type BaseContent struct {
	ContentType string `json:"type"`
}

func NewBaseContent(contentType string) BaseContent {
	return BaseContent{
		ContentType: contentType,
	}
}

func (c BaseContent) Type() string {
	return c.ContentType
}

/********************************************************/
/*************** Stream Content Structure ***************/
/********************************************************/
type StreamContent interface {
	Type() string
}

// 文本流式消息
type StreamTextContent struct {
	BaseContent

	Delta string `json:"delta"`
}

func NewStreamTextContent(delta string) StreamTextContent {
	return StreamTextContent{
		BaseContent: NewBaseContent(ContentTypeText),
		Delta:       delta,
	}
}

// 思考流式消息
type StreamThinkingContent struct {
	BaseContent

	Delta string `json:"delta"`
}

func NewStreamThinkingContent(delta string) StreamThinkingContent {
	return StreamThinkingContent{
		BaseContent: NewBaseContent(ContentTypeThinking),
		Delta:       delta,
	}
}

// 工具调用流式消息
type StreamToolCallContent struct {
	BaseContent

	CallID   string `json:"call_id"`
	ToolName string `json:"tool_name"`
}

func NewStreamToolCallContent(callID string, toolName string) StreamToolCallContent {
	return StreamToolCallContent{
		BaseContent: NewBaseContent(ContentTypeToolCall),
		CallID:      callID,
		ToolName:    toolName,
	}
}

// 工具调用参数流式消息
type StreamToolArgsContent struct {
	BaseContent

	CallID string `json:"call_id"`
	Delta  string `json:"delta"`
}

func NewStreamToolArgsContent(callID string, delta string) StreamToolArgsContent {
	return StreamToolArgsContent{
		BaseContent: NewBaseContent(ContentTypeToolArgs),
		CallID:      callID,
		Delta:       delta,
	}
}

// 工具调用结果流式消息
type StreamToolResultContent struct {
	BaseContent

	CallID string `json:"call_id"`
	Delta  string `json:"delta"`
}

func NewStreamToolResultContent(callID string, delta string) StreamToolResultContent {
	return StreamToolResultContent{
		BaseContent: NewBaseContent(ContentTypeToolResult),
		CallID:      callID,
		Delta:       delta,
	}
}

// 文件流式消息
type StreamFileContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewStreamFileContent(mimeType string, fileID string) StreamFileContent {
	return StreamFileContent{
		BaseContent: NewBaseContent(ContentTypeFile),
		MimeType:    mimeType,
		FileID:      fileID,
	}
}

// 数据流式消息
type StreamDataContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // base64
}

func NewStreamDataContent(mimeType string, data []byte) StreamDataContent {
	return StreamDataContent{
		BaseContent: NewBaseContent(ContentTypeData),
		MimeType:    mimeType,
		Data:        base64.StdEncoding.EncodeToString(data),
	}
}

// 制品流式消息
type StreamArtifactContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewStreamArtifactContent(mimeType string, fileID string) StreamArtifactContent {
	return StreamArtifactContent{
		BaseContent: NewBaseContent(ContentTypeArtifact),
		MimeType:    mimeType,
		FileID:      fileID,
	}
}

// 变量流式消息
type StreamVariableContent struct {
	BaseContent

	Variables map[string]any `json:"variables"`
}

func NewStreamVariableContent(variables map[string]any) StreamVariableContent {
	return StreamVariableContent{
		BaseContent: NewBaseContent(ContentTypeVariable),
		Variables:   variables,
	}
}

// 交互流式消息
type StreamInteractionContent struct {
	BaseContent
}

// 变更流式消息
type StreamPatchContent struct {
	BaseContent
}

// 自定义流式消息
type StreamCustomContent struct {
	BaseContent

	Raw string `json:"raw"`
}

func NewStreamCustomContent(raw string) StreamCustomContent {
	return StreamCustomContent{
		BaseContent: NewBaseContent(ContentTypeCustom),
		Raw:         raw,
	}
}

// MCP调用流式消息
type StreamMCPCallContent struct {
	BaseContent
}

// MCP参数流式消息
type StreamMCPArgsContent struct {
	BaseContent
}

// MCP结果流式消息
type StreamMCPResultContent struct {
	BaseContent
}

// 命令执行流式消息
type StreamCommandContent struct {
	BaseContent

	CallID  string `json:"call_id"`
	Command string `json:"delta"`
}

func NewStreamCommandContent(callID string, delta string) StreamCommandContent {
	return StreamCommandContent{
		BaseContent: NewBaseContent(ContentTypeCommandExecution),
		Command:     delta,
		CallID:      callID,
	}
}

// 命令执行结果流式消息
type StreamCommandResultContent struct {
	BaseContent

	CallID   string `json:"call_id"`
	Delta    string `json:"delta"`
	ExitCode int    `json:"exit_code"`
}

func NewStreamCommandResultContent(callID string, delta string, exitCode int) StreamCommandResultContent {
	return StreamCommandResultContent{
		BaseContent: NewBaseContent(ContentTypeCommandExecutionResult),
		Delta:       delta,
		ExitCode:    exitCode,
		CallID:      callID,
	}
}

// 网络搜索
type StreamNetworkSearchContent struct {
	BaseContent
}

// 网络搜索结果
type StreamNetworkSearchResultContent struct {
	BaseContent
}

// 代办列表
type StreamTodoContent struct {
	BaseContent
}

/********************************************************/
/*************** Session Content Structure **************/
/********************************************************/

type Content interface {
	Type() string
}

// 文本消息
type TextContent struct {
	BaseContent

	Text string `json:"text"`
}

func NewTextContent(id, text string) *TextContent {
	return &TextContent{
		BaseContent: NewBaseContent(ContentTypeText),
		Text:        text,
	}
}

func (c *TextContent) Append(delta string) { c.Text += delta }

// 思考消息
type ThinkingContent struct {
	BaseContent

	Text string `json:"text"`
}

func NewThinkingContent(id, text string) *ThinkingContent {
	return &ThinkingContent{
		BaseContent: NewBaseContent(ContentTypeThinking),
		Text:        text,
	}
}

func (c *ThinkingContent) Append(delta string) { c.Text += delta }

// 工具调用
type ToolCallContent struct {
	BaseContent

	CallID     string `json:"call_id"`
	ToolName   string `json:"tool_name"`
	ToolArgs   string `json:"tool_args"`
	ToolResult string `json:"tool_result"`
}

func NewToolCallContent(callID, toolName string) *ToolCallContent {
	return &ToolCallContent{
		BaseContent: NewBaseContent(ContentTypeToolCall),
		CallID:      callID,
		ToolName:    toolName,
	}
}

// 文件消息
type FileContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewFileContent(mimeType string, fileID string) *FileContent {
	return &FileContent{
		BaseContent: NewBaseContent(ContentTypeFile),
		MimeType:    mimeType,
		FileID:      fileID,
	}
}

// 数据消息
type DataContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // base64
}

// 制品消息
type ArtifactContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

// 变量消息
type VariableContent struct {
	BaseContent

	Variables map[string]any `json:"variables"`
}

// 交互消息

// 自定义消息

// MCP消息
type MCPContent struct {
	BaseContent
}

// 命令执行
type CommandContent struct {
	BaseContent

	CallID   string `json:"call_id"`
	Command  string `json:"command"`
	Result   string `json:"result"`
	ExitCode int    `json:"exit_code"`
}

func NewCommandContent(callID, command string) *CommandContent {
	return &CommandContent{
		BaseContent: NewBaseContent(ContentTypeCommandExecution),
		CallID:      callID,
		Command:     command,
	}
}

// 代码执行
type CodeExecutionContent struct {
	BaseContent
}

// 网络搜索
type WebSearchContent struct {
	BaseContent
}

// 代办列表
type TodoListContent struct {
	BaseContent
}

func unmarshalContent(data []byte) (Content, error) {
	var base BaseContent
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type() {
	case ContentTypeText:
		var c TextContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeThinking:
		var c ThinkingContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeToolCall:
		var c ToolCallContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeFile:
		var c FileContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeData:
		var c DataContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeArtifact:
		var c ArtifactContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeMcpCall, ContentTypeMcpArgs, ContentTypeMcpResult:
		var c MCPContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeCommandExecution, ContentTypeCommandExecutionResult:
		var c CommandContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeCodeExecution, ContentTypeCodeExecutionResult:
		var c CodeExecutionContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeWebSearch, ContentTypeWebSearchResult:
		var c WebSearchContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	case ContentTypeTodoList:
		var c TodoListContent
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return &c, nil
	default:
		return nil, fmt.Errorf("unsupported content type: %s", base.Type())
	}
}
