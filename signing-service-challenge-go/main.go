package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/config"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

func main() {
	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Initialize the repository (in-memory storage)
	repository := persistence.NewInMemoryDeviceRepository()

	// Initialize the service layer
	deviceService := service.NewDeviceService(repository)

	// Initialize the HTTP server
	server := api.NewServer(cfg.GetListenAddr(), deviceService)

	// Create HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         cfg.GetListenAddr(),
		Handler:      server.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting signature service on %s", cfg.GetListenAddr())
		log.Printf("Configuration: RSA_KEY_SIZE=%d, ECC_CURVE=%s", cfg.RSAKeySize, cfg.ECCCurve)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
