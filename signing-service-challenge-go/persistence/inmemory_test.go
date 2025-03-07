package persistence

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestAddDevice(t *testing.T) {
	db := MakeInMemoryDB()

	device := Device{ID: "1", Algorithm: "RSA", Label: nil}
	err := db.AddDevice(device)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Try adding the same device again
	err = db.AddDevice(device)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetDevice(t *testing.T) {
	db := MakeInMemoryDB()

	device := Device{ID: "1", Algorithm: "RSA", Label: nil}
	db.AddDevice(device)

	retrievedDevice, err := db.GetDevice("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(retrievedDevice, &device) {
		t.Fatalf("expected %v, got %v", retrievedDevice, device)
	}

	// Try getting a non-existent device
	_, err = db.GetDevice("2")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetAndSetSignature(t *testing.T) {
	db := MakeInMemoryDB()

	device := Device{ID: "1", Algorithm: "RSA", Label: nil}
	db.AddDevice(device)

	// Set initial signature
	err := db.CompareAndSwapSignature("1", "", "hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Update signature
	err = db.CompareAndSwapSignature("1", "hello", "world")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Try updating with wrong previous signature
	err = db.CompareAndSwapSignature("1", "hello", "world")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetAndSetSignatureParalell(t *testing.T) {
	db := MakeInMemoryDB()

	for i := 0; i < 5; i++ {
		d := Device{ID: ID(strconv.Itoa(i)), Algorithm: "RSA", Label: nil}
		db.AddDevice(d)
	}

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			for j := 0; j < 1000; j++ {
				// Test fails on Fatal error if conrurrent read and write occurs.
				// so no error handling is done here.
				id := ID(strconv.Itoa(i))
				d, err := db.GetDevice(id)
				if err != nil {
					break
				}
				_ = db.CompareAndSwapSignature(id, d.LastSignature, Signature(strconv.Itoa(j)))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestListDevices(t *testing.T) {
	db := MakeInMemoryDB()

	// Add devices
	for i := 0; i < 3; i++ {
		device := Device{ID: ID(strconv.Itoa(i)), Algorithm: "RSA", Label: nil}
		db.AddDevice(device)
	}

	// List devices
	ids, err := db.ListDevices()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("expected 3 devices, got %v", len(ids))
	}

	expectedIDs := []ID{"0", "1", "2"}
	for i, id := range ids {
		if id != expectedIDs[i] {
			t.Fatalf("expected id %v, got %v", expectedIDs[i], id)
		}
	}
}
