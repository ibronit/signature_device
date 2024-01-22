package persistence

import (
	"log/slog"
	"sync"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/google/uuid"
)

type InMemory struct {
	sync.RWMutex
	entities map[uuid.UUID]api.DeviceRequest
}

func NewInMemory() *InMemory {
	return &InMemory{entities: make(map[uuid.UUID]api.DeviceRequest)}
}

func (im *InMemory) Save(entity api.DeviceRequest) {
	im.Lock()
	defer im.Unlock()
	_, exists := im.entities[entity.Id]
	if exists {
		slog.Info("Entity already exists!", "id", entity.Id)
		return
	}

	entity.Counter = 1
	entity.LastSignature = time.Now()
	im.entities[entity.Id] = entity
}

func (im *InMemory) Update(entity api.DeviceRequest) {
	im.Lock()
	defer im.Unlock()
	existingEntity, exists := im.entities[entity.Id]
	if !exists {
		slog.Error("Entity doesn't exist! Entity cannot be updated!", "id", entity.Id)
		return
	}

	existingEntity.Counter++
	existingEntity.LastSignature = time.Now()
	im.entities[entity.Id] = existingEntity
}

func (im *InMemory) FindAll() map[uuid.UUID]api.DeviceRequest {
	return im.entities
}
