package api

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// Tools and variables used to test the APIs of the server

const (
	MIN_PORT = 4000
	MAX_PORT = 10000
)

type respTypeDevice struct {
	Data deviceOut `json:"data"`
}
type respTypeDeviceArr struct {
	Data []deviceOut `json:"data"`
}
type respTypeHealth struct {
	Data HealthResponse `json:"data"`
}
type respTypeSign struct {
	Data signatureOut `json:"data"`
}

func setupServer() *Server {

	var randPort int

	// check if port is available
	for {

		// get a random port number withing the interval
		randPort = rand.Intn(MAX_PORT-MIN_PORT) + MIN_PORT

		// try to connect
		var addr *net.TCPAddr
		var err error
		if addr, err = net.ResolveTCPAddr("tcp", net.JoinHostPort("localhost", fmt.Sprint(randPort))); err == nil {
			var l *net.TCPListener
			if l, err = net.ListenTCP("tcp", addr); err == nil {
				l.Close()
				log.Printf("found empty port on '%d'", randPort)
				break
			}
		}
	}

	// let the previous connection close
	time.Sleep(200 * time.Millisecond)

	// create server
	server := NewServer(fmt.Sprintf(":%d", randPort))

	// run server in goroutine
	go func() {
		if err := server.Run(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Fatalf("error while running server on '%d': %s", randPort, err)
			}
		}
	}()

	// let the server spawn
	time.Sleep(200 * time.Millisecond)

	// check server is working
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v0/health", fmt.Sprint(randPort)))
	if err != nil {
		log.Fatal("server not responding")
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal("server not in goos state")
	}

	return server
}

func closeServer(server *Server) {
	if err := server.Close(); err != nil {
		log.Fatalf("error closing the server: %s", err)
	}
}
