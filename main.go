package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/device"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/health"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/signature"
)

const (
	ListenAddress = ":8000"
	// TODO: add further configuration parameters here ...
)

func main() {
	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	deviceRepository := device.NewDeviceRepository()
	keyPairGeneratorGetter := crypto.GeneratorGetter{}
	marshalerGetter := crypto.MarshalerGetter{}
	signatureGetter := crypto.SignerGetter{}
	deviceService := device.NewDeviceService(deviceRepository, keyPairGeneratorGetter, marshalerGetter, logger)
	signature.NewSignatureService(deviceRepository, marshalerGetter, signatureGetter, logger)

	mux := http.NewServeMux()

	healthHandler := health.HealthHandler{}
	deviceHandler := device.NewDeviceHandler(deviceService, logger)
	mux.Handle("/api/v0/health", &healthHandler)
	mux.Handle("/api/v1/device", deviceHandler)

	err := http.ListenAndServe(ListenAddress, mux)
	if err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
