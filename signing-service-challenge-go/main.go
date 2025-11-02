package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	repo := persistence.NewInMemorySignatureDeviceRepository()
	signatureService := service.NewSignatureService(repo)
	deviceHandler := api.NewDeviceHandler(signatureService)

	server := api.NewServer(ListenAddress, deviceHandler)

	log.Printf("Signature service listening on %s", ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatalf("could not start server on %s: %v", ListenAddress, err)
	}
}
