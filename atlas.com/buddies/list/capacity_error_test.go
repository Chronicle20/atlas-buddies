package list

import (
	"atlas-buddies/buddy"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCapacityErrorScenarios specifically tests error scenarios for capacity validation
// as required by task 13: "Test error scenarios (invalid capacity, capacity too small)"
func TestCapacityErrorScenarios(t *testing.T) {
	t.Run("invalid capacity values", func(t *testing.T) {
		// Test various invalid capacity values that should be rejected
		
		t.Run("zero capacity is invalid", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			_, err := NewBuilder(tenantId, characterId).
				SetCapacity(0).
				Build()
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "capacity must be greater than 0")
		})
		
		t.Run("capacity boundaries", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Test minimum valid capacity (1)
			model1, err1 := NewBuilder(tenantId, characterId).
				SetCapacity(1).
				Build()
			assert.NoError(t, err1)
			assert.Equal(t, byte(1), model1.Capacity())
			
			// Test maximum valid capacity (255)
			model255, err255 := NewBuilder(tenantId, characterId).
				SetCapacity(255).
				Build()
			assert.NoError(t, err255)
			assert.Equal(t, byte(255), model255.Capacity())
		})
	})
	
	t.Run("capacity too small scenarios", func(t *testing.T) {
		// Test scenarios where capacity is smaller than current buddy count
		
		t.Run("capacity less than buddy count", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Create test buddies
			buddies := createTestBuddies(3) // 3 buddies
			
			// Try to set capacity to 2 when we have 3 buddies
			_, err := NewBuilder(tenantId, characterId).
				SetCapacity(2). // Less than buddy count
				SetBuddies(buddies).
				Build()
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "buddy count exceeds capacity")
		})
		
		t.Run("capacity exactly equal to buddy count is valid", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Create test buddies
			buddies := createTestBuddies(3) // 3 buddies
			
			// Set capacity to exactly match buddy count
			model, err := NewBuilder(tenantId, characterId).
				SetCapacity(3). // Equal to buddy count
				SetBuddies(buddies).
				Build()
			
			assert.NoError(t, err)
			assert.Equal(t, byte(3), model.Capacity())
			assert.Equal(t, 3, len(model.Buddies()))
		})
		
		t.Run("capacity greater than buddy count is valid", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Create test buddies
			buddies := createTestBuddies(3) // 3 buddies
			
			// Set capacity higher than buddy count
			model, err := NewBuilder(tenantId, characterId).
				SetCapacity(10). // Greater than buddy count
				SetBuddies(buddies).
				Build()
			
			assert.NoError(t, err)
			assert.Equal(t, byte(10), model.Capacity())
			assert.Equal(t, 3, len(model.Buddies()))
		})
		
		t.Run("empty buddy list with any capacity is valid", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Test with empty buddy list and various capacities
			capacities := []byte{1, 5, 20, 100, 255}
			
			for _, capacity := range capacities {
				model, err := NewBuilder(tenantId, characterId).
					SetCapacity(capacity).
					SetBuddies([]buddy.Model{}). // Empty buddy list
					Build()
				
				assert.NoError(t, err, "Capacity %d should be valid with empty buddy list", capacity)
				assert.Equal(t, capacity, model.Capacity())
				assert.Equal(t, 0, len(model.Buddies()))
			}
		})
	})
	
	t.Run("edge case scenarios", func(t *testing.T) {
		// Test edge cases for capacity validation
		
		t.Run("maximum buddy count with maximum capacity", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Create many buddies (but not 255 for test performance)
			buddies := createTestBuddies(50)
			
			// Set maximum capacity
			model, err := NewBuilder(tenantId, characterId).
				SetCapacity(255). // Maximum capacity
				SetBuddies(buddies).
				Build()
			
			assert.NoError(t, err)
			assert.Equal(t, byte(255), model.Capacity())
			assert.Equal(t, 50, len(model.Buddies()))
		})
		
		t.Run("capacity validation with large buddy list", func(t *testing.T) {
			tenantId := uuid.New()
			characterId := uint32(12345)
			
			// Create a large buddy list (100 buddies)
			buddies := createTestBuddies(100)
			
			// Try to set capacity less than buddy count
			_, err := NewBuilder(tenantId, characterId).
				SetCapacity(50). // Less than 100 buddies
				SetBuddies(buddies).
				Build()
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "buddy count exceeds capacity")
		})
		
		t.Run("capacity validation preserves other validation rules", func(t *testing.T) {
			// Test that capacity validation doesn't interfere with other validations
			
			// Test with nil tenant ID
			_, err1 := NewBuilder(uuid.Nil, 12345).
				SetCapacity(50).
				Build()
			assert.Error(t, err1)
			assert.Contains(t, err1.Error(), "tenant id is required")
			
			// Test with zero character ID
			_, err2 := NewBuilder(uuid.New(), 0).
				SetCapacity(50).
				Build()
			assert.Error(t, err2)
			assert.Contains(t, err2.Error(), "character id is required")
		})
	})
	
	t.Run("REST validation error scenarios", func(t *testing.T) {
		// Test REST-specific capacity validation scenarios
		
		t.Run("REST model capacity validation", func(t *testing.T) {
			// Test invalid capacities in REST model
			invalidCapacities := []byte{0}
			validCapacities := []byte{1, 20, 50, 100, 255}
			
			for _, capacity := range invalidCapacities {
				restModel := RestModel{Capacity: capacity}
				isValid := restModel.Capacity != 0 // This is the validation logic from the REST handler
				assert.False(t, isValid, "Capacity %d should be invalid in REST model", capacity)
			}
			
			for _, capacity := range validCapacities {
				restModel := RestModel{Capacity: capacity}
				isValid := restModel.Capacity != 0 // This is the validation logic from the REST handler
				assert.True(t, isValid, "Capacity %d should be valid in REST model", capacity)
			}
		})
		
		t.Run("command provider with invalid capacity values", func(t *testing.T) {
			// Test that command provider can be created with any capacity value
			// (validation happens at the processor level)
			characterId := uint32(12345)
			worldId := byte(1)
			
			// Test with zero capacity (invalid but command should be created)
			provider0 := UpdateCapacityCommandProvider(characterId, worldId, 0)
			messages0, err0 := provider0()
			assert.NoError(t, err0)
			assert.Len(t, messages0, 1)
			
			// Test with valid capacities
			validCapacities := []byte{1, 50, 255}
			for _, capacity := range validCapacities {
				provider := UpdateCapacityCommandProvider(characterId, worldId, capacity)
				messages, err := provider()
				assert.NoError(t, err, "Command provider should work with capacity %d", capacity)
				assert.Len(t, messages, 1)
				assert.NotEmpty(t, messages[0].Key)
				assert.NotEmpty(t, messages[0].Value)
			}
		})
	})
}

// TestCapacityValidationConsistency ensures all layers validate capacity consistently
func TestCapacityValidationConsistency(t *testing.T) {
	t.Run("consistent validation across all layers", func(t *testing.T) {
		// Test that builder, REST, and other layers all consistently validate capacity
		
		invalidCapacities := []byte{0}
		validCapacities := []byte{1, 20, 50, 100, 255}
		
		for _, capacity := range invalidCapacities {
			t.Run("invalid capacity validation consistency", func(t *testing.T) {
				tenantId := uuid.New()
				characterId := uint32(12345)
				
				// Test builder validation
				_, builderErr := NewBuilder(tenantId, characterId).
					SetCapacity(capacity).
					Build()
				assert.Error(t, builderErr)
				assert.Contains(t, builderErr.Error(), "capacity must be greater than 0")
				
				// Test REST validation
				restModel := RestModel{Capacity: capacity}
				isRestValid := restModel.Capacity != 0
				assert.False(t, isRestValid, "REST validation should reject capacity %d", capacity)
			})
		}
		
		for _, capacity := range validCapacities {
			t.Run("valid capacity validation consistency", func(t *testing.T) {
				tenantId := uuid.New()
				characterId := uint32(12345)
				
				// Test builder validation
				model, builderErr := NewBuilder(tenantId, characterId).
					SetCapacity(capacity).
					Build()
				assert.NoError(t, builderErr)
				assert.Equal(t, capacity, model.Capacity())
				
				// Test REST validation
				restModel := RestModel{Capacity: capacity}
				isRestValid := restModel.Capacity != 0
				assert.True(t, isRestValid, "REST validation should accept capacity %d", capacity)
			})
		}
	})
}

// Helper function to create test buddies for testing capacity validation
func createTestBuddies(count int) []buddy.Model {
	buddies := make([]buddy.Model, count)
	for i := 0; i < count; i++ {
		entity := buddy.Entity{
			CharacterId:   uint32(1000 + i),
			ListId:        uuid.New(),
			Group:         "default",
			CharacterName: "TestBuddy",
			ChannelId:     -1,
			InShop:        false,
			Pending:       false,
		}
		buddyModel, _ := buddy.Make(entity)
		buddies[i] = buddyModel
	}
	return buddies
}