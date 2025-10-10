package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// DeviceRepository defines the contract for device persistence operations
// Design Decision: Using repository pattern to abstract storage implementation
// This allows easy migration to a database in the future
type DeviceRepository interface {
	Create(device *domain.SignatureDevice) error
	GetByID(id string) (*domain.SignatureDevice, error)
	List() ([]*domain.SignatureDevice, error)
	Update(device *domain.SignatureDevice) error
	Delete(id string) error
}

// InMemoryDeviceRepository implements DeviceRepository using in-memory storage
// Design Decision: Using sync.RWMutex for thread-safety
// - Multiple concurrent reads are allowed (RLock)
// - Writes are exclusive (Lock)
// This is critical for the requirement: "The system will be used by many concurrent clients"
type InMemoryDeviceRepository struct {
	devices map[string]*domain.SignatureDevice
	mu      sync.RWMutex
}

// NewInMemoryDeviceRepository creates a new in-memory repository
func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// Create stores a new device
func (r *InMemoryDeviceRepository) Create(device *domain.SignatureDevice) error {
	if device == nil {
		return errors.New("device cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return errors.New("device with this ID already exists")
	}

	r.devices[device.ID] = device
	return nil
}

// GetByID retrieves a device by its ID
func (r *InMemoryDeviceRepository) GetByID(id string) (*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	return device, nil
}

// List retrieves all devices
// Design Decision: Returns a slice of devices (not a map) for easier API serialization
func (r *InMemoryDeviceRepository) List() ([]*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*domain.SignatureDevice, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

// Update updates an existing device
// Design Decision: This is primarily for updating the device state after signing
// The actual device object is stored by reference, so updates are reflected immediately
func (r *InMemoryDeviceRepository) Update(device *domain.SignatureDevice) error {
	if device == nil {
		return errors.New("device cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; !exists {
		return errors.New("device not found")
	}

	r.devices[device.ID] = device
	return nil
}

// Delete removes a device by its ID
func (r *InMemoryDeviceRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[id]; !exists {
		return errors.New("device not found")
	}

	delete(r.devices, id)
	return nil
}
