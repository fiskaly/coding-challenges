package main

import (
	"log"
	"signing-service-challenge/api"
	"signing-service-challenge/crypto"
	"signing-service-challenge/persistence"
	"signing-service-challenge/service"
)

const (
	ListenAddress = ":8081"
	// TODO: add further configuration parameters here ...
)

func main() {

	//initialize Inmemory
	inmemoryDeviceRepository := persistence.NewInmemoryDeviceRepository()
	// initialize Generators
	eccGenerator := &crypto.DefaultECCGenerator{}
	rsaGenerator := &crypto.DefaultRSAGenerator{}

	// initialize Marshalers
	rsaMarshaler := &crypto.DefaultRSAMarshaler{}
	ecdsaMarshaler := &crypto.DefaultECCMarshaler{}

	//Initialize Services
	deviceService := service.NewDefaultDeviceService(inmemoryDeviceRepository,
		eccGenerator,
		rsaGenerator,
		ecdsaMarshaler,
		rsaMarshaler)
	transactionService := service.NewDefaultTransactionService(inmemoryDeviceRepository)

	//Initialize Handlers
	deviceHandler := api.NewDeviceHandler(deviceService)
	transactionHandler := api.NewTransactionHandler(transactionService)

	//Initialize Server
	server := api.NewServer(ListenAddress, deviceHandler, transactionHandler)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
