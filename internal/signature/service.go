package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	fiskalycrypto "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
)

type SignatureService interface {
	SignData(request SignatureRequest) SignatureResponse
	createSignature(msg []byte, privateKey *rsa.PrivateKey) []byte
}

type signatureService struct {
	deviceRepository device.DeviceRepository
	rsaMarshaler     fiskalycrypto.RSAMarshaler
}

func NewSignatureService(deviceRepository device.DeviceRepository, rsaMarshaler fiskalycrypto.RSAMarshaler) SignatureService {
	return &signatureService{deviceRepository: deviceRepository, rsaMarshaler: rsaMarshaler}
}

func (ss *signatureService) SignData(request SignatureRequest) SignatureResponse {
	device, err := ss.deviceRepository.FindById(request.DeviceId)
	if err != nil {
		// TODO: Log
	}

	rsaKeyPair, err := ss.rsaMarshaler.Unmarshal(device.PrivateKey)
	if err != nil {
		// TODO: Log
	}

	dataByte := []byte(request.DataToBeSigned)
	signature := ss.createSignature(dataByte, rsaKeyPair.Private)
	ss.deviceRepository.Update(device)

	response := SignatureResponse{Signature: base64.StdEncoding.EncodeToString(signature)}

	return response
}

func (ss *signatureService) createSignature(msg []byte, privateKey *rsa.PrivateKey) []byte {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		panic(err) // TODO: error handling
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	return signature
}
