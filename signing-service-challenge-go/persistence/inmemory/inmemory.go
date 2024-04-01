package inmemory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// TODO: in-memory persistence ...
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

// // AddTransaction stores a new Transaction.
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
		return nil, fmt.Errorf("no device found with the id: %s", deviceID)
	}

	return sigDevice, nil
}

// // GetDevice retrieves a SignatureDevice by ID.
// func (store *InMemoryStore) GetDevice(id string) (*domain.SignatureDevice, bool) {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()

// 	device, exists := store.Devices[id]
// 	return device, exists
// }

// // GetTransaction retrieves a Transaction by ID.
// func (store *InMemoryStore) GetTransaction(id string) (*domain.Transaction, bool) {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()

// 	transaction, exists := store.Transactions[id]
// 	return transaction, exists
// }

// IncrementSignatureCounter increments the signature counter for a given device ID.
func (ms *InMemoryStore) IncrementSignatureCounter(deviceID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	device, exists := ms.Devices[deviceID]
	if !exists {
		return fmt.Errorf("signature device with ID %s not found", deviceID)
	}

	device.SignatureCounter++
	return nil
}
