package main

import (
	"encoding/json"
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
	rsaGenerator := crypto.RSAGenerator{}
	rsaMarshaler := crypto.NewRSAMarshaler()
	deviceService := device.NewDeviceService(deviceRepository, rsaGenerator, rsaMarshaler, logger)
	signature.NewSignatureService(deviceRepository, rsaMarshaler)
	deviceJson := `{"id": "bc890106-641e-41fc-aed8-b77cca0b42b9", "algorithm": "ECC", "label": "DELETED"}`
	var deviceRequest device.DeviceRequest
	if err := json.Unmarshal([]byte(deviceJson), &deviceRequest); err != nil {
		panic(err)
	}
	// deviceService.CreateSignatureDevice(deviceRequest)
	// Signature
	// signatureJson := `{"device_id": "bc890106-641e-41fc-aed8-b77cca0b42b9", "data_to_be_signed": "valami"}`
	// var signatureRequest signature.SignatureRequest
	// if err := json.Unmarshal([]byte(signatureJson), &signatureRequest); err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < 100; i++ {
	// 	go signatureService.SignData(signatureRequest) // TODO: Give simpler params here, we don't need the whole request struct
	// }
	// time.Sleep(1 * time.Second)
	// fmt.Println(deviceRepository.FindById(deviceRequest.Id))

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
