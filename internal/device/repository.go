package device

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

// DeviceRepository performs database operations on DeviceEntities.
type DeviceRepository interface {
	Save(entity DeviceEntity)
	Update(id uuid.UUID, lastSignature string) (DeviceEntity, error)
	FindById(id uuid.UUID) (DeviceEntity, error)
	FindAll() map[uuid.UUID]DeviceEntity
}

type deviceRepository struct {
	sync.RWMutex
	entities map[uuid.UUID]DeviceEntity
}

func NewDeviceRepository() DeviceRepository {
	return &deviceRepository{entities: make(map[uuid.UUID]DeviceEntity)}
}

// Saves a new DeviceEntity.
func (repo *deviceRepository) Save(entity DeviceEntity) {
	repo.Lock()
	defer repo.Unlock()
	_, exists := repo.entities[entity.Uuid]
	if exists {
		slog.Info("Entity already exists!", "id", entity.Uuid)
		return
	}

	entity.Counter = 0
	entity.LastSignature = ""
	repo.entities[entity.Uuid] = entity
}

// Updates an existing DeviceEntity and bumps the Counter by one.
func (repo *deviceRepository) Update(id uuid.UUID, lastSignature string) (DeviceEntity, error) {
	repo.Lock()
	defer repo.Unlock()
	existingEntity, exists := repo.entities[id]
	if !exists {
		slog.Error("Entity doesn't exist! Entity cannot be updated!", "id", id)
		return existingEntity, errors.New("Entity doesn't exist")
	}

	existingEntity.Counter++
	existingEntity.LastSignature = lastSignature
	repo.entities[id] = existingEntity

	return existingEntity, nil
}

// Retrieves a DeviceEntity if it exists
func (repo *deviceRepository) FindById(id uuid.UUID) (DeviceEntity, error) {
	existingEntity, exists := repo.entities[id]
	if !exists {
		slog.Error("Entity doesn't exist!", "id", id)
		return DeviceEntity{}, errors.New("Entity doesn't exist!")
	}

	return existingEntity, nil
}

// Returns all saved DeviceEntities
func (repo *deviceRepository) FindAll() map[uuid.UUID]DeviceEntity {
	return repo.entities
}
