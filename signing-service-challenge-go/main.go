package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	p := persistence.MakeInMemoryDB()
	d := domain.MakeSignatureDomain(&p)
	server := api.NewServer(ListenAddress, domain.SignatureDomain(&d))

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
