package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// InMemorySignatureDeviceRepository represents in-memory storage for signature devices
type InMemorySignatureDeviceRepository struct {
	devices map[string]*domain.SignatureDevice
	sync.RWMutex
}

// NewInMemorySignatureDeviceRepository creates a new instance of InMemorySignatureDeviceRepository
func NewInMemorySignatureDeviceRepository() *InMemorySignatureDeviceRepository {
	return &InMemorySignatureDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// AddDevice adds a new signature device to the repository
func (r *InMemorySignatureDeviceRepository) AddDevice(device *domain.SignatureDevice) error {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.devices[device.ID.String()]; exists {
		return errors.New("device ID already exists")
	}
	r.devices[device.ID.String()] = device
	return nil
}

// GetDeviceByID retrieves a signature device by its ID
func (r *InMemorySignatureDeviceRepository) GetDeviceByID(id string) (*domain.SignatureDevice, error) {
	r.RLock()
	defer r.RUnlock()
	device, exists := r.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

// ListDevices returns a list of all signature devices
func (r *InMemorySignatureDeviceRepository) ListDevices() ([]*domain.SignatureDevice, error) {
	r.RLock()
	defer r.RUnlock()
	devices := make([]*domain.SignatureDevice, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	return devices, nil
}
