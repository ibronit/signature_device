package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence/device"
)

const (
	ListenAddress = ":8000"
	// TODO: add further configuration parameters here ...
)

func main() {
	// playground
	// TODO: Move to handler logic
	inMemory := persistence.NewInMemory()
	deviceService := domain.NewDeviceService(inMemory)
	deviceJson := `{"id": "bc890106-641e-41fc-aed8-b77cca0b42b9", "algorithm": "ECC", "label": "DELETED"}`
	var device api.DeviceRequest
	if err := json.Unmarshal([]byte(deviceJson), &device); err != nil {
		panic(err)
	}
	go deviceService.CreateSignatureDevice(device)
	for i := 0; i < 100; i++ {
		go deviceService.SignData(device)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(inMemory.FindAll())

	// do not change
	server := api.NewServer(ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
