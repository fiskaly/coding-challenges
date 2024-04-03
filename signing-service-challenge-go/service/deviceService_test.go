package service

import (
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence/inmemory"
)

func TestConcurrentSigning(t *testing.T) {
	// Setup
	repo := inmemory.NewInMemoryStore()
	service := NewSignatureDeviceService(repo)

	deviceID := "dev-1"
	sigDevice, err := service.CreateSignatureDevice(deviceID, "RSA", "label-1")
	if err != nil {
		t.Error("failed to create signature device: ", err)
		return
	}

	// Number of concurent clients
	const clients = 20
	var wg sync.WaitGroup
	wg.Add(clients)

	// channel to collect counter values
	counters := make(chan int, clients)

	// Simulate concurrent clients
	for i := 0; i < clients; i++ {
		go func() {
			defer wg.Done()
			_, err := service.SignTransaction(deviceID, "data to be signed")
			if err != nil {
				t.Error("failed to sign transaction:", err)
				return
			}
			counters <- sigDevice.SignatureCounter
		}()
	}

	wg.Wait()
	close(counters)

	// Verify the counters
	seen := make(map[int]bool)
	for counter := range counters {
		t.Logf("counter: %d", counter)

		if seen[counter] {
			t.Errorf("duplicate signature_counter detected: %d", counter)
		}
		seen[counter] = true
	}

	if len(seen) != clients {
		t.Errorf("expected %d unique counters, got %d", clients, len(seen))
	} else {
		t.Logf("Success: %d unique signature_counter values were detected, as expected.", clients)
	}
}
