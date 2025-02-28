package domain

import (
	"strings"
	"sync"
	"testing"
)

func TestCreateNewDevice_Valid(t *testing.T) {
	device, err := CreateNewDevice("RSA", "Test Device")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if device.GetAlgorithm() != "RSA" {
		t.Errorf("Expected algorithm to be 'RSA', got %s", device.GetAlgorithm())
	}

	if device.GetLabel() != "Test Device" {
		t.Errorf("Expected label to be 'Test Device', got %s", device.GetLabel())
	}

	if device.GetSignatureCounter() != 0 {
		t.Errorf("Expected signature counter to be 0, got %d", device.GetSignatureCounter())
	}
}

func TestCreateNewDevice_InvalidAlgorithm(t *testing.T) {
	_, err := CreateNewDevice("InvalidAlgo", "Test Device")
	if err == nil {
		t.Fatalf("Expected error for invalid algorithm, got nil")
	}
}

func TestSignData_Valid(t *testing.T) {
	device, err := CreateNewDevice("RSA", "Test Device")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	signature, securedData, err := device.SignData("test_data")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if signature == "" {
		t.Errorf("Expected non-empty signature, got empty string")
	}

	if !strings.Contains(securedData, "test_data") {
		t.Errorf("Expected secured data to contain 'test_data', got %s", securedData)
	}

	if device.GetSignatureCounter() != 1 {
		t.Errorf("Expected signature counter to be 1, got %d", device.GetSignatureCounter())
	}
}

func TestSignData_Concurrency(t *testing.T) {
	device, err := CreateNewDevice("RSA", "Concurrent Device")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var wg sync.WaitGroup
	numRoutines := 100
	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			_, _, err := device.SignData("concurrent_data")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	if device.GetSignatureCounter() != numRoutines {
		t.Errorf("Expected signature counter to be %d, got %d", numRoutines, device.GetSignatureCounter())
	}
}
