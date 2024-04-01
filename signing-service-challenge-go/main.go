package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence/inmemory"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {

	repo := inmemory.NewInMemoryStore()
	service := service.NewSignatureDeviceService(repo)
	deviceHandler := api.NewDeviceHandler(service)

	server := api.NewServer(ListenAddress, deviceHandler)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
