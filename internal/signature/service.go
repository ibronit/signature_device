package signature

import (
	"encoding/base64"
	"log/slog"

	fiskalycrypto "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
	"github.com/google/uuid"
)

// SignatureService covers the business logic of signatures
type SignatureService interface {
	SignData(deviceId uuid.UUID, dataToBeSigned string) (device.DeviceEntity, string, string, error)
}

type signatureService struct {
	deviceRepository device.DeviceRepository
	marshalerGetter  *fiskalycrypto.MarshalerGetter
	signerGetter     *fiskalycrypto.SignerGetter
	logger           *slog.Logger
}

func NewSignatureService(
	deviceRepository device.DeviceRepository,
	marshalerGetter *fiskalycrypto.MarshalerGetter,
	signerGetter *fiskalycrypto.SignerGetter,
	logger *slog.Logger) SignatureService {
	return &signatureService{deviceRepository: deviceRepository, marshalerGetter: marshalerGetter, signerGetter: signerGetter, logger: logger}
}

// SignData signs data with the provided signature device.
func (ss *signatureService) SignData(deviceId uuid.UUID, dataToBeSigned string) (device.DeviceEntity, string, string, error) {
	device, err := ss.deviceRepository.FindById(deviceId)
	if err != nil {
		ss.logger.Error("Couldn't find the signature device by uuid", "uuid", deviceId)
		return device, "", "", err
	}

	marshaler, err := ss.marshalerGetter.GetMarshalerByAlgorithm(device.Algorithm)
	if err != nil {
		ss.logger.Error("Couldn't get the right marshaler", "error", err)
		return device, "", "", err
	}

	keyPair, err := marshaler.Unmarshal(device.PrivateKey)
	if err != nil {
		ss.logger.Error("Couldn't unmarshal the key-pair", "error", err)
		return device, "", "", err
	}

	signer, err := ss.signerGetter.GetSignerByAlgorithm(device.Algorithm)
	if err != nil {
		ss.logger.Error("Couldn't get the right marshaler", "error", err)
		return device, "", "", err
	}

	dataByte := []byte(dataToBeSigned)
	signature, err := signer.CreateSignature(dataByte, keyPair)
	if err != nil {
		ss.logger.Error("Couldn't create the signature", "error", err)
		return device, "", "", err
	}
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	updatedDevice, err := ss.deviceRepository.Update(deviceId, base64Signature)
	if err != nil {
		ss.logger.Error("Couldn't update the signature device", "error", err)
		return device, "", "", err
	}
	var lastSignature string
	if device.Counter == 0 {
		binaryUuid, err := deviceId.MarshalBinary()
		if err != nil {
			ss.logger.Error("Couldn't marshal uuid to binary", "error", err)
			return device, "", "", err
		}
		lastSignature = base64.StdEncoding.EncodeToString(binaryUuid)
	} else {
		lastSignature = device.LastSignature
	}

	return updatedDevice, lastSignature, base64Signature, nil
}
