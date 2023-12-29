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

var lock = &sync.Mutex{}
var instance *inMemoryDeviceRepository

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

func (dr *inMemoryDeviceRepository) FindDeviceById(id string) (domain.SignatureDevice, error) {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	device, exists := dr.repository[id]
	if !exists {
		return domain.SignatureDevice{}, fmt.Errorf("[FindDeviceById] device with specified ID doesn't exist: \"%s\"", id)
	}
	return device, nil
}

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

func (dr *inMemoryDeviceRepository) GetAllDevices() []string {
	dr.lock.Lock()
	defer dr.lock.Unlock()
	keys := make([]string, 0, len(dr.repository))
	for k := range dr.repository {
		keys = append(keys, k)
	}

	return keys
}
