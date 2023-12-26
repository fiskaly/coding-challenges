package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMemoryDeviceRepository struct {
	repository map[string]domain.SignatureDevice
	sync.Mutex
}

func (dr *InMemoryDeviceRepository) FindDeviceById(id string) (domain.SignatureDevice, error) {
	device, exists := dr.repository[id]
	if !exists {
		return domain.SignatureDevice{}, errors.New("device with specified ID doesn't exist")
	}
	return device, nil
}

func (dr *InMemoryDeviceRepository) NewDevice(device domain.SignatureDevice) error {
	_, exists := dr.repository[device.Id]
	if exists {
		return errors.New("device with specified ID already exists")
	}
	dr.repository[device.Id] = device
	return nil
}

func (dr *InMemoryDeviceRepository) UpdateDevice(device domain.SignatureDevice) error {
	_, exists := dr.repository[device.Id]
	if !exists {
		return errors.New("device with specified ID doesn't exist")
	}
	dr.repository[device.Id] = device
	return nil
}
