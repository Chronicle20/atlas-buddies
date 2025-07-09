package list

import (
	"atlas-buddies/buddy"
	"errors"
	"github.com/google/uuid"
)

type Builder struct {
	tenantId    uuid.UUID
	id          uuid.UUID
	characterId uint32
	capacity    byte
	buddies     []buddy.Model
}

func NewBuilder(tenantId uuid.UUID, characterId uint32) *Builder {
	return &Builder{
		tenantId:    tenantId,
		characterId: characterId,
		capacity:    20, // default capacity
		buddies:     []buddy.Model{},
	}
}

func (b *Builder) SetId(id uuid.UUID) *Builder {
	b.id = id
	return b
}

func (b *Builder) SetCapacity(capacity byte) *Builder {
	b.capacity = capacity
	return b
}

func (b *Builder) SetBuddies(buddies []buddy.Model) *Builder {
	b.buddies = buddies
	return b
}

func (b *Builder) Build() (Model, error) {
	if err := b.validate(); err != nil {
		return Model{}, err
	}

	return Model{
		tenantId:    b.tenantId,
		id:          b.id,
		characterId: b.characterId,
		capacity:    b.capacity,
		buddies:     b.buddies,
	}, nil
}

func (b *Builder) validate() error {
	if b.tenantId == uuid.Nil {
		return errors.New("tenant id is required")
	}

	if b.characterId == 0 {
		return errors.New("character id is required")
	}

	if b.capacity == 0 {
		return errors.New("capacity must be greater than 0")
	}

	if len(b.buddies) > int(b.capacity) {
		return errors.New("buddy count exceeds capacity")
	}

	return nil
}

func (m Model) Builder() *Builder {
	return &Builder{
		tenantId:    m.tenantId,
		id:          m.id,
		characterId: m.characterId,
		capacity:    m.capacity,
		buddies:     m.buddies,
	}
}