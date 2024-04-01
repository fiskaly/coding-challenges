package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Repository interface {
	AddSignatureDevice(device *domain.SignatureDevice) error
	GetSignatureDevice(deviceID string) (*domain.SignatureDevice, error)
	AddTransaction(transaction *domain.Transaction)
	IncrementSignatureCounter(deviceID string) error
}
