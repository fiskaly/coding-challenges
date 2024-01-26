package persistence

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	t.Run("persists the device in memory", func(t *testing.T) {
		device := domain.SignatureDevice{
			ID:                uuid.New(),
			AlgorithmName:     "RSA",
			EncodedPrivateKey: []byte("SOME_RSA_KEY"),
			Label:             "my rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()

		if len(repository.devices) != 0 {
			t.Errorf("new repository should have 0 devices")
		}

		err := repository.Create(device)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if len(repository.devices) != 1 {
			t.Errorf("expected repository to contain 1 device, got: %d", len(repository.devices))
		}

		persistedDevice, ok := repository.devices[device.ID]
		if !ok {
			t.Error("expected device with id to be persisted")
		}
		diff := cmp.Diff(persistedDevice, device)
		if diff != "" {
			t.Errorf("unexpected difference between original and persisted device: %s", diff)
		}
	})

	t.Run("does not persist when id is not unique", func(t *testing.T) {
		id := uuid.New()
		alreadyExistingDevice := domain.SignatureDevice{
			ID:            id,
			AlgorithmName: "RSA",
			Label:         "already existing rsa key",
		}
		duplicateIdDevice := domain.SignatureDevice{
			ID:            id,
			AlgorithmName: "RSA",
			Label:         "new rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()
		repository.devices[id] = alreadyExistingDevice
		if len(repository.devices) != 1 {
			t.Errorf("repository should contain 1 device")
		}

		err := repository.Create(duplicateIdDevice)
		if err == nil {
			t.Error("expected error")
		}
		if len(repository.devices) != 1 {
			t.Errorf("expected repository to contain 1 device, got: %d", len(repository.devices))
		}

		persistedDevice, ok := repository.devices[id]
		if !ok {
			t.Error("expected device with id to be present")
		}
		diff := cmp.Diff(persistedDevice, alreadyExistingDevice)
		if diff != "" {
			t.Errorf("expected persisted device to not have changed. diff: %s", diff)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("returns the device when device with id exists", func(t *testing.T) {
		device := domain.SignatureDevice{
			ID:                uuid.New(),
			AlgorithmName:     "RSA",
			EncodedPrivateKey: []byte("SOME_RSA_KEY"),
			Label:             "my rsa key",
		}

		repository := NewInMemorySignatureDeviceRepository()
		repository.devices[device.ID] = device

		foundDevice, found, err := repository.Find(device.ID)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if !found {
			t.Error("expected device to be found")
		}
		diff := cmp.Diff(foundDevice, device)
		if diff != "" {
			t.Errorf("unexpected difference between original and found device: %s", diff)
		}
	})

	t.Run("returns false when device with id does not exist", func(t *testing.T) {
		id := uuid.New()
		repository := NewInMemorySignatureDeviceRepository()

		_, found, err := repository.Find(id)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
		if found {
			t.Error("expected found: false")
		}
	})
}
