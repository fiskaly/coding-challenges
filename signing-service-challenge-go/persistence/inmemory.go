package persistence

import (
	"signing-service-challenge/domain"
	"sync"
)

type DeviceRepository interface {
	GetDeviceById(deviceId string) (*domain.Device, bool)
	UpdateDevice(Device *domain.Device)
	ListDevices() ([]*domain.Device, error)
}

// TODO: in-memory persistence ...
type InmemoryDeviceRepository struct {
	mu      *sync.RWMutex
	Devices map[string]*domain.Device
}

func NewInmemoryDeviceRepository() *InmemoryDeviceRepository {
	return &InmemoryDeviceRepository{
		mu:      &sync.RWMutex{},
		Devices: make(map[string]*domain.Device),
	}
}

func (r *InmemoryDeviceRepository) GetDeviceById(deviceId string) (*domain.Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	Device, found := r.Devices[deviceId]
	if !found {
		return nil, false
	}

	return Device, true
}

func (r *InmemoryDeviceRepository) UpdateDevice(Device *domain.Device) { //Will keep update as upsert operation for the inmemory implementation
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Devices[Device.Id] = Device
}

func (r *InmemoryDeviceRepository) ListDevices() ([]*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var devices []*domain.Device
	for _, device := range r.Devices {
		devices = append(devices, device)
	}
	return devices, nil
}
