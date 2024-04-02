package service

import (
	"encoding/base64"
	"fmt"
	"strconv"

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

func (s *SignatureDeviceService) SignTransaction(deviceId string, data string) (*domain.Signature, error) {
	device, err := s.store.GetSignatureDevice(deviceId)
	if err != nil {
		return nil, err
	}

	// construct data
	dataToBeSigned := fmt.Sprintf("%s_%s_%s", strconv.Itoa(device.SignatureCounter), data, device.LastSignature)

	// signe transaction
	transaction, err := device.SignTransaction(dataToBeSigned)
	if err != nil {
		return nil, err
	}

	// update last signature
	device.LastSignature = string(transaction.Signature)

	// update the signatureCounter
	s.store.IncrementSignatureCounter(deviceId)

	// store transaction
	s.store.AddTransaction(transaction)

	return &domain.Signature{
		Signature:  base64.StdEncoding.EncodeToString(transaction.Signature),
		SignedData: dataToBeSigned,
	}, nil
}
