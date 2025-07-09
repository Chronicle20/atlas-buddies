package list

import (
	"encoding/json"
	list2 "atlas-buddies/kafka/message/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCapacityCommandProvider(t *testing.T) {
	characterId := uint32(12345)
	worldId := byte(1)
	capacity := byte(50)

	t.Run("create update capacity command", func(t *testing.T) {
		provider := UpdateCapacityCommandProvider(characterId, worldId, capacity)
		messages, err := provider()
		assert.NoError(t, err)
		assert.Len(t, messages, 1)

		message := messages[0]
		assert.NotEmpty(t, message.Key)

		// Verify the message can be deserialized
		var command list2.Command[list2.UpdateCapacityCommandBody]
		err = json.Unmarshal(message.Value, &command)
		assert.NoError(t, err)
		assert.Equal(t, worldId, command.WorldId)
		assert.Equal(t, characterId, command.CharacterId)
		assert.Equal(t, list2.CommandTypeUpdateCapacity, command.Type)
		assert.Equal(t, capacity, command.Body.Capacity)
	})

	t.Run("different capacity values", func(t *testing.T) {
		testCases := []struct {
			capacity byte
			expected byte
		}{
			{1, 1},
			{20, 20},
			{100, 100},
			{255, 255},
		}

		for _, tc := range testCases {
			provider := UpdateCapacityCommandProvider(characterId, worldId, tc.capacity)
			messages, err := provider()
			assert.NoError(t, err)
			assert.Len(t, messages, 1)

			var command list2.Command[list2.UpdateCapacityCommandBody]
			err = json.Unmarshal(messages[0].Value, &command)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, command.Body.Capacity)
		}
	})

	t.Run("message key generation", func(t *testing.T) {
		provider1 := UpdateCapacityCommandProvider(12345, worldId, capacity)
		provider2 := UpdateCapacityCommandProvider(12345, worldId, capacity)
		provider3 := UpdateCapacityCommandProvider(54321, worldId, capacity)

		messages1, err := provider1()
		assert.NoError(t, err)
		messages2, err := provider2()
		assert.NoError(t, err)
		messages3, err := provider3()
		assert.NoError(t, err)

		// Same character ID should generate same key
		assert.Equal(t, messages1[0].Key, messages2[0].Key)
		// Different character ID should generate different key
		assert.NotEqual(t, messages1[0].Key, messages3[0].Key)
	})
}