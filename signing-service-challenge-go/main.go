package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	ListenAddress = ":8080"
	ReadTimeout   = 5
	WriteTimeout  = 10
	IdleTimeout   = 120
)

func main() {
	server := api.NewServer(ListenAddress, ReadTimeout, WriteTimeout, IdleTimeout)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
