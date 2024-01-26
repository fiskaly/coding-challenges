package persistence

import (
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type InMemorySignatureDeviceRepository struct {
	devices map[uuid.UUID]domain.SignatureDevice
}

func (repository InMemorySignatureDeviceRepository) Create(device domain.SignatureDevice) error {
	_, ok := repository.devices[device.ID]
	if ok {
		return errors.New(fmt.Sprintf("duplicate id: %s", device.ID))
	}

	repository.devices[device.ID] = device
	return nil
}

func (repository InMemorySignatureDeviceRepository) Find(id uuid.UUID) (domain.SignatureDevice, bool, error) {
	device, ok := repository.devices[id]
	if !ok {
		return domain.SignatureDevice{}, false, nil
	}

	return device, true, nil
}

func NewInMemorySignatureDeviceRepository() InMemorySignatureDeviceRepository {
	return InMemorySignatureDeviceRepository{
		devices: map[uuid.UUID]domain.SignatureDevice{},
	}
}
