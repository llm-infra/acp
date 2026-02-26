package acp

import "encoding/json"

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type Message struct {
	ID        string  `json:"id"`
	Role      string  `json:"role"`
	Blocks    []Block `json:"blocks"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	Errors    string  `json:"errors,omitempty"`
}

func (m *Message) GetInputs() (*TextContent, []*FileContent) {
	var text *TextContent
	var files = []*FileContent{}

	for _, b := range m.Blocks {
		for _, c := range b.Contents {
			switch c := c.(type) {
			case *TextContent:
				text = c
			case *FileContent:
				files = append(files, c)
			}
		}
	}

	return text, files
}

func (m *Message) GetVariables() (*TextContent, *VariableContent) {
	var text *TextContent
	var variables *VariableContent

	for _, b := range m.Blocks {
		for _, c := range b.Contents {
			switch c := c.(type) {
			case *TextContent:
				text = c
			case *VariableContent:
				variables = c
			}
		}
	}
	return text, variables
}

type Block struct {
	ID            string         `json:"id"`
	Contents      []Content      `json:"contents"`
	Usage         *Usage         `json:"usage,omitempty"`
	IsParallel    bool           `json:"is_parallel,omitempty"`
	IsSubagent    bool           `json:"is_subagent,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	ParentBlockID string         `json:"parent_block_id,omitempty"`
}

func (b *Block) UnmarshalJSON(data []byte) error {
	type rawBlock struct {
		ID            string            `json:"id"`
		Contents      []json.RawMessage `json:"contents"`
		Usage         *Usage            `json:"usage,omitempty"`
		IsParallel    bool              `json:"is_parallel,omitempty"`
		IsSubagent    bool              `json:"is_subagent,omitempty"`
		Metadata      map[string]any    `json:"metadata,omitempty"`
		ParentBlockID string            `json:"parent_block_id,omitempty"`
	}

	var rb rawBlock
	if err := json.Unmarshal(data, &rb); err != nil {
		return err
	}

	b.ID = rb.ID
	b.Usage = rb.Usage
	b.IsParallel = rb.IsParallel
	b.IsSubagent = rb.IsSubagent
	b.Metadata = rb.Metadata
	b.ParentBlockID = rb.ParentBlockID
	b.Contents = make([]Content, 0, len(rb.Contents))

	for _, raw := range rb.Contents {
		content, err := unmarshalContent(raw)
		if err != nil {
			return err
		}
		b.Contents = append(b.Contents, content)
	}

	return nil
}
