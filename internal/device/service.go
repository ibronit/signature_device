package device

import (
	"log/slog"

	fiskalycrypto "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/google/uuid"
)

type DeviceService interface {
	CreateSignatureDevice(uuid uuid.UUID, algorithm Algorithm, label string) (uuid.UUID, error)
}

type deviceService struct {
	repository   DeviceRepository
	rsaGenerator fiskalycrypto.RSAGenerator
	rsaMarshaler fiskalycrypto.RSAMarshaler
	logger       *slog.Logger
}

func NewDeviceService(
	repository DeviceRepository,
	rsaGenerator fiskalycrypto.RSAGenerator,
	rsaMarshaler fiskalycrypto.RSAMarshaler,
	logger *slog.Logger) DeviceService {
	return &deviceService{
		repository:   repository,
		rsaGenerator: rsaGenerator,
		rsaMarshaler: rsaMarshaler,
		logger:       logger,
	}
}

func (ds *deviceService) CreateSignatureDevice(uuid uuid.UUID, algorithm Algorithm, label string) (uuid.UUID, error) {
	entity := DeviceEntity{
		Uuid:      uuid,
		Algorithm: algorithm,
		Label:     label,
	}

	keyPair, err := ds.rsaGenerator.Generate()
	if err != nil {
		ds.logger.Error("Couldn't generate keypair", "error", err)
		return uuid, err
	}
	publicKey, privateKey, err := ds.rsaMarshaler.Marshal(*keyPair)
	if err != nil {
		ds.logger.Error("Couldn't marshal keypair", "error", err)
		return uuid, err
	}
	entity.PrivateKey = privateKey
	entity.PublicKey = publicKey

	ds.repository.Save(entity)
	return uuid, nil
}
