package service

import (
	"errors"
	"signing-service-challenge/crypto"
	"signing-service-challenge/domain"
	"signing-service-challenge/mocks"
	"signing-service-challenge/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDeviceWithRSASuccess(t *testing.T) {

	mockRSAGenerator := &mocks.MockRSAGenerator{
		GenerateFunc: func() (*crypto.RSAKeyPair, error) {
			return &crypto.RSAKeyPair{}, nil
		},
	}

	mockRSAMarshaler := &mocks.MockRSAMarshaler{
		MarshalFunc: func(keyPair crypto.RSAKeyPair) ([]byte, []byte, error) {
			return []byte("publicKey"), []byte("privateKey"), nil
		},
	}

	mockDeviceRepository := mocks.NewMockDeviceRepository()

	deviceService := NewDefaultDeviceService(mockDeviceRepository, nil, mockRSAGenerator, nil, mockRSAMarshaler)

	device, err := deviceService.CreateDevice("rsa-sign-test", domain.RSAAlgorithm)

	assert.NoError(t, err, "expect no error when successfully created device")
	assert.NotNil(t, device, "expect created device to be not nil")
	assert.GreaterOrEqual(t, mockDeviceRepository.UpdateDeviceCallsCount, 1, "expect device to be updated in the repository")
	assert.Equal(t, "rsa-sign-test", device.Label, "expect device label to be rsa-sign-test")
	assert.Equal(t, domain.RSAAlgorithm, device.Algorithm, "expect device algorithm to be RSA")
	assert.Equal(t, []byte("publicKey"), device.PublicKey, "expect device public key to be publicKey")
	assert.Equal(t, []byte("privateKey"), device.PrivateKey, "expect device private key to be privateKey")
}

func TestCreateDeviceWithECCSuccess(t *testing.T) {

	mockECCGenerator := &mocks.MockECCGenerator{
		GenerateFunc: func() (*crypto.ECCKeyPair, error) {
			return &crypto.ECCKeyPair{}, nil
		},
	}

	mockECCMarshaler := &mocks.MockECCMarshaler{
		EncodeFunc: func(keyPair crypto.ECCKeyPair) ([]byte, []byte, error) {
			return []byte("publicKey"), []byte("privateKey"), nil
		},
	}

	mockDeviceRepository := mocks.NewMockDeviceRepository()

	deviceService := NewDefaultDeviceService(mockDeviceRepository, mockECCGenerator, nil, mockECCMarshaler, nil)

	device, err := deviceService.CreateDevice("ecdsa-sign-test", domain.ECCAlgorithm)

	assert.NoError(t, err, "expect no error when successfully created device")
	assert.NotNil(t, device, "expect created device to be not nil")
	assert.GreaterOrEqual(t, mockDeviceRepository.UpdateDeviceCallsCount, 1, "expect device to be updated in the repository")
	assert.Equal(t, "ecdsa-sign-test", device.Label, "expect device label to be ecdsa-sign-test")
	assert.Equal(t, domain.ECCAlgorithm, device.Algorithm, "expect device algorithm to be ECC")
	assert.Equal(t, []byte("publicKey"), device.PublicKey, "expect device public key to be publicKey")
	assert.Equal(t, []byte("privateKey"), device.PrivateKey, "expect device private key to be privateKey")
}

func TestCreateDeviceWithEncodingError(t *testing.T) {
	mockRSAGenerator := &mocks.MockRSAGenerator{
		GenerateFunc: func() (*crypto.RSAKeyPair, error) {
			return &crypto.RSAKeyPair{}, nil
		},
	}

	mockRSAMarshaler := &mocks.MockRSAMarshaler{
		MarshalFunc: func(keyPair crypto.RSAKeyPair) ([]byte, []byte, error) {
			return nil, nil, errors.New("error while encoding using rsa")
		},
	}

	mockDeviceRepository := persistence.NewInmemoryDeviceRepository()

	deviceService := NewDefaultDeviceService(mockDeviceRepository, nil, mockRSAGenerator, nil, mockRSAMarshaler)

	device, err := deviceService.CreateDevice("rsa-sign-test", domain.RSAAlgorithm)

	assert.Nil(t, device, "expect device to be nil when error occurs")
	assert.Error(t, err, "expect error when encoding error occurs")

	var encodingError *KeysEncodingError
	isKeysEncodingError := errors.As(err, &encodingError)
	assert.True(t, isKeysEncodingError, "expect error to be of type KeysEncodingError")

	if isKeysEncodingError {
		assert.Equal(t, "RSA", encodingError.Algorithm, "expect algorithm to be RSA")
		assert.EqualError(t, encodingError.Err, "error while encoding using rsa", "expect correct error message")
	}
}

func TestCreateDeviceWithKeysGenerationError(t *testing.T) {
	mockRSAGenerator := &mocks.MockRSAGenerator{
		GenerateFunc: func() (*crypto.RSAKeyPair, error) {
			return nil, errors.New("error while generating keys using rsa")
		},
	}

	mockDeviceRepository := persistence.NewInmemoryDeviceRepository()

	deviceService := NewDefaultDeviceService(mockDeviceRepository, nil, mockRSAGenerator, nil, nil)

	device, err := deviceService.CreateDevice("rsa-sign-test", domain.RSAAlgorithm)

	assert.Nil(t, device, "expect device to be nil when error occurs")
	assert.Error(t, err, "expect error when keys generation error occurs")

	var encodingError *KeysGenerationError
	isKeysGenerationError := errors.As(err, &encodingError)
	assert.True(t, isKeysGenerationError, "expect error to be of type KeysGenerationError")

	if isKeysGenerationError {
		assert.Equal(t, "RSA", encodingError.Algorithm)
		assert.EqualError(t, encodingError.Err, "error while generating keys using rsa", "expect correct error message")
	}
}
