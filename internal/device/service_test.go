package device

import (
	// "testing"

	"github.com/google/uuid"
)

type DeviceRepositoryMock struct {
	SaveInvoked int
	UpdateInvoked int
	FindByIdInvoked int
}

func (drm *DeviceRepositoryMock) Save(entity DeviceEntity) {
	drm.SaveInvoked++
}

func (drm *DeviceRepositoryMock) Update(entity DeviceEntity) {
	drm.UpdateInvoked++
}

func (drm *DeviceRepositoryMock) FindById(id uuid.UUID) (DeviceEntity, error) {
	drm.FindByIdInvoked++

	return DeviceEntity{}, nil
}

// func TestCreateSignatureDevice(t *testing.T) {
// 	deviceRepositoryMock := DeviceRepositoryMock
// 	deviceService := NewDeviceService()

// 	t.Fail()
// }
