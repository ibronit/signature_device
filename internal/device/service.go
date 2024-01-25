package device

import (
	"log/slog"

	fiskalycrypto "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
	"github.com/google/uuid"
)

type DeviceService interface {
	CreateSignatureDevice(uuid uuid.UUID, algorithm enum.Algorithm, label string) (uuid.UUID, error)
	GetAllDevices() map[uuid.UUID]DeviceEntity
}

type deviceService struct {
	repository      DeviceRepository
	generatorGetter *fiskalycrypto.GeneratorGetter
	marshalerGetter *fiskalycrypto.MarshalerGetter
	logger          *slog.Logger
}

func NewDeviceService(
	repository DeviceRepository,
	generatorGetter *fiskalycrypto.GeneratorGetter,
	marshalerGetter *fiskalycrypto.MarshalerGetter,
	logger *slog.Logger) DeviceService {
	return &deviceService{
		repository:      repository,
		generatorGetter: generatorGetter,
		marshalerGetter: marshalerGetter,
		logger:          logger,
	}
}

// CreateSignatureDevice creates signature device and saves the signature device into memory
func (ds *deviceService) CreateSignatureDevice(uuid uuid.UUID, algorithm enum.Algorithm, label string) (uuid.UUID, error) {
	entity := DeviceEntity{
		Uuid:      uuid,
		Algorithm: algorithm,
		Label:     label,
	}

	publicKey, privateKey, err := ds.generateKeyPairByAlgorithm(algorithm)
	if err != nil {
		ds.logger.Error("Couldn't generate keypair with algortihm", "error", err, "algorithm", algorithm)
		return uuid, err
	}
	entity.PrivateKey = privateKey
	entity.PublicKey = publicKey

	ds.repository.Save(entity)
	return uuid, nil
}

func (ds *deviceService) generateKeyPairByAlgorithm(algorithm enum.Algorithm) ([]byte, []byte, error) {
	generator, err := ds.generatorGetter.GetGeneratorByAlgorithm(algorithm)
	if err != nil {
		ds.logger.Error("Couldn't get the right key-pair generator", "error", err)
	}

	keyPair, err := generator.Generate()
	if err != nil {
		ds.logger.Error("Couldn't generate keypair", "error", err)
		return nil, nil, err
	}

	marshaler, err := ds.marshalerGetter.GetMarshalerByAlgorithm(algorithm)
	if err != nil {
		ds.logger.Error("Couldn't get the right marshaler", "error", err)
	}

	publicKey, privateKey, err := marshaler.Marshal(keyPair)
	if err != nil {
		ds.logger.Error("Couldn't marshal keypair", "error", err)
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

// GetAllDevices gets every previously saved signature devices.
func (ds *deviceService) GetAllDevices() map[uuid.UUID]DeviceEntity {
	return ds.repository.FindAll()
}
