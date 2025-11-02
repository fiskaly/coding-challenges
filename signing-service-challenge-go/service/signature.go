package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// Errors returned by the signature service layer.
var (
	ErrInvalidDeviceID = errors.New("service: device id must be non-empty")
	ErrInvalidData     = errors.New("service: data to be signed must be non-empty")
	ErrMissingSigner   = errors.New("service: signing material not found for device")
)

// SignatureService orchestrates signature device lifecycle and signing operations.
type SignatureService struct {
	repo domain.SignatureDeviceRepository

	rsaGenerator *crypto.RSAGenerator
	eccGenerator *crypto.ECCGenerator

	signers sync.Map // map[string]crypto.Signer
	locks   sync.Map // map[string]*sync.Mutex
}

// NewSignatureService constructs a SignatureService.
func NewSignatureService(repo domain.SignatureDeviceRepository) *SignatureService {
	return &SignatureService{
		repo:         repo,
		rsaGenerator: &crypto.RSAGenerator{},
		eccGenerator: &crypto.ECCGenerator{},
	}
}

// CreateSignatureDeviceRequest encapsulates input parameters for device creation.
type CreateSignatureDeviceRequest struct {
	ID        string
	Algorithm string
	Label     string
}

// CreateSignatureDeviceResponse represents the result of a creation operation.
type CreateSignatureDeviceResponse struct {
	ID               string
	Label            string
	Algorithm        domain.Algorithm
	SignatureCounter uint64
	LastSignature    string
	PublicKeyPEM     []byte
}

// SignatureDeviceDTO represents a view of a signature device for read operations.
type SignatureDeviceDTO struct {
	ID               string
	Label            string
	Algorithm        domain.Algorithm
	SignatureCounter uint64
	LastSignature    string
}

// CreateSignatureDevice validates the request, provisions keys/signers, and persists the device.
func (s *SignatureService) CreateSignatureDevice(ctx context.Context, req CreateSignatureDeviceRequest) (*CreateSignatureDeviceResponse, error) {
	if strings.TrimSpace(req.ID) == "" {
		return nil, ErrInvalidDeviceID
	}

	algorithm, err := domain.ParseAlgorithm(req.Algorithm)
	if err != nil {
		return nil, err
	}

	var (
		publicKey []byte
		signer    crypto.Signer
	)

	switch algorithm {
	case domain.AlgorithmRSA:
		keyPair, err := s.rsaGenerator.Generate()
		if err != nil {
			return nil, fmt.Errorf("service: generate RSA key pair: %w", err)
		}

		signer, err = crypto.NewRSASigner(keyPair)
		if err != nil {
			return nil, fmt.Errorf("service: initialise RSA signer: %w", err)
		}

		marshaler := crypto.NewRSAMarshaler()
		var private []byte
		publicKey, private, err = marshaler.Marshal(*keyPair)
		if err != nil {
			return nil, fmt.Errorf("service: marshal RSA key pair: %w", err)
		}

		_ = private // private key would be persisted by a real implementation.

	case domain.AlgorithmECC:
		keyPair, err := s.eccGenerator.Generate()
		if err != nil {
			return nil, fmt.Errorf("service: generate ECC key pair: %w", err)
		}

		signer, err = crypto.NewECDSASigner(keyPair)
		if err != nil {
			return nil, fmt.Errorf("service: initialise ECDSA signer: %w", err)
		}

		marshaler := crypto.NewECCMarshaler()
		var private []byte
		publicKey, private, err = marshaler.Encode(*keyPair)
		if err != nil {
			return nil, fmt.Errorf("service: marshal ECC key pair: %w", err)
		}
		_ = private
	default:
		return nil, domain.ErrUnsupportedAlgorithm
	}

	device, err := domain.NewSignatureDevice(req.ID, algorithm, req.Label)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, device); err != nil {
		return nil, err
	}

	s.signers.Store(device.ID(), signer)

	return &CreateSignatureDeviceResponse{
		ID:               device.ID(),
		Label:            device.Label(),
		Algorithm:        device.Algorithm(),
		SignatureCounter: device.SignatureCounter(),
		LastSignature:    device.LastSignature(),
		PublicKeyPEM:     publicKey,
	}, nil
}

// SignTransactionRequest encapsulates the parameters for signing.
type SignTransactionRequest struct {
	DeviceID string
	Data     string
}

// SignTransactionResponse contains the resulting signature and metadata.
type SignTransactionResponse struct {
	DeviceID         string
	Signature        string
	SignedData       string
	SignatureCounter uint64
	LastSignature    string
}

// SignTransaction coordinates payload preparation, signing, and state persistence.
func (s *SignatureService) SignTransaction(ctx context.Context, req SignTransactionRequest) (*SignTransactionResponse, error) {
	if strings.TrimSpace(req.DeviceID) == "" {
		return nil, ErrInvalidDeviceID
	}

	if strings.TrimSpace(req.Data) == "" {
		return nil, ErrInvalidData
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	lock := s.obtainDeviceLock(req.DeviceID)
	lock.Lock()
	defer lock.Unlock()

	device, err := s.repo.Get(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}

	secured, err := device.SecuredPayload(req.Data)
	if err != nil {
		return nil, err
	}

	signer, err := s.getSigner(req.DeviceID)
	if err != nil {
		return nil, err
	}

	signatureBytes, err := signer.Sign([]byte(secured))
	if err != nil {
		return nil, fmt.Errorf("service: sign payload: %w", err)
	}

	encodedSignature, err := device.RecordSignature(signatureBytes)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, err
	}

	return &SignTransactionResponse{
		DeviceID:         device.ID(),
		Signature:        encodedSignature,
		SignedData:       secured,
		SignatureCounter: device.SignatureCounter(),
		LastSignature:    device.LastSignature(),
	}, nil
}

// GetSignatureDevice fetches a device by identifier.
func (s *SignatureService) GetSignatureDevice(ctx context.Context, deviceID string) (*SignatureDeviceDTO, error) {
	if strings.TrimSpace(deviceID) == "" {
		return nil, ErrInvalidDeviceID
	}

	device, err := s.repo.Get(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	dto := toDeviceDTO(device)
	return &dto, nil
}

// ListSignatureDevices returns every registered device.
func (s *SignatureService) ListSignatureDevices(ctx context.Context) ([]SignatureDeviceDTO, error) {
	devices, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]SignatureDeviceDTO, 0, len(devices))
	for _, device := range devices {
		result = append(result, toDeviceDTO(device))
	}

	return result, nil
}

func toDeviceDTO(device *domain.SignatureDevice) SignatureDeviceDTO {
	return SignatureDeviceDTO{
		ID:               device.ID(),
		Label:            device.Label(),
		Algorithm:        device.Algorithm(),
		SignatureCounter: device.SignatureCounter(),
		LastSignature:    device.LastSignature(),
	}
}

func (s *SignatureService) getSigner(deviceID string) (crypto.Signer, error) {
	value, ok := s.signers.Load(deviceID)
	if !ok {
		return nil, ErrMissingSigner
	}

	signer, ok := value.(crypto.Signer)
	if !ok || signer == nil {
		return nil, ErrMissingSigner
	}

	return signer, nil
}

func (s *SignatureService) obtainDeviceLock(deviceID string) *sync.Mutex {
	value, _ := s.locks.LoadOrStore(deviceID, &sync.Mutex{})
	return value.(*sync.Mutex)
}
