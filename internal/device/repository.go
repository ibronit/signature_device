package device

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type DeviceRepository interface {
	Save(entity DeviceEntity)
	Update(entity DeviceEntity)
	FindById(id uuid.UUID) (DeviceEntity, error)
}

type deviceRepository struct {
	sync.RWMutex
	entities map[uuid.UUID]DeviceEntity
}

func NewDeviceRepository() DeviceRepository {
	return &deviceRepository{entities: make(map[uuid.UUID]DeviceEntity)}
}

func (repo *deviceRepository) Save(entity DeviceEntity) {
	repo.Lock()
	defer repo.Unlock()
	_, exists := repo.entities[entity.Uuid]
	if exists {
		slog.Info("Entity already exists!", "id", entity.Uuid)
		return
	}

	entity.Counter = 1
	entity.LastSignature = time.Now()
	repo.entities[entity.Uuid] = entity
}

func (repo *deviceRepository) Update(entity DeviceEntity) {
	repo.Lock()
	defer repo.Unlock()
	existingEntity, exists := repo.entities[entity.Uuid]
	if !exists {
		slog.Error("Entity doesn't exist! Entity cannot be updated!", "id", entity.Uuid)
		return
	}

	existingEntity.Counter++
	existingEntity.LastSignature = time.Now()
	repo.entities[entity.Uuid] = existingEntity
}

func (repo *deviceRepository) FindById(id uuid.UUID) (DeviceEntity, error) {
	existingEntity, exists := repo.entities[id]
	if !exists {
		slog.Error("Entity doesn't exist!", "id", id)
		return DeviceEntity{}, errors.New("Entity doesn't exist!")
	}

	return existingEntity, nil
}
