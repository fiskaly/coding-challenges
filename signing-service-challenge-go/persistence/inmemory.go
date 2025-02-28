package persistence

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type InMemoryStore struct {
	devices map[string]*domain.SignatureDevice
	mu      sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

func (s *InMemoryStore) Save(device *domain.SignatureDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[device.GetID()]; exists {
		return errors.New("device with this ID already exists")
	}
	s.devices[device.GetID()] = device
	return nil
}

func (s *InMemoryStore) FindByID(id string) (*domain.SignatureDevice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, exists := s.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func (s *InMemoryStore) FindAll() ([]*domain.SignatureDevice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	devices := make([]*domain.SignatureDevice, 0, len(s.devices))
	for _, device := range s.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *InMemoryStore) Update(device *domain.SignatureDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[device.GetID()]; !exists {
		return fmt.Errorf("device not found")
	}
	// Avoid copying the device (and its embedded mutex)
	s.devices[device.GetID()] = device
	return nil
}
