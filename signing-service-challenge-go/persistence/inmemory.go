package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// InMemorySignatureDeviceRepository represents in-memory storage for signature devices
type InMemorySignatureDeviceRepository struct {
	devices      map[string]*domain.SignatureDevice
	transactions map[string][]*domain.SignTransaction
	mu           sync.RWMutex
}

// NewInMemorySignatureDeviceRepository creates a new instance of InMemorySignatureDeviceRepository
func NewInMemorySignatureDeviceRepository() *InMemorySignatureDeviceRepository {
	return &InMemorySignatureDeviceRepository{
		devices:      make(map[string]*domain.SignatureDevice),
		transactions: make(map[string][]*domain.SignTransaction),
	}
}

// AddDevice adds a new signature device to the repository
func (r *InMemorySignatureDeviceRepository) AddDevice(device *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.devices[device.ID.String()]; exists {
		return errors.New("device ID already exists")
	}
	r.devices[device.ID.String()] = device
	return nil
}
func (r *InMemorySignatureDeviceRepository) UpdateDevice(device *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.devices[device.ID.String()]; exists {

		r.devices[device.ID.String()] = device
		return nil
	}
	return errors.New("device not found")
}

// GetDeviceByID retrieves a signature device by its ID
func (r *InMemorySignatureDeviceRepository) GetDeviceByID(id string) (*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	device, exists := r.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

// ListDevices returns a list of all signature devices
func (r *InMemorySignatureDeviceRepository) ListDevices() ([]*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	devices := make([]*domain.SignatureDevice, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *InMemorySignatureDeviceRepository) SaveSignTransaction(device *domain.SignatureDevice, tx *domain.SignTransaction) error {
	r.UpdateDevice(device)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions[tx.ID] = append(r.transactions[tx.ID], tx)
	return nil
}

func (r *InMemorySignatureDeviceRepository) GetSignTransactionsForDevice(id string) ([]*domain.SignTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	transactions, ok := r.transactions[id]
	if ok {
		return transactions, nil
	}
	return nil, errors.New("invalid device id")
}
