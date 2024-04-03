package service

import (
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

	// update last signature & increment signature counter
	err = s.store.UpdateSignatureDevice(deviceId, transaction.Signature)
	if err != nil {
		return nil, err
	}

	// store transaction
	s.store.AddTransaction(transaction)

	return &domain.Signature{
		Signature:  transaction.Signature,
		SignedData: dataToBeSigned,
	}, nil
}

func (s *SignatureDeviceService) GetSignatureDeviceByID(deviceId string) (*domain.SignatureDevice, error) {
	device, err := s.store.GetSignatureDevice(deviceId)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *SignatureDeviceService) GetTransactionsByDeviceID(deviceId string) (*domain.TransactionsByDeviceResp, error) {
	device, err := s.GetSignatureDeviceByID(deviceId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.store.GetTransactionsByDeviceID(device.ID)
	if err != nil {
		return nil, err
	}

	transactionResp := make([]*domain.TransactionResp, 0)
	for _, tr := range transactions {
		transactionResp = append(transactionResp, &domain.TransactionResp{
			ID:        tr.ID,
			Signature: tr.Signature,
			CreatedAt: tr.Timestamp,
		})
	}

	return &domain.TransactionsByDeviceResp{
		ID:           device.ID,
		Algorithm:    device.Algorithm,
		Label:        device.Label,
		CreatedAt:    device.CreatedAt,
		Transactions: transactionResp,
	}, nil
}
