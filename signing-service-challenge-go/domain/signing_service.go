package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

const (
	AlgorithmRSA = "RSA"
	AlgorithmECC = "ECC"
)

type SigningService struct {
	deviceRepo SignatureDeviceRepository
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

	// TODO: Revise locking
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
	err = service.deviceRepo.StoreSignatureDevice(device)
	if err != nil {
		return nil, err
	}

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
	device := service.deviceRepo.GetSignatureDeviceById(deviceId)
	if device == nil {
		return nil, ErrorDeviceNotFound(deviceId)
	}

	result, err := device.Sign(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}
