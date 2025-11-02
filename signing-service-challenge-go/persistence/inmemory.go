package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// InMemorySignatureDeviceRepository provides a thread-safe map-backed storage.
type InMemorySignatureDeviceRepository struct {
	mu      sync.RWMutex
	devices map[string]*domain.SignatureDevice
}

// NewInMemorySignatureDeviceRepository constructs a repository instance.
func NewInMemorySignatureDeviceRepository() *InMemorySignatureDeviceRepository {
	return &InMemorySignatureDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// Save persists a new device, returning ErrDeviceAlreadyExists when the id is already taken.
func (r *InMemorySignatureDeviceRepository) Save(ctx context.Context, device *domain.SignatureDevice) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if device == nil {
		return errors.New("persistence: cannot save nil device")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID()]; exists {
		return domain.ErrDeviceAlreadyExists
	}

	clone, err := cloneDevice(device)
	if err != nil {
		return err
	}

	r.devices[device.ID()] = clone
	return nil
}

// Update stores the modified device state, returning ErrDeviceNotFound if the device is unknown.
func (r *InMemorySignatureDeviceRepository) Update(ctx context.Context, device *domain.SignatureDevice) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if device == nil {
		return errors.New("persistence: cannot update nil device")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID()]; !exists {
		return domain.ErrDeviceNotFound
	}

	clone, err := cloneDevice(device)
	if err != nil {
		return err
	}

	r.devices[device.ID()] = clone
	return nil
}

// Get retrieves a device by id, cloning it to decouple mutations from the storage.
func (r *InMemorySignatureDeviceRepository) Get(ctx context.Context, id string) (*domain.SignatureDevice, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, domain.ErrDeviceNotFound
	}

	return cloneDevice(device)
}

// List returns all known devices, each as an independent clone.
func (r *InMemorySignatureDeviceRepository) List(ctx context.Context) ([]*domain.SignatureDevice, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.SignatureDevice, 0, len(r.devices))
	for _, device := range r.devices {
		clone, err := cloneDevice(device)
		if err != nil {
			return nil, err
		}
		result = append(result, clone)
	}

	return result, nil
}

func cloneDevice(device *domain.SignatureDevice) (*domain.SignatureDevice, error) {
	if device == nil {
		return nil, errors.New("persistence: cannot clone nil device")
	}

	return domain.RestoreSignatureDevice(
		device.ID(),
		device.Algorithm(),
		device.Label(),
		device.SignatureCounter(),
		device.LastSignature(),
	)
}
