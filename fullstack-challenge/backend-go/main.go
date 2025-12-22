package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/fullstack-challenge/api"
)

const (
	ListenAddress = ":8080"
	// TODO: add database configuration parameters here ...
)

func main() {
	server := api.NewServer(ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
