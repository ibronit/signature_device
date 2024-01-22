package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type DeviceService struct {
	repository persistence.Repository
}

func NewDeviceService(repository persistence.Repository) *DeviceService {
	return &DeviceService{repository: repository}
}

func (ds *DeviceService) CreateSignatureDevice(device api.DeviceRequest) {
	ds.repository.Save(device)
}

func (ds *DeviceService) SignData(device api.DeviceRequest) {
	ds.repository.Update(device)
}
