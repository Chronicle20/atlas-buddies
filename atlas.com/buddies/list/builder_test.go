package list

import (
	"atlas-buddies/buddy"
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

func TestBuilderValidation(t *testing.T) {
	tenantId := uuid.New()
	characterId := uint32(12345)

	// Helper function to create test buddy models
	createBuddy := func(id uint32, name string) buddy.Model {
		entity := buddy.Entity{
			CharacterId:   id,
			ListId:        uuid.New(),
			Group:         "default",
			CharacterName: name,
			ChannelId:     -1,
			InShop:        false,
			Pending:       false,
		}
		model, _ := buddy.Make(entity)
		return model
	}

	t.Run("validation passes with valid parameters", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(20).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, tenantId, model.TenantId())
		assert.Equal(t, characterId, model.CharacterId())
		assert.Equal(t, byte(20), model.Capacity())
	})

	t.Run("validation fails with nil tenant id", func(t *testing.T) {
		_, err := NewBuilder(uuid.Nil, characterId).
			SetId(uuid.New()).
			SetCapacity(20).
			Build()

		assert.Error(t, err)
		assert.Equal(t, "tenant id is required", err.Error())
	})

	t.Run("validation fails with zero character id", func(t *testing.T) {
		_, err := NewBuilder(tenantId, 0).
			SetId(uuid.New()).
			SetCapacity(20).
			Build()

		assert.Error(t, err)
		assert.Equal(t, "character id is required", err.Error())
	})

	t.Run("validation fails with zero capacity", func(t *testing.T) {
		_, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(0).
			Build()

		assert.Error(t, err)
		assert.Equal(t, "capacity must be greater than 0", err.Error())
	})

	t.Run("validation fails when buddy count exceeds capacity", func(t *testing.T) {
		// Create 3 buddies but set capacity to 2
		buddies := []buddy.Model{
			createBuddy(1001, "Friend1"),
			createBuddy(1002, "Friend2"),
			createBuddy(1003, "Friend3"),
		}

		_, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(2).
			SetBuddies(buddies).
			Build()

		assert.Error(t, err)
		assert.Equal(t, "buddy count exceeds capacity", err.Error())
	})

	t.Run("validation passes when buddy count equals capacity", func(t *testing.T) {
		// Create 2 buddies with capacity of 2
		buddies := []buddy.Model{
			createBuddy(1001, "Friend1"),
			createBuddy(1002, "Friend2"),
		}

		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(2).
			SetBuddies(buddies).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(2), model.Capacity())
		assert.Len(t, model.Buddies(), 2)
	})

	t.Run("validation passes when buddy count is less than capacity", func(t *testing.T) {
		// Create 1 buddy with capacity of 5
		buddies := []buddy.Model{
			createBuddy(1001, "Friend1"),
		}

		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(5).
			SetBuddies(buddies).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(5), model.Capacity())
		assert.Len(t, model.Buddies(), 1)
	})

	t.Run("validation passes with empty buddies list", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(10).
			SetBuddies([]buddy.Model{}).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(10), model.Capacity())
		assert.Empty(t, model.Buddies())
	})

	t.Run("validation with maximum capacity 255", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(255).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(255), model.Capacity())
	})

	t.Run("validation with minimum capacity 1", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(uuid.New()).
			SetCapacity(1).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, byte(1), model.Capacity())
	})
}

func TestBuilderDefaults(t *testing.T) {
	tenantId := uuid.New()
	characterId := uint32(12345)

	t.Run("new builder sets default capacity to 20", func(t *testing.T) {
		builder := NewBuilder(tenantId, characterId)
		assert.Equal(t, byte(20), builder.capacity)
	})

	t.Run("new builder sets empty buddies list", func(t *testing.T) {
		builder := NewBuilder(tenantId, characterId)
		assert.Empty(t, builder.buddies)
	})

	t.Run("new builder sets provided tenant id", func(t *testing.T) {
		builder := NewBuilder(tenantId, characterId)
		assert.Equal(t, tenantId, builder.tenantId)
	})

	t.Run("new builder sets provided character id", func(t *testing.T) {
		builder := NewBuilder(tenantId, characterId)
		assert.Equal(t, characterId, builder.characterId)
	})
}

func TestBuilderChaining(t *testing.T) {
	tenantId := uuid.New()
	characterId := uint32(12345)
	id := uuid.New()

	t.Run("all setter methods return builder for chaining", func(t *testing.T) {
		builder := NewBuilder(tenantId, characterId)
		
		// Test that all setters return *Builder
		result1 := builder.SetId(id)
		assert.Equal(t, builder, result1)
		
		result2 := builder.SetCapacity(50)
		assert.Equal(t, builder, result2)
		
		result3 := builder.SetBuddies([]buddy.Model{})
		assert.Equal(t, builder, result3)
	})

	t.Run("fluent interface works correctly", func(t *testing.T) {
		model, err := NewBuilder(tenantId, characterId).
			SetId(id).
			SetCapacity(100).
			SetBuddies([]buddy.Model{}).
			Build()

		assert.NoError(t, err)
		assert.Equal(t, id, model.Id())
		assert.Equal(t, byte(100), model.Capacity())
		assert.Empty(t, model.Buddies())
	})
}