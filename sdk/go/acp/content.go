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

/********************************************************/
/*************** Stream Content Structure ***************/
/********************************************************/
type StreamContent interface {
	SType() string
}

type StreamBaseContent struct {
	ContentType string `json:"type"`
}

func NewStreamBaseContent(contentType string) StreamBaseContent {
	return StreamBaseContent{
		ContentType: contentType,
	}
}

func (c StreamBaseContent) SType() string {
	return c.ContentType
}

// 文本流式消息
type StreamTextContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamTextContent(delta string) StreamTextContent {
	return StreamTextContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeText),
		Delta:             delta,
	}
}

// 思考流式消息
type StreamThinkingContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamThinkingContent(delta string) StreamThinkingContent {
	return StreamThinkingContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeThinking),
		Delta:             delta,
	}
}

// 工具调用流式消息
type StreamToolCallContent struct {
	StreamBaseContent

	ToolName string `json:"tool_name"`
}

func NewStreamToolCallContent(toolName string) StreamToolCallContent {
	return StreamToolCallContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeToolCall),
		ToolName:          toolName,
	}
}

// 工具调用参数流式消息
type StreamToolArgsContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamToolArgsContent(delta string) StreamToolArgsContent {
	return StreamToolArgsContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeToolArgs),
		Delta:             delta,
	}
}

// 工具调用结果流式消息
type StreamToolResultContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamToolResultContent(delta string) StreamToolResultContent {
	return StreamToolResultContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeToolResult),
		Delta:             delta,
	}
}

// 文件流式消息
type StreamFileContent struct {
	StreamBaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewStreamFileContent(mimeType string, fileID string) StreamFileContent {
	return StreamFileContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeFile),
		MimeType:          mimeType,
		FileID:            fileID,
	}
}

// 数据流式消息
type StreamDataContent struct {
	StreamBaseContent

	MimeType string `json:"mime_type"`
	Delta    string `json:"delta"` // base64
}

func NewStreamDataContent(mimeType string, data []byte) StreamDataContent {
	return StreamDataContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeData),
		MimeType:          mimeType,
		Delta:             base64.StdEncoding.EncodeToString(data),
	}
}

// 制品流式消息
type StreamArtifactContent struct {
	StreamBaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewStreamArtifactContent(mimeType string, fileID string) StreamArtifactContent {
	return StreamArtifactContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeArtifact),
		MimeType:          mimeType,
		FileID:            fileID,
	}
}

// 变量流式消息
type StreamVariableContent struct {
	StreamBaseContent

	Variables map[string]any `json:"variables"`
}

func NewStreamVariableContent(variables map[string]any) StreamVariableContent {
	return StreamVariableContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeVariable),
		Variables:         variables,
	}
}

// 交互流式消息
type StreamInteractionContent struct {
	StreamBaseContent
}

// 变更流式消息
type StreamPatchContent struct {
	StreamBaseContent
}

// 自定义流式消息
type StreamCustomContent struct {
	StreamBaseContent

	Raw string `json:"raw"`
}

func NewStreamCustomContent(raw string) StreamCustomContent {
	return StreamCustomContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeCustom),
		Raw:               raw,
	}
}

// MCP调用流式消息
type StreamMCPCallContent struct {
	StreamBaseContent

	McpName  string `json:"mcp_name"`
	ToolName string `json:"tool_name"`
}

func NewStreamMCPCallContent(mcpName string, toolName string) StreamMCPCallContent {
	return StreamMCPCallContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeMcpCall),
		McpName:           mcpName,
		ToolName:          toolName,
	}
}

// MCP参数流式消息
type StreamMCPArgsContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamMCPArgsContent(delta string) StreamMCPArgsContent {
	return StreamMCPArgsContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeMcpArgs),
		Delta:             delta,
	}
}

// MCP结果流式消息
type StreamMCPResultContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamMCPResultContent(delta string) StreamMCPResultContent {
	return StreamMCPResultContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeMcpResult),
		Delta:             delta,
	}
}

// 命令执行流式消息
type StreamCommandContent struct {
	StreamBaseContent

	Command string `json:"delta"`
}

func NewStreamCommandContent(delta string) StreamCommandContent {
	return StreamCommandContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeCommandExecution),
		Command:           delta,
	}
}

// 命令执行结果流式消息
type StreamCommandResultContent struct {
	StreamBaseContent

	Delta    string `json:"delta"`
	ExitCode int    `json:"exit_code"`
}

func NewStreamCommandResultContent(delta string, exitCode int) StreamCommandResultContent {
	return StreamCommandResultContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeCommandExecutionResult),
		Delta:             delta,
		ExitCode:          exitCode,
	}
}

// 代码执行
type StreamCodeExecutionContent struct {
	StreamBaseContent

	Lang  string `json:"lang"`
	Delta string `json:"delta"`
}

func NewStreamCodeContent(lang string, delta string) StreamCodeExecutionContent {
	return StreamCodeExecutionContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeCodeExecution),
		Lang:              lang,
		Delta:             delta,
	}
}

// 代码执行结果
type StreamCodeExecutionResultContent struct {
	StreamBaseContent

	Delta string `json:"delta"`
}

func NewStreamCodeResultContent(delta string) StreamCodeExecutionResultContent {
	return StreamCodeExecutionResultContent{
		StreamBaseContent: NewStreamBaseContent(ContentTypeCodeExecutionResult),
		Delta:             delta,
	}
}

// 网络搜索
type StreamNetworkSearchContent struct {
	StreamBaseContent
}

// 网络搜索结果
type StreamNetworkSearchResultContent struct {
	StreamBaseContent
}

// 代办列表
type StreamTodoContent struct {
	StreamBaseContent
}

/********************************************************/
/*************** Session Content Structure **************/
/********************************************************/

type Content interface {
	Type() string
}

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

	ToolName   string `json:"tool_name"`
	ToolArgs   string `json:"tool_args"`
	ToolResult string `json:"tool_result"`
}

func NewToolCallContent(toolName string) *ToolCallContent {
	return &ToolCallContent{
		BaseContent: NewBaseContent(ContentTypeToolCall),
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

func NewDataContent(mimeType string, data []byte) *DataContent {
	return &DataContent{
		BaseContent: NewBaseContent(ContentTypeData),
		MimeType:    mimeType,
		Data:        base64.StdEncoding.EncodeToString(data),
	}
}

// 制品消息
type ArtifactContent struct {
	BaseContent

	MimeType string `json:"mime_type"`
	FileID   string `json:"file_id"`
}

func NewArtifactContent(mimeType string, fileID string) *ArtifactContent {
	return &ArtifactContent{
		BaseContent: NewBaseContent(ContentTypeArtifact),
		MimeType:    mimeType,
		FileID:      fileID,
	}
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

	McpName    string `json:"mcp_name"`
	ToolName   string `json:"tool_name"`
	ToolArgs   string `json:"tool_args"`
	ToolResult string `json:"tool_result"`
}

func NewMCPContent(mcpName, toolName string) *MCPContent {
	return &MCPContent{
		BaseContent: NewBaseContent(ContentTypeMcpCall),
		McpName:     mcpName,
		ToolName:    toolName,
	}
}

// 命令执行
type CommandContent struct {
	BaseContent

	Command  string `json:"command"`
	Result   string `json:"result"`
	ExitCode int    `json:"exit_code"`
}

func NewCommandContent(command string) *CommandContent {
	return &CommandContent{
		BaseContent: NewBaseContent(ContentTypeCommandExecution),
		Command:     command,
	}
}

// 代码执行
type CodeExecutionContent struct {
	BaseContent

	Lang   string `json:"lang"`
	Code   string `json:"code"`
	Result string `json:"result"`
}

func NewCodeExecutionContent(lang, code string) *CodeExecutionContent {
	return &CodeExecutionContent{
		BaseContent: NewBaseContent(ContentTypeCodeExecution),
		Lang:        lang,
		Code:        code,
	}
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
