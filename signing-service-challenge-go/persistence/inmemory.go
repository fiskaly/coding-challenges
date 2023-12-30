package persistence

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
)

type inMemoryDeviceRepository struct {
	repository map[string]domain.SignatureDevice
	lock       sync.Mutex
}

// insance functions as a singleton
var instance *inMemoryDeviceRepository
var lock = &sync.Mutex{}

// Gets the instance of the in memory device repository.
// If no instance exists, a new instance is created.
func Get() *inMemoryDeviceRepository {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &inMemoryDeviceRepository{
				repository: make(map[string]domain.SignatureDevice),
				lock:       sync.Mutex{},
			}
		}
	}
	return instance
}

// Find device by id. Returns device and boolean marking if the device was found
func (dr *inMemoryDeviceRepository) FindDeviceById(id string) (domain.SignatureDevice, bool) {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	device, exists := dr.repository[id]
	return device, exists
}

// Create new device. Will return error if device with same ID already exists.
func (dr *inMemoryDeviceRepository) NewDevice(device domain.SignatureDevice) error {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	_, exists := dr.repository[device.Id]
	if exists {
		return fmt.Errorf("[NewDevice] device with specified ID already exists: \"%s\"", device.Id)
	}
	dr.repository[device.Id] = device
	return nil
}

// Update existing device. Will return error if evice with same ID does not exist.
func (dr *inMemoryDeviceRepository) UpdateDevice(device domain.SignatureDevice) error {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	_, exists := dr.repository[device.Id]
	if !exists {
		return fmt.Errorf("[UpdateDevice] device with specified ID doesn't exist: \"%s\"", device.Id)
	}
	dr.repository[device.Id] = device
	return nil
}

// Returns ID's of all existing devices
func (dr *inMemoryDeviceRepository) GetAllDevices() []string {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	keys := make([]string, 0, len(dr.repository))
	for k := range dr.repository {
		keys = append(keys, k)
	}

	return keys
}
