package list

import (
	"atlas-buddies/buddy"
	"errors"
	"fmt"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func create(db *gorm.DB, t tenant.Model, characterId uint32, capacity byte) (Model, error) {
	e := &Entity{
		TenantId:    t.Id(),
		CharacterId: characterId,
		Capacity:    capacity,
	}

	err := db.Create(e).Error
	if err != nil {
		return Model{}, err
	}
	return Make(*e)
}

func addPendingBuddy(db *gorm.DB, tenantId uuid.UUID, characterId uint32, targetId uint32, targetName string, group string) error {
	return addBuddy(db, tenantId, characterId, targetId, targetName, group, true)
}

func addBuddy(db *gorm.DB, tenantId uuid.UUID, characterId uint32, targetId uint32, targetName string, group string, pending bool) error {
	e, err := byCharacterIdEntityProvider(tenantId, characterId)(db)()
	if err != nil {
		return err
	}

	nb := buddy.Entity{
		CharacterId:   targetId,
		ListId:        e.Id,
		Group:         group,
		CharacterName: targetName,
		ChannelId:     -1,
		Pending:       pending,
	}
	return db.Create(&nb).Error
}

func removeBuddy(db *gorm.DB, tenantId uuid.UUID, characterId uint32, targetId uint32) error {
	e, err := byCharacterIdEntityProvider(tenantId, characterId)(db)()
	if err != nil {
		return err
	}

	var rb buddy.Entity
	for _, b := range e.Buddies {
		if b.CharacterId == targetId {
			rb = b
			break
		}
	}

	if rb.ListId == uuid.Nil {
		return gorm.ErrRecordNotFound
	}

	return db.Delete(&rb).Error
}

func updateBuddyChannel(db *gorm.DB, tenantId uuid.UUID, characterId uint32, targetId uint32, channelId int8) (bool, error) {
	bbl, err := byCharacterIdEntityProvider(tenantId, targetId)(db)()
	if err != nil {
		return false, err
	}

	var meAsBuddy *buddy.Entity
	for _, pm := range bbl.Buddies {
		if pm.CharacterId == characterId {
			meAsBuddy = &pm
		}
	}
	if meAsBuddy == nil {
		return false, nil
	}
	meAsBuddy.ChannelId = channelId

	err = db.Save(meAsBuddy).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func updateBuddyShopStatus(db *gorm.DB, tenantId uuid.UUID, characterId uint32, targetId uint32, inShop bool) (bool, error) {
	bbl, err := byCharacterIdEntityProvider(tenantId, targetId)(db)()
	if err != nil {
		return false, err
	}

	var meAsBuddy *buddy.Entity
	for _, pm := range bbl.Buddies {
		if pm.CharacterId == characterId {
			meAsBuddy = &pm
		}
	}
	if meAsBuddy == nil {
		return false, nil
	}
	meAsBuddy.InShop = inShop

	err = db.Save(meAsBuddy).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func deleteEntityWithBuddies(db *gorm.DB, tenantId uuid.UUID, characterId uint32) error {
	var entity Entity

	// Step 1: Find the Entity
	if err := db.
		Where("tenant_id = ? AND character_id = ?", tenantId, characterId).
		First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // No-op if not found
		}
		return fmt.Errorf("failed to find entity: %w", err)
	}

	// Step 2: Delete associated Buddies
	if err := db.
		Where("list_id = ?", entity.Id).
		Delete(&buddy.Entity{}).Error; err != nil {
		return fmt.Errorf("failed to delete buddies: %w", err)
	}

	// Step 3: Delete the Entity
	if err := db.Delete(&entity).Error; err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}

	return nil
}

func updateCapacity(db *gorm.DB) func(tenantId uuid.UUID) func(characterId uint32) func(capacity byte) error {
	return func(tenantId uuid.UUID) func(characterId uint32) func(capacity byte) error {
		return func(characterId uint32) func(capacity byte) error {
			return func(capacity byte) error {
				if capacity == 0 {
					return errors.New("capacity must be greater than 0")
				}

				// Find the entity
				e, err := byCharacterIdEntityProvider(tenantId, characterId)(db)()
				if err != nil {
					return err
				}

				// Check if new capacity is less than current buddy count
				currentBuddyCount := len(e.Buddies)
				if int(capacity) < currentBuddyCount {
					return fmt.Errorf("new capacity %d is less than current buddy count %d", capacity, currentBuddyCount)
				}

				// Update capacity
				e.Capacity = capacity
				return db.Model(&Entity{}).Where("id = ?", e.Id).Update("capacity", capacity).Error
			}
		}
	}
}
