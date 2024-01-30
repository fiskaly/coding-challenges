package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/application"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"

	//"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

const (
	ListenAddress = ":8081"
	// TODO: add further configuration parameters here ...
)

func main() {

	// Initialize repository
	var repository = persistence.NewInMemorySignatureDeviceRepository()
	keyPairFactory := crypto.NewKeyPairFactoryImpl()
	signerFactory := crypto.NewSignerFactoryImpl()
	// Initialize service layer
	service := application.NewSignatureDeviceService(repository, keyPairFactory, signerFactory)
	//api handler
	apihandler := api.NewDeviceHTTPHandler(*service)
	server := api.NewServer(ListenAddress)
	server.SetupDeviceApiHandlers(*apihandler)

	if err := server.Run(); err != nil {
		//log.Fatal("Could not start server on ", ListenAddress)
		log.Fatal("Error: ", err.Error())

	}
}
