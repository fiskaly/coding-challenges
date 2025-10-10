package service

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
)

// DeviceService handles business logic for signature devices
// Design Decision: Service layer separates business logic from HTTP handlers
// This improves testability and allows reuse across different interfaces (HTTP, gRPC, CLI, etc.)
type DeviceService struct {
	repository persistence.DeviceRepository
}

// NewDeviceService creates a new device service
func NewDeviceService(repository persistence.DeviceRepository) *DeviceService {
	return &DeviceService{
		repository: repository,
	}
}

// CreateDeviceRequest represents the request to create a new device
type CreateDeviceRequest struct {
	Algorithm domain.SignatureAlgorithm `json:"algorithm"`
	Label     string                    `json:"label,omitempty"`
}

// CreateDeviceResponse represents the response after creating a device
type CreateDeviceResponse struct {
	Device *domain.SignatureDeviceView `json:"device"`
}

// SignTransactionRequest represents the request to sign transaction data
type SignTransactionRequest struct {
	DeviceID string `json:"device_id"`
	Data     string `json:"data"`
}

// SignTransactionResponse represents the signature response
type SignTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

// CreateDevice creates a new signature device with the specified algorithm
// ID is generated using google/uuid library
func (s *DeviceService) CreateDevice(req CreateDeviceRequest) (*CreateDeviceResponse, error) {
	// Validate algorithm
	if err := req.Algorithm.Validate(); err != nil {
		return nil, err
	}

	// Generate UUID for the device (server-side generation)
	deviceID := uuid.New().String()

	var device *domain.SignatureDevice
	var err error

	// Generate key pair using existing generators and create device
	switch req.Algorithm {
	case domain.AlgorithmRSA:
		generator := &crypto.RSAGenerator{}
		keyPair, genErr := generator.Generate()
		if genErr != nil {
			return nil, fmt.Errorf("failed to generate RSA key pair: %w", genErr)
		}
		device, err = domain.NewSignatureDeviceWithRSA(deviceID, req.Label, keyPair)

	case domain.AlgorithmECC:
		generator := &crypto.ECCGenerator{}
		keyPair, genErr := generator.Generate()
		if genErr != nil {
			return nil, fmt.Errorf("failed to generate ECC key pair: %w", genErr)
		}
		device, err = domain.NewSignatureDeviceWithECC(deviceID, req.Label, keyPair)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", req.Algorithm)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Store the device
	if err := s.repository.Create(device); err != nil {
		return nil, fmt.Errorf("failed to store device: %w", err)
	}

	return &CreateDeviceResponse{
		Device: device.ToPublicView(),
	}, nil
}

// SignTransaction signs transaction data using the specified device
func (s *DeviceService) SignTransaction(req SignTransactionRequest) (*SignTransactionResponse, error) {
	if req.DeviceID == "" {
		return nil, errors.New("device ID is required")
	}

	if req.Data == "" {
		return nil, errors.New("data to sign is required")
	}

	// Retrieve the device
	device, err := s.repository.GetByID(req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device: %w", err)
	}

	// Prepare the data to be signed according to spec:
	// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
	securedData := device.PrepareDataToSign(req.Data)

	var signer crypto.Signer
	switch device.Algorithm {
	case domain.AlgorithmRSA:
		keyPair := device.GetRSAKeyPair()
		if keyPair == nil {
			return nil, errors.New("RSA key pair not found for device")
		}
		signer, err = crypto.NewRSASigner(keyPair)
		if err != nil {
			return nil, fmt.Errorf("failed to create RSA signer: %w", err)
		}

	case domain.AlgorithmECC:
		keyPair := device.GetECCKeyPair()
		if keyPair == nil {
			return nil, errors.New("ECC key pair not found for device")
		}
		signer, err = crypto.NewECDSASigner(keyPair)
		if err != nil {
			return nil, fmt.Errorf("failed to create ECDSA signer: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", device.Algorithm)
	}

	// Sign the secured data using the Signer interface
	signatureBytes, err := signer.Sign([]byte(securedData))
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	// Encode signature as base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes)

	// Update the device state (increment counter, update last signature)
	// Design Decision: This is atomic and thread-safe thanks to the mutex in the device
	device.UpdateAfterSigning(signatureBase64)

	return &SignTransactionResponse{
		Signature:  signatureBase64,
		SignedData: securedData,
	}, nil
}

// GetDevice retrieves a device by its ID
func (s *DeviceService) GetDevice(id string) (*domain.SignatureDeviceView, error) {
	if id == "" {
		return nil, errors.New("device ID is required")
	}

	device, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return device.ToPublicView(), nil
}

// ListDevices retrieves all devices
func (s *DeviceService) ListDevices() ([]*domain.SignatureDeviceView, error) {
	devices, err := s.repository.List()
	if err != nil {
		return nil, err
	}

	// Convert to public views
	views := make([]*domain.SignatureDeviceView, len(devices))
	for i, device := range devices {
		views[i] = device.ToPublicView()
	}

	return views, nil
}
