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
	require.Len(t, msg.Blocks[0].Contents, 2)

	_, ok1 := msg.Blocks[0].Contents[0].(*InteractionContent)
	_, ok2 := msg.Blocks[0].Contents[1].(*CustomContent)
	assert.True(t, ok1)
	assert.True(t, ok2)
}
