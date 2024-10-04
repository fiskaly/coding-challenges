package service

import (
	"errors"
	"signing-service-challenge/domain"
	"signing-service-challenge/mocks"
	"signing-service-challenge/persistence"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTransactionSuccess(t *testing.T) {
	mockDeviceRepository := mocks.NewMockDeviceRepository()
	mockSigner := &mocks.MockSigner{
		SignFunc: func(dataToBeSigned []byte) ([]byte, error) {
			return []byte("test-signature"), nil
		},
	}

	mockDevice := &domain.Device{
		Id:               "test-id",
		Algorithm:        domain.RSAAlgorithm,
		Label:            "Test Device",
		Signer:           mockSigner,
		SignatureCounter: 0,
		LastSignature:    nil,
	}

	mockDeviceRepository.DeviceToReturn = mockDevice
	mockDeviceRepository.GetDeviceByIdFound = true
	transactionService := NewDefaultTransactionService(mockDeviceRepository)
	signedData, securedData, err := transactionService.SignTransaction("test-id", "test data")

	assert.NoError(t, err, "expect no error when signed data successfully")
	assert.Equal(t, []byte("test-signature"), signedData, "expect signed data to be correct")
	assert.Contains(t, string(securedData), "test data", "expect secured data to contain test data")
	assert.GreaterOrEqual(t, mockDeviceRepository.UpdateDeviceCallsCount, 1, "expect device to be updated in the repository")

	assert.Equal(t, 1, mockDevice.SignatureCounter, "expect signature device counter augmented")
	assert.Equal(t, []byte("test-signature"), mockDevice.LastSignature, "expect last signature to be correct")
}

func TestSignTransactionWithDeviceNotFound(t *testing.T) {
	mockDeviceRepository := mocks.NewMockDeviceRepository()
	mockDeviceRepository.GetDeviceByIdFound = false

	transactionService := NewDefaultTransactionService(mockDeviceRepository)

	signedData, securedData, err := transactionService.SignTransaction("invalid-id", "test-data")

	assert.Error(t, err, "expect error when device not found")
	assert.Nil(t, signedData, "expect signed data to be nil")
	assert.Nil(t, securedData, "expect secured signed data to be nil")

	var deviceNotFoundError *DeviceNotFoundError
	isDeviceNotFoundError := errors.As(err, &deviceNotFoundError)
	assert.True(t, isDeviceNotFoundError, "expect error to be of type DeviceNotFoundError")

}

func TestSignTransactionWithSignError(t *testing.T) {
	mockDeviceRepository := mocks.NewMockDeviceRepository()
	mockSigner := &mocks.MockSigner{
		SignFunc: func(dataToBeSigned []byte) ([]byte, error) {
			return nil, errors.New("some signing error")
		},
	}

	mockDevice := &domain.Device{
		Id:               "test-id",
		Algorithm:        domain.RSAAlgorithm,
		Label:            "Test device",
		Signer:           mockSigner,
		SignatureCounter: 0,
		LastSignature:    nil,
	}

	mockDeviceRepository.DeviceToReturn = mockDevice
	mockDeviceRepository.GetDeviceByIdFound = true

	transactionService := NewDefaultTransactionService(mockDeviceRepository)

	signedData, securedData, err := transactionService.SignTransaction("test-id", "test-data")

	assert.Error(t, err, "expect error when signing fails")
	assert.Nil(t, signedData, "expect signed data to be nil")
	assert.Nil(t, securedData, "expect secured data to be nil")

	assert.EqualError(t, err, "error while signing the data: some signing error", "expect correct error message")
	assert.Equal(t, 0, mockDeviceRepository.UpdateDeviceCallsCount, "expect device not be updated in memory when signing fails")
}

func TestDefaultTransactionService_ConcurrentSignTransaction(t *testing.T) {
	deviceRepository := persistence.NewInmemoryDeviceRepository()
	service := NewDefaultTransactionService(deviceRepository)

	device := &domain.Device{
		Id:               "test-id",
		Label:            "test-device",
		Algorithm:        domain.RSAAlgorithm,
		SignatureCounter: 0,
		Signer:           &mocks.MockSigner{},
	}
	deviceRepository.UpdateDevice(device)

	const concurrentTransactions = 50
	var wg sync.WaitGroup
	wg.Add(concurrentTransactions)

	for i := 0; i < concurrentTransactions; i++ {
		go func() {
			defer wg.Done()
			_, _, err := service.SignTransaction("test-id", "test-data")
			assert.NoError(t, err)
		}()
	}

	wg.Wait()

	updatedDevice, _ := deviceRepository.GetDeviceById("test-id")
	assert.Equal(t, concurrentTransactions, updatedDevice.SignatureCounter, "expect all transactions to be processed")
}
