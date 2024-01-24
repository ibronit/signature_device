package signature

import (
	"encoding/base64"
	"log/slog"

	fiskalycrypto "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
)

type SignatureService interface {
	SignData(request SignatureRequest) (SignatureResponse, error)
}

type signatureService struct {
	deviceRepository device.DeviceRepository
	marshalerGetter  fiskalycrypto.MarshalerGetter
	signerGetter     fiskalycrypto.SignerGetter
	logger           *slog.Logger
}

func NewSignatureService(
	deviceRepository device.DeviceRepository,
	marshalerGetter fiskalycrypto.MarshalerGetter,
	signerGetter fiskalycrypto.SignerGetter,
	logger *slog.Logger) SignatureService {
	return &signatureService{deviceRepository: deviceRepository, marshalerGetter: marshalerGetter, signerGetter: signerGetter, logger: logger}
}

func (ss *signatureService) SignData(request SignatureRequest) (SignatureResponse, error) {
	device, err := ss.deviceRepository.FindById(request.DeviceId)
	if err != nil {
		ss.logger.Error("Couldn't find the signature device by uuid", "uuid", request.DeviceId)
		return SignatureResponse{}, err
	}

	marshaler, err := ss.marshalerGetter.GetMarshalerByAlgorithm(fiskalycrypto.Algorithm(device.Algorithm))
	if err != nil {
		ss.logger.Error("Couldn't get the right marshaler", "error", err)
		return SignatureResponse{}, err
	}

	keyPair, err := marshaler.Unmarshal(device.PrivateKey)
	if err != nil {
		ss.logger.Error("Couldn't unmarshal the key-pair", "error", err)
		return SignatureResponse{}, err
	}

	signer, err := ss.signerGetter.GetSignatureByAlgorithm(fiskalycrypto.Algorithm(device.Algorithm))
	if err != nil {
		ss.logger.Error("Couldn't get the right marshaler", "error", err)
		return SignatureResponse{}, err
	}

	dataByte := []byte(request.DataToBeSigned)
	signature, err := signer.CreateSignature(dataByte, keyPair)
	if err != nil {
		ss.logger.Error("Couldn't create the signature", "error", err)
		return SignatureResponse{}, err
	}
	ss.deviceRepository.Update(device)

	response := SignatureResponse{Signature: base64.StdEncoding.EncodeToString(signature)}

	return response, nil
}

// func (ss *signatureService) createSignature(msg []byte, privateKey *rsa.PrivateKey) []byte {
// 	msgHash := sha256.New()
// 	_, err := msgHash.Write(msg)
// 	if err != nil {
// 		panic(err) // TODO: error handling
// 	}
// 	msgHashSum := msgHash.Sum(nil)

// 	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return signature
// }
