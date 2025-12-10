package acp

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mel2oo/go-dkit/json"
	"github.com/stretchr/testify/assert"
)

func TestUserMessage(t *testing.T) {
	msgs := []Message{{
		ID:   uuid.NewString(),
		Role: RoleUser,
		Blocks: []Block{
			{
				ID: uuid.NewString(),
				Contents: []Content{
					&TextContent{
						BaseContent: NewBaseContent(ContentTypeText),
						Text:        "hello",
					},
				},
			},
		},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMicro(),
	}}

	fmt.Println(json.MarshalString(msgs))
}

func TestUnmarshalMessage(t *testing.T) {
	data := `[{"id":"0aa2bd03-f726-41b1-b46e-6d8c5eacb1db","role":"user","block":[{"id":"9a6e53d1-0f4c-4aeb-96af-140cbf11c0d3","content":[{"type":"text","text":"hello"}]}],"created_at":1764811393006,"updated_at":1764811393006124}]
`

	var msgs []Message
	err := json.Unmarshal([]byte(data), &msgs)
	assert.NoError(t, err)
}

func TestBase64(t *testing.T) {
	data := `{
  "private_envs": [],
  "session_round": 5,
  "public_envs": [],
  "org": "1802543778616180738",
  "input_parameters": [
    {
      "name": "input",
      "type": "string"
    }
  ],
  "edges": [],
  "description": "测试测试测试测试",
  "envs": [],
  "cot_version": "v2",
  "uid": "1832953843586764801",
  "entry": "node1",
  "nodes": [
    {
      "max_retries": 3,
      "input_parameters": [
        {
          "name": "input",
          "type": "string",
          "strict": true,
          "value": {
            "expr": "Lazarus的ID是多少",
            "type": "literal"
          }
        }
      ],
      "iteration_values": [],
      "description": "这是一个大模型_CLkOr0",
      "max_working": 30,
      "mcps": [],
      "type": "llm_mixup",
      "tools": [
        {
          "id": "c46eb992-e601-0028-44f7-7454e13364e0"
        },
        {
          "id": "e91738ba-f954-8a80-bf71-53a20105c4bd"
        }
      ],
      "use_memory": false,
      "timeout": 300,
      "agents": [],
      "log_memory": false,
      "result_key": "result_key",
      "output_parameters": [],
      "max_reflections": 5,
      "name": "大模型_CLkOr0",
      "iteration": false,
      "on_error": 0,
      "model": {
        "top_p": 1,
        "max_tokens": 1024,
        "reasoning": 1,
        "top_k": -1,
        "temperature": 0.7,
        "model": "HengNao-v4"
      },
      "id": "node1",
      "no_action": false,
      "no_repeat": false
    }
  ],
  "output_parameters": [],
  "name": "测试测试",
  "model": {
    "top_p": 1,
    "max_tokens": 1024,
    "reasoning": 1,
    "top_k": -1,
    "temperature": 0.7,
    "model": "HengNao-v4"
  },
  "id": "0aa2bd03-f726-41b1-b46e-6d8c5eacb1db",
  "prologue": ""
}`

	fmt.Println(base64.StdEncoding.EncodeToString([]byte(data)))
}
