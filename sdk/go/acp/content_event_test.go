package acp

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatorAggregatesInteractionContent(t *testing.T) {
	creator := NewCreator(nil)
	runID := uuid.NewString()
	blockID := uuid.NewString()
	contentID := uuid.NewString()

	require.NoError(t, creator.AddEvent(NewRunStartedEvent(runID)))
	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))

	msg1 := map[string]any{"beginRendering": map[string]any{"surfaceId": "booking"}}
	msg2 := map[string]any{"deleteSurface": map[string]any{"surfaceId": "booking"}}

	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamInteractionContent("itx_1", "0.8", msg1))))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamInteractionContent("itx_1", "0.8", msg2))))
	require.NoError(t, creator.AddEvent(NewContentEndEvent(contentID)))

	require.Len(t, creator.Blocks, 1)
	require.Len(t, creator.Blocks[0].Contents, 1)

	interaction, ok := creator.Blocks[0].Contents[0].(*InteractionContent)
	require.True(t, ok)
	assert.Equal(t, "itx_1", interaction.InteractionID)
	assert.Equal(t, "0.8", interaction.A2UIVersion)
	assert.Len(t, interaction.A2UIMessages, 2)
}

func TestCreatorAggregatesCustomContent(t *testing.T) {
	creator := NewCreator(nil)
	blockID := uuid.NewString()
	contentID := uuid.NewString()

	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamCustomContent("hello "))))
	require.NoError(t, creator.AddEvent(NewContentEndEvent(contentID)))

	require.Len(t, creator.Blocks, 1)
	require.Len(t, creator.Blocks[0].Contents, 1)

	custom, ok := creator.Blocks[0].Contents[0].(*CustomContent)
	require.True(t, ok)
	assert.Equal(t, "hello ", custom.Raw)
}

func TestCreatorCustomContentRejectsDuplicateDelta(t *testing.T) {
	creator := NewCreator(nil)
	blockID := uuid.NewString()
	contentID := uuid.NewString()

	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamCustomContent("full payload"))))

	err := creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamCustomContent("another payload")))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrContentEvent)
}

func TestCreatorAggregatesSkillLoadedContent(t *testing.T) {
	creator := NewCreator(nil)
	blockID := uuid.NewString()
	contentID := uuid.NewString()

	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamSkillLoadedContent("content-type-skill-load"))))
	require.NoError(t, creator.AddEvent(NewContentEndEvent(contentID)))

	require.Len(t, creator.Blocks, 1)
	require.Len(t, creator.Blocks[0].Contents, 1)

	skillLoaded, ok := creator.Blocks[0].Contents[0].(*SkillLoadedContent)
	require.True(t, ok)
	assert.Equal(t, "content-type-skill-load", skillLoaded.Name)
}

func TestCreatorAggregatesQAContent(t *testing.T) {
	creator := NewCreator(nil)
	blockID := uuid.NewString()
	contentID := uuid.NewString()
	options := map[string]any{
		"choices": []string{"A", "B"},
	}

	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamQAContent("qa_1", "confirm", "deploy", "continue?", options))))
	require.NoError(t, creator.AddEvent(NewContentEndEvent(contentID)))

	require.Len(t, creator.Blocks, 1)
	require.Len(t, creator.Blocks[0].Contents, 1)

	qa, ok := creator.Blocks[0].Contents[0].(*QAContent)
	require.True(t, ok)
	assert.Equal(t, "qa_1", qa.QAID)
	assert.Equal(t, "confirm", qa.QAType)
	assert.Equal(t, "deploy", qa.QAName)
	assert.Equal(t, "continue?", qa.Message)
	assert.Equal(t, options, qa.Options)
}

func TestCreatorAggregatesQAResultContentIntoQAContent(t *testing.T) {
	creator := NewCreator(nil)
	blockID := uuid.NewString()
	contentID := uuid.NewString()
	options := map[string]any{
		"choices": []string{"A", "B"},
	}
	answer := map[string]any{
		"choice": "A",
		"confirmed": true,
	}

	require.NoError(t, creator.AddEvent(NewBlockStartEvent(blockID)))
	require.NoError(t, creator.AddEvent(NewContentStartEvent(contentID, blockID)))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamQAContent("qa_1", "confirm", "deploy", "continue?", options))))
	require.NoError(t, creator.AddEvent(NewContentDeltaEvent(contentID, NewStreamQAResultContent(answer))))
	require.NoError(t, creator.AddEvent(NewContentEndEvent(contentID)))

	require.Len(t, creator.Blocks, 1)
	require.Len(t, creator.Blocks[0].Contents, 1)

	qa, ok := creator.Blocks[0].Contents[0].(*QAContent)
	require.True(t, ok)
	assert.Equal(t, answer, qa.Answer)
	assert.Equal(t, options, qa.Options)
}

func TestUnmarshalMessageWithInteractionAndCustom(t *testing.T) {
	payload := `{
		"id": "m1",
		"role": "assistant",
		"blocks": [
			{
				"id": "b1",
				"contents": [
					{
						"type": "interaction",
						"interaction_id": "itx_1",
						"a2ui_version": "0.8",
						"a2ui_messages": [
							{"beginRendering": {"surfaceId": "booking"}}
						]
					},
					{
							"type": "custom",
							"raw": "{\"k\":\"v\"}"
						},
						{
							"type": "skill_loaded",
							"name": "content-type-skill-load"
						}
					]
				}
			],
		"created_at": 1,
		"updated_at": 2
	}`

	var msg Message
	require.NoError(t, json.Unmarshal([]byte(payload), &msg))
	require.Len(t, msg.Blocks, 1)
	require.Len(t, msg.Blocks[0].Contents, 3)

	_, ok1 := msg.Blocks[0].Contents[0].(*InteractionContent)
	_, ok2 := msg.Blocks[0].Contents[1].(*CustomContent)
	_, ok3 := msg.Blocks[0].Contents[2].(*SkillLoadedContent)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
}

func TestUnmarshalMessageWithQAContent(t *testing.T) {
	payload := `{
		"id": "m1",
		"role": "assistant",
		"blocks": [
			{
				"id": "b1",
				"contents": [
					{
						"type": "qa",
						"id": "qa_1",
						"name": "deploy",
						"message": "continue?",
						"options": {
							"choices": ["A", "B"]
						}
					}
				]
			}
		],
		"created_at": 1,
		"updated_at": 2
	}`

	var msg Message
	require.NoError(t, json.Unmarshal([]byte(payload), &msg))
	require.Len(t, msg.Blocks, 1)
	require.Len(t, msg.Blocks[0].Contents, 1)

	qa, ok := msg.Blocks[0].Contents[0].(*QAContent)
	require.True(t, ok)
	assert.Equal(t, "qa_1", qa.QAID)
	assert.Equal(t, "deploy", qa.QAName)
	assert.Equal(t, "continue?", qa.Message)
	assert.Equal(t, map[string]any{"choices": []any{"A", "B"}}, qa.Options)
}

func TestUnmarshalMessageWithQAAnswer(t *testing.T) {
	payload := `{
		"id": "m1",
		"role": "assistant",
		"blocks": [
			{
				"id": "b1",
				"contents": [
					{
						"type": "qa",
						"qa_id": "qa_1",
						"qa_type": "confirm",
						"qa_name": "deploy",
						"message": "continue?",
						"options": {
							"choices": ["A", "B"]
						},
						"answer": {
							"choice": "A",
							"confirmed": true
						}
					}
				]
			}
		],
		"created_at": 1,
		"updated_at": 2
	}`

	var msg Message
	require.NoError(t, json.Unmarshal([]byte(payload), &msg))
	require.Len(t, msg.Blocks, 1)
	require.Len(t, msg.Blocks[0].Contents, 1)

	qa, ok := msg.Blocks[0].Contents[0].(*QAContent)
	require.True(t, ok)
	assert.Equal(t, map[string]any{
		"choice":    "A",
		"confirmed": true,
	}, qa.Answer)
}

func TestUnmarshalMessageWithQAResultContent(t *testing.T) {
	payload := `{
		"id": "m1",
		"role": "assistant",
		"blocks": [
			{
				"id": "b1",
				"contents": [
					{
						"type": "qa_result",
						"answer": {
							"choice": "A"
						}
					}
				]
			}
		],
		"created_at": 1,
		"updated_at": 2
	}`

	var msg Message
	require.NoError(t, json.Unmarshal([]byte(payload), &msg))
	require.Len(t, msg.Blocks, 1)
	require.Len(t, msg.Blocks[0].Contents, 1)

	qa, ok := msg.Blocks[0].Contents[0].(*QAContent)
	require.True(t, ok)
	assert.Equal(t, map[string]any{"choice": "A"}, qa.Answer)
}
