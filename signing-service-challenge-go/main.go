package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	DEFAULT_SERVER_ADDR = ":8080"
)

func main() {

	// get the address from the env variables, otherwise set the
	// default value
	serverAddr := os.Getenv("FISAKLY_SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = DEFAULT_SERVER_ADDR
	}

	server := api.NewServer(serverAddr)

	if err := server.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Could not start server on ", serverAddr)
	}
}
