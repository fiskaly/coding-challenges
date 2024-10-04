package service

import (
	"fmt"
	"signing-service-challenge/persistence"
)

type TransactionService interface {
	SignTransaction(deviceId string, data string) ([]byte, []byte, error)
}

type DefaultTransactionService struct {
	deviceRepository persistence.DeviceRepository
}

func NewDefaultTransactionService(deviceRepository persistence.DeviceRepository) *DefaultTransactionService {
	return &DefaultTransactionService{deviceRepository: deviceRepository}
}

func (s *DefaultTransactionService) SignTransaction(deviceId string, data string) ([]byte, []byte, error) {
	device, found := s.deviceRepository.GetDeviceById(deviceId) //Get signature device to sign the data
	if !found {
		return nil, nil, NewDeviceNotFoundError(deviceId)
	}

	securedDataToBeSigned := device.BuildSecuredDataToBeSigned(data)            //Build secured data to be signed string
	signedSecuredData, err := device.Signer.Sign([]byte(securedDataToBeSigned)) //Sign secured data
	if err != nil {
		return nil, nil, fmt.Errorf("error while signing the data: %w", err)
	}

	//Update signature device values and update it in the repository if data was signed successfully
	device.IncrementSignatureCounter()
	device.UpdateLastSignature(signedSecuredData)

	s.deviceRepository.UpdateDevice(device)

	return signedSecuredData, []byte(securedDataToBeSigned), nil
}
