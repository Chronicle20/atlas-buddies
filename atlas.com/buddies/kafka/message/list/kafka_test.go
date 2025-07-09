package list

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCapacityCommand(t *testing.T) {
	t.Run("create update capacity command", func(t *testing.T) {
		command := Command[UpdateCapacityCommandBody]{
			WorldId:     1,
			CharacterId: 12345,
			Type:        CommandTypeUpdateCapacity,
			Body: UpdateCapacityCommandBody{
				Capacity: 50,
			},
		}

		// Test serialization
		data, err := json.Marshal(command)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)

		// Test deserialization
		var decoded Command[UpdateCapacityCommandBody]
		err = json.Unmarshal(data, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, command.WorldId, decoded.WorldId)
		assert.Equal(t, command.CharacterId, decoded.CharacterId)
		assert.Equal(t, command.Type, decoded.Type)
		assert.Equal(t, command.Body.Capacity, decoded.Body.Capacity)
	})

	t.Run("command type constant", func(t *testing.T) {
		assert.Equal(t, "UPDATE_CAPACITY", CommandTypeUpdateCapacity)
	})

	t.Run("update capacity command body", func(t *testing.T) {
		body := UpdateCapacityCommandBody{
			Capacity: 100,
		}

		data, err := json.Marshal(body)
		assert.NoError(t, err)

		var decoded UpdateCapacityCommandBody
		err = json.Unmarshal(data, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, body.Capacity, decoded.Capacity)
	})
}