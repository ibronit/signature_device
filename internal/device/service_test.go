package device

import (
	"log/slog"
	"os"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/enum"
	"github.com/google/uuid"
)

func TestCreateSignatureDevice(t *testing.T) {
	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	deviceRepository := NewDeviceRepository()
	keyPairGeneratorGetter := crypto.NewGeneratorGetter()
	marshalerGetter := crypto.NewMarshalerGetter()
	deviceService := NewDeviceService(deviceRepository, keyPairGeneratorGetter, marshalerGetter, logger)

	testUUID := uuid.New()
	testAlgorithm := enum.RSA
	testLabel := "TestLabel"

	uuid, err := deviceService.CreateSignatureDevice(testUUID, testAlgorithm, testLabel)

	// Assert that there are no errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Assert the uuids match
	if uuid != testUUID {
		t.Errorf("TestUUID and uuid should match: testUUID: %v, uuid: %v", testUUID, uuid)
	}
	// Assert we have 1 saved signature device
	devices := deviceRepository.FindAll()
	if len(devices) != 1 {
		t.Error("There should be exactly 1 devices")
	}
}
