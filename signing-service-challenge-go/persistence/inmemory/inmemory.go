package inmemory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMemoryStore struct {
	Devices      map[string]*domain.SignatureDevice
	Transactions map[string][]*domain.Transaction
	mu           sync.RWMutex
}

// NewInMemoryStore initializes a new in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Devices:      make(map[string]*domain.SignatureDevice),
		Transactions: make(map[string][]*domain.Transaction),
	}
}

// AddDevice stores a new SignatureDevice.
func (ms *InMemoryStore) AddSignatureDevice(device *domain.SignatureDevice) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, exists := ms.Devices[device.ID]
	if exists {
		return errors.New("device already exists")
	}

	ms.Devices[device.ID] = device
	return nil
}

// AddTransaction stores a new Transaction.
func (ms *InMemoryStore) AddTransaction(transaction *domain.Transaction) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.Transactions[transaction.DeviceID] = append(ms.Transactions[transaction.DeviceID], transaction)
}

func (ms *InMemoryStore) GetSignatureDevice(deviceID string) (*domain.SignatureDevice, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	sigDevice, exists := ms.Devices[deviceID]
	if !exists {
		return nil, fmt.Errorf("no device found with id: %s", deviceID)
	}

	return sigDevice, nil
}

// IncrementSignatureCounter increments the signature counter for a given device ID.
func (ms *InMemoryStore) GetTransactionsByDeviceID(deviceID string) ([]*domain.Transaction, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	transactions, exists := ms.Transactions[deviceID]
	if !exists {
		return nil, fmt.Errorf("no transactions found for device ID %s", deviceID)
	}

	return transactions, nil
}

func (ms *InMemoryStore) UpdateSignatureDevice(deviceID string, signature string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	device, exists := ms.Devices[deviceID]
	if !exists {
		return fmt.Errorf("device with ID %s not found", deviceID)
	}

	device.LastSignature = signature
	device.SignatureCounter++

	return nil
}
