package domain

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

const (
	AlgorithmRSA = "RSA"
	AlgorithmECC = "ECC"
)

type SigningService struct {
	deviceRepo SignatureDeviceRepository
	signLock   sync.Mutex
	createLock sync.Mutex
}

func NewSigningService(deviceRepo SignatureDeviceRepository) *SigningService {
	return &SigningService{
		deviceRepo: deviceRepo,
	}
}

func (service *SigningService) CreateSignatureDevice(id string, label string, algorithm string) (*SignatureDevice, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid id")
	}

	// If we do not lock here then multiple devices with the same id might be created
	service.createLock.Lock()
	defer service.createLock.Unlock()

	repo := service.deviceRepo
	existing := repo.GetSignatureDeviceById(id)
	if existing != nil {
		return nil, fmt.Errorf("device with id %s already exists", id)
	}

	device, err := service.createSignatureDevice(id, label, algorithm)
	if err != nil {
		return nil, err
	}

	err = repo.StoreSignatureDevice(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (service *SigningService) createSignatureDevice(id string, label string, algorithm string) (*SignatureDevice, error) {
	signer, err := createSigner(algorithm)
	if err != nil {
		return nil, err
	}

	device := NewSignatureDevice(id, label, algorithm, signer)

	return device, nil
}

func createSigner(algorithm string) (crypto.Signer, error) {
	switch strings.ToUpper(algorithm) {
	case AlgorithmRSA:
		return crypto.NewRSASigner()
	case AlgorithmECC:
		return crypto.NewECCSigner()
	default:
		return nil, errors.New("invalid algorithm")
	}
}

func (service *SigningService) ListSignatureDevices() ([]SignatureDevice, error) {
	return service.deviceRepo.ListSignatureDevices()
}

func (service *SigningService) GetSignatureDeviceById(id string) (*SignatureDevice, error) {
	device := service.deviceRepo.GetSignatureDeviceById(id)
	if device == nil {
		return nil, ErrorDeviceNotFound(id)
	}

	return device, nil
}

func (service *SigningService) SignTransaction(deviceId string, data []byte) (*SignDataResult, error) {
	// TODO: Improve performance
	// Locking indiscriminately on each sign is bad for performance because is locks out all devices from signing
	// One solution would be to keep track of active devices, counting each instance by id (map[string]uint)
	// For each device we would create a mutex to control signing
	service.signLock.Lock()
	defer service.signLock.Unlock()

	device := service.deviceRepo.GetSignatureDeviceById(deviceId)
	if device == nil {
		return nil, ErrorDeviceNotFound(deviceId)
	}

	result, err := device.Sign(data)
	if err != nil {
		return nil, err
	}

	service.deviceRepo.StoreSignatureDevice(device)
	if err != nil {
		return nil, err
	}

	return result, nil
}
