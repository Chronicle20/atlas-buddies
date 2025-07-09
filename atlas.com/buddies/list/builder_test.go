package list

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	tenantId := uuid.New()
	characterId := uint32(12345)

	t.Run("create new buddy list with defaults", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, tenantId, model.TenantId())
		assert.Equal(t, characterId, model.CharacterId())
		assert.Equal(t, byte(20), model.Capacity())
		assert.Empty(t, model.Buddies())
	})

	t.Run("create buddy list with custom capacity", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(50).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(50), model.Capacity())
	})

	t.Run("fail with zero capacity", func(t *testing.T) {
		_, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(0).
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "capacity must be greater than 0")
	})

	t.Run("fail with nil tenant id", func(t *testing.T) {
		_, err := NewBuilder(uuid.Nil, characterId).
			SetId(uuid.New()).
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tenant id is required")
	})

	t.Run("fail with zero character id", func(t *testing.T) {
		_, err := NewBuilder(tenantId, 0).
			SetId(uuid.New()).
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "character id is required")
	})

	t.Run("maximum capacity", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(255).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(255), model.Capacity())
	})

	t.Run("builder from existing model", func(t *testing.T) {
		original, _ := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(30).
			Build()

		modified, err := original.Builder().
			SetCapacity(40).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, original.Id(), modified.Id())
		assert.Equal(t, original.TenantId(), modified.TenantId())
		assert.Equal(t, original.CharacterId(), modified.CharacterId())
		assert.Equal(t, byte(40), modified.Capacity())
	})
}