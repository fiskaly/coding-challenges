package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

type Repository interface {
	AddSignatureDevice(device *domain.SignatureDevice) error
	GetSignatureDevice(deviceID string) (*domain.SignatureDevice, error)
	UpdateSignatureDevice(deviceID string, signature string) error
	AddTransaction(transaction *domain.Transaction)
	GetTransactionsByDeviceID(deviceID string) ([]*domain.Transaction, error)
}
