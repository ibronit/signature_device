package signature

import (
	"log/slog"
	"os"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
	"github.com/google/uuid"
)

func TestSignData(t *testing.T) {
	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	deviceRepository := device.NewDeviceRepository()
	keyPairGeneratorGetter := crypto.NewGeneratorGetter()
	marshalerGetter := crypto.NewMarshalerGetter()
	signatureGetter := crypto.NewSignerGetter()
	deviceService := device.NewDeviceService(deviceRepository, keyPairGeneratorGetter, marshalerGetter, logger)
	signatureService := NewSignatureService(deviceRepository, marshalerGetter, signatureGetter, logger)

	// Create test signature device
	testUUID := uuid.New()
	testAlgorithm := enum.ECC
	testLabel := "TestLabel"

	_, err := deviceService.CreateSignatureDevice(testUUID, testAlgorithm, testLabel)
	// Assert that there are no errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	testDataToBeSign := "Test data"
	_, signature, lastSignature, err := signatureService.SignData(testUUID, testDataToBeSign)
	// Assert that there are no errors
	if err != nil {
		t.Errorf("Unexpected error during signature creation: %v", err)
	}
	// Assert that signature is there
	if len(signature) == 0 {
		t.Errorf("Signature shouldn't be empty!")
	}
	// Assert that last signature is there
	if len(lastSignature) == 0 {
		t.Errorf("Signature shouldn't be empty!")
	}
}
