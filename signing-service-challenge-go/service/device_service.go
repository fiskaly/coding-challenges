package service

import (
	"signing-service-challenge/crypto"
	"signing-service-challenge/domain"
	"signing-service-challenge/persistence"

	"github.com/google/uuid"
)

type DeviceService interface {
	CreateDevice(label string, algorithmType domain.AlgorithmType) (*domain.Device, error)
	ListDevices() ([]*domain.Device, error)
	GetDeviceById(deviceId string) (*domain.Device, error)
}

type DefaultDeviceService struct {
	deviceRepository persistence.DeviceRepository
	rsaGenerator     crypto.RSAGenerator
	eccGenerator     crypto.ECCGenerator
	rsaMarshaler     crypto.RSAMarshaler
	ecdsaMarshaler   crypto.ECCMarshaler
}

func NewDefaultDeviceService(deviceRepository persistence.DeviceRepository,
	eccGenerator crypto.ECCGenerator,
	rsaGenerator crypto.RSAGenerator,
	eccMarshaler crypto.ECCMarshaler,
	rsaMarshaler crypto.RSAMarshaler) *DefaultDeviceService {
	return &DefaultDeviceService{deviceRepository: deviceRepository, eccGenerator: eccGenerator, rsaGenerator: rsaGenerator, ecdsaMarshaler: eccMarshaler, rsaMarshaler: rsaMarshaler}
}

func (s *DefaultDeviceService) CreateDevice(label string, algorithmType domain.AlgorithmType) (*domain.Device, error) {
	id := uuid.NewString()
	var publicKey, privateKey []byte
	var signer crypto.Signer
	switch algorithmType {
	case domain.RSAAlgorithm:
		rsakey, err := s.rsaGenerator.Generate()
		if err != nil {
			return nil, NewKeysGenerationError(string(algorithmType), err)
		}

		signer = crypto.NewRSASigner(rsakey.Private)
		publicKey, privateKey, err = s.rsaMarshaler.Marshal(*rsakey)
		if err != nil {
			return nil, NewKeysEncodingError(string(algorithmType), err)
		}

	case domain.ECCAlgorithm:
		eccKey, err := s.eccGenerator.Generate()
		if err != nil {
			return nil, NewKeysGenerationError(string(algorithmType), err)
		}

		signer = crypto.NewECCSigner(eccKey.Private)
		publicKey, privateKey, err = s.ecdsaMarshaler.Encode(*eccKey)
		if err != nil {
			return nil, NewKeysEncodingError(string(algorithmType), err)
		}
	default:
		return nil, NewInvalidAlgorithmError(string(algorithmType))
	}

	createdDevice := &domain.Device{
		Id:               id,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		SignatureCounter: 0,
		Algorithm:        algorithmType,
		Signer:           signer,
		Label:            label,
	}
	s.deviceRepository.UpdateDevice(createdDevice)
	return createdDevice, nil
}

func (s *DefaultDeviceService) GetDeviceById(deviceId string) (*domain.Device, error) {
	device, found := s.deviceRepository.GetDeviceById(deviceId)
	if !found {
		return nil, NewDeviceNotFoundError(deviceId)
	}
	return device, nil
}

func (s *DefaultDeviceService) ListDevices() ([]*domain.Device, error) {
	return s.deviceRepository.ListDevices()
}
