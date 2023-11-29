package domain

import (
	"encoding/base64"
	"errors"
	"fmt"

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

	signer, err := createSigner(algorithm)
	if err != nil {
		return nil, err
	}

	lastSignature := base64.StdEncoding.EncodeToString([]byte(id))
	device := SignatureDevice{
		Id:               id,
		Label:            label,
		Algorithm:        algorithm,
		signer:           signer,
		signatureCounter: 0,
		lastSignature:    lastSignature,
	}

	err = repo.StoreSignatureDevice(&device)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func createSigner(algorithm string) (crypto.Signer, error) {
	switch algorithm {
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
		return nil, fmt.Errorf("device with id %s does not exist", id)
	}

	return device, nil
}
