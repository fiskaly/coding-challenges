package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

// Storage interface for SignatureDevice management
type Storage interface {
	Save(device *domain.SignatureDevice) error
	FindByID(id string) (*domain.SignatureDevice, error)
	FindAll() ([]*domain.SignatureDevice, error)
	Update(device *domain.SignatureDevice) error
}
