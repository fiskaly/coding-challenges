package service

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

// TestCreateDevice tests device creation with different algorithms
func TestCreateDevice(t *testing.T) {
	tests := []struct {
		name      string
		algorithm domain.SignatureAlgorithm
		label     string
		wantError bool
	}{
		{
			name:      "Create RSA device",
			algorithm: domain.AlgorithmRSA,
			label:     "Test RSA Device",
			wantError: false,
		},
		{
			name:      "Create ECC device",
			algorithm: domain.AlgorithmECC,
			label:     "Test ECC Device",
			wantError: false,
		},
		{
			name:      "Create device with empty label",
			algorithm: domain.AlgorithmRSA,
			label:     "",
			wantError: false,
		},
		{
			name:      "Invalid algorithm",
			algorithm: "INVALID",
			label:     "Test",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := persistence.NewInMemoryDeviceRepository()
			svc := NewDeviceService(repo)

			req := CreateDeviceRequest{
				Algorithm: tt.algorithm,
				Label:     tt.label,
			}

			resp, err := svc.CreateDevice(req)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if resp == nil {
					t.Errorf("Expected response but got nil")
				}
				if resp != nil {
					// Verify UUID was generated
					if resp.Device.ID == "" {
						t.Errorf("Device ID should be generated, got empty string")
					}
					if resp.Device.SignatureCounter != 0 {
						t.Errorf("Initial signature counter should be 0, got %d", resp.Device.SignatureCounter)
					}
					if resp.Device.Label != tt.label {
						t.Errorf("Label mismatch: got %s, want %s", resp.Device.Label, tt.label)
					}
				}
			}
		})
	}
}

// TestSignTransaction tests the signing workflow
func TestSignTransaction(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	// Create a device first (ID is auto-generated now)
	createResp, err := svc.CreateDevice(CreateDeviceRequest{
		Algorithm: domain.AlgorithmRSA,
		Label:     "Test Device",
	})
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}
	deviceID := createResp.Device.ID

	// Sign first transaction
	req1 := SignTransactionRequest{
		DeviceID: deviceID,
		Data:     "transaction_data_1",
	}

	resp1, err := svc.SignTransaction(req1)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	if resp1.Signature == "" {
		t.Error("Signature should not be empty")
	}

	// Verify signed data format: <counter>_<data>_<last_signature>
	expectedPrefix := "0_transaction_data_1_"
	if len(resp1.SignedData) < len(expectedPrefix) {
		t.Errorf("SignedData too short: %s", resp1.SignedData)
	}

	// Sign second transaction
	req2 := SignTransactionRequest{
		DeviceID: deviceID,
		Data:     "transaction_data_2",
	}

	resp2, err := svc.SignTransaction(req2)
	if err != nil {
		t.Fatalf("Failed to sign second transaction: %v", err)
	}

	// Verify counter incremented and last signature is chained
	expectedPrefix2 := "1_transaction_data_2_" + resp1.Signature
	if resp2.SignedData != expectedPrefix2 {
		t.Errorf("SignedData mismatch:\ngot:  %s\nwant: %s", resp2.SignedData, expectedPrefix2)
	}

	// Verify signature counter increased
	device, err := svc.GetDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}

	if device.SignatureCounter != 2 {
		t.Errorf("Expected counter to be 2, got %d", device.SignatureCounter)
	}
}

// TestConcurrentSigning tests thread-safety of signing
func TestConcurrentSigning(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	// Create a device (ID is auto-generated)
	createResp, err := svc.CreateDevice(CreateDeviceRequest{
		Algorithm: domain.AlgorithmECC,
		Label:     "Concurrent Test Device",
	})
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}
	deviceID := createResp.Device.ID

	// Sign 100 transactions concurrently
	const numSignatures = 100
	errChan := make(chan error, numSignatures)
	doneChan := make(chan bool, numSignatures)

	for i := 0; i < numSignatures; i++ {
		go func(n int) {
			req := SignTransactionRequest{
				DeviceID: deviceID,
				Data:     "concurrent_transaction",
			}
			_, err := svc.SignTransaction(req)
			if err != nil {
				errChan <- err
			}
			doneChan <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numSignatures; i++ {
		<-doneChan
	}
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("Concurrent signing error: %v", err)
	}

	// Verify counter is exactly numSignatures (no gaps, strictly monotonic)
	device, err := svc.GetDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}

	if device.SignatureCounter != numSignatures {
		t.Errorf("Expected counter to be %d, got %d", numSignatures, device.SignatureCounter)
	}
}

// TestListDevices tests listing devices
func TestListDevices(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	// Initially empty
	devices, err := svc.ListDevices()
	if err != nil {
		t.Fatalf("Failed to list devices: %v", err)
	}
	if len(devices) != 0 {
		t.Errorf("Expected 0 devices, got %d", len(devices))
	}

	// Create multiple devices (IDs auto-generated)
	for i := 0; i < 5; i++ {
		_, err := svc.CreateDevice(CreateDeviceRequest{
			Algorithm: domain.AlgorithmRSA,
			Label:     "Device " + string(rune('A'+i)),
		})
		if err != nil {
			t.Fatalf("Failed to create device: %v", err)
		}
	}

	// List devices
	devices, err = svc.ListDevices()
	if err != nil {
		t.Fatalf("Failed to list devices: %v", err)
	}
	if len(devices) != 5 {
		t.Errorf("Expected 5 devices, got %d", len(devices))
	}
}

// TestGetDevice tests retrieving a specific device
func TestGetDevice(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	// Create a device
	createResp, err := svc.CreateDevice(CreateDeviceRequest{
		Algorithm: domain.AlgorithmECC,
		Label:     "Test Device",
	})
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}
	deviceID := createResp.Device.ID

	// Get existing device
	device, err := svc.GetDevice(deviceID)
	if err != nil {
		t.Fatalf("Failed to get device: %v", err)
	}
	if device.ID != deviceID {
		t.Errorf("Device ID mismatch: got %s, want %s", device.ID, deviceID)
	}

	// Get non-existing device
	_, err = svc.GetDevice("non-existing-id")
	if err == nil {
		t.Error("Expected error for non-existing device")
	}
}

// TestSignWithoutDevice tests signing with non-existing device
func TestSignWithoutDevice(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	req := SignTransactionRequest{
		DeviceID: "non-existing-device",
		Data:     "test data",
	}

	_, err := svc.SignTransaction(req)
	if err == nil {
		t.Error("Expected error when signing with non-existing device")
	}
}

// TestMultipleDevicesUniqueness tests that each device gets a unique UUID
func TestMultipleDevicesUniqueness(t *testing.T) {
	repo := persistence.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo)

	// Create multiple devices
	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		resp, err := svc.CreateDevice(CreateDeviceRequest{
			Algorithm: domain.AlgorithmRSA,
			Label:     "Test Device",
		})
		if err != nil {
			t.Fatalf("Failed to create device: %v", err)
		}

		// Check for duplicate IDs
		if ids[resp.Device.ID] {
			t.Errorf("Duplicate device ID generated: %s", resp.Device.ID)
		}
		ids[resp.Device.ID] = true
	}

	if len(ids) != 10 {
		t.Errorf("Expected 10 unique IDs, got %d", len(ids))
	}
}
