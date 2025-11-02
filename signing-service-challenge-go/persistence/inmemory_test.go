package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func TestInMemoryRepositorySaveAndGet(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()
	device, err := domain.NewSignatureDevice("device-1", domain.AlgorithmRSA, "Label")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if err := repo.Save(context.Background(), device); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	stored, err := repo.Get(context.Background(), "device-1")
	if err != nil {
		t.Fatalf("expected to retrieve device, got %v", err)
	}

	if stored == device {
		t.Fatal("expected repository to return a clone, got same pointer")
	}

	if stored.Label() != "Label" {
		t.Fatalf("expected label %q, got %q", "Label", stored.Label())
	}

	stored.SetLabel("Mutated Label")
	reloaded, err := repo.Get(context.Background(), "device-1")
	if err != nil {
		t.Fatalf("expected to retrieve device again, got %v", err)
	}

	if reloaded.Label() != "Label" {
		t.Fatalf("expected stored device label to remain %q, got %q", "Label", reloaded.Label())
	}
}

func TestInMemoryRepositoryDetectsDuplicates(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()
	device, err := domain.NewSignatureDevice("device-dup", domain.AlgorithmECC, "")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if err := repo.Save(context.Background(), device); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	if err := repo.Save(context.Background(), device); err != domain.ErrDeviceAlreadyExists {
		t.Fatalf("expected ErrDeviceAlreadyExists, got %v", err)
	}
}

func TestInMemoryRepositoryUpdate(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()
	device, err := domain.NewSignatureDevice("device-update", domain.AlgorithmRSA, "initial")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if err := repo.Save(context.Background(), device); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	device.SetLabel("updated")
	if err := repo.Update(context.Background(), device); err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	updated, err := repo.Get(context.Background(), "device-update")
	if err != nil {
		t.Fatalf("expected to retrieve updated device, got %v", err)
	}

	if updated.Label() != "updated" {
		t.Fatalf("expected label to be updated, got %q", updated.Label())
	}
}

func TestInMemoryRepositoryUpdateFailsForUnknownDevice(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()
	device, err := domain.NewSignatureDevice("unknown-device", domain.AlgorithmRSA, "")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	if err := repo.Update(context.Background(), device); err != domain.ErrDeviceNotFound {
		t.Fatalf("expected ErrDeviceNotFound, got %v", err)
	}
}

func TestInMemoryRepositoryListReturnsClones(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()

	for i := 0; i < 3; i++ {
		device, err := domain.NewSignatureDevice(
			uniqueID(i),
			domain.AlgorithmRSA,
			"",
		)
		if err != nil {
			t.Fatalf("failed to create device %d: %v", i, err)
		}

		if err := repo.Save(context.Background(), device); err != nil {
			t.Fatalf("failed to save device %d: %v", i, err)
		}
	}

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}

	if len(list) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(list))
	}

	list[0].SetLabel("should not affect stored state")
	reloaded, err := repo.Get(context.Background(), uniqueID(0))
	if err != nil {
		t.Fatalf("expected to reload device, got %v", err)
	}

	if reloaded.Label() != "" {
		t.Fatalf("expected stored device label to remain empty, got %q", reloaded.Label())
	}
}

func uniqueID(i int) string {
	return fmt.Sprintf("device-%d", i)
}
