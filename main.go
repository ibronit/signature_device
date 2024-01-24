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
)

func main() {
	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	deviceRepository := device.NewDeviceRepository()
	keyPairGeneratorGetter := crypto.NewGeneratorGetter()
	marshalerGetter := crypto.NewMarshalerGetter()
	signatureGetter := crypto.NewSignerGetter()
	deviceService := device.NewDeviceService(deviceRepository, keyPairGeneratorGetter, marshalerGetter, logger)
	signatureService := signature.NewSignatureService(deviceRepository, marshalerGetter, signatureGetter, logger)

	mux := http.NewServeMux()

	healthHandler := health.HealthHandler{}
	deviceHandler := device.NewDeviceHandler(deviceService, logger)
	signatureHandler := signature.NewSignatureHandler(signatureService, logger)
	mux.Handle("/api/v0/health", &healthHandler)
	mux.Handle("/api/v1/device", deviceHandler)
	mux.Handle("/api/v1/signature", signatureHandler)

	logger.Info("Server is listening on port", "port", ListenAddress)
	err := http.ListenAndServe(ListenAddress, mux)
	if err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
