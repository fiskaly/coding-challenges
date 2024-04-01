package service

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type SignatureDeviceService struct {
	store persistence.Repository
}

func NewSignatureDeviceService(store persistence.Repository) *SignatureDeviceService {
	return &SignatureDeviceService{store: store}
}

func (s *SignatureDeviceService) CreateSignatureDevice(id string, algorithm string, label string) (*domain.SignatureDevice, error) {
	device, err := domain.NewSignatureDevice(id, algorithm, label)
	if err != nil {
		return nil, err
	}

	// store the new device in the repository.
	if err := s.store.AddSignatureDevice(device); err != nil {
		return nil, err
	}

	return device, nil
}

func (s *SignatureDeviceService) SignTransaction(deviceId string, data string) (*domain.Transaction, error) {
	// Using a real db all these operations would be executed
	// in a transaction in order to maintain data integrity
	device, err := s.store.GetSignatureDevice(deviceId)
	if err != nil {
		return nil, err
	}

	// signe transaction
	transaction, err := device.SignTransaction(data)
	if err != nil {
		return nil, err
	}

	// update the signatureCounter
	s.store.IncrementSignatureCounter(deviceId)

	// store transaction
	s.store.AddTransaction(transaction)

	return transaction, nil
}
