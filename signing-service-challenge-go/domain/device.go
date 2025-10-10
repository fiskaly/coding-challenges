package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// SignatureAlgorithm represents the cryptographic algorithm used for signing
type SignatureAlgorithm string

const (
	AlgorithmRSA SignatureAlgorithm = "RSA"
	AlgorithmECC SignatureAlgorithm = "ECC"
)

// Validate checks if the algorithm is supported
func (a SignatureAlgorithm) Validate() error {
	switch a {
	case AlgorithmRSA, AlgorithmECC:
		return nil
	default:
		return fmt.Errorf("unsupported algorithm: %s", a)
	}
}

// SignatureDevice represents a device capable of signing transaction data
// Design Decision: Using sync.RWMutex to ensure thread-safety for concurrent access
type SignatureDevice struct {
	ID               string             `json:"id"`
	Algorithm        SignatureAlgorithm `json:"algorithm"`
	Label            string             `json:"label,omitempty"`
	SignatureCounter int                `json:"signature_counter"`
	LastSignature    string             `json:"last_signature"` // base64 encoded

	// Private fields for internal use - now using the existing crypto package types
	rsaKeyPair *crypto.RSAKeyPair // For RSA algorithm
	eccKeyPair *crypto.ECCKeyPair // For ECC algorithm
	mu         sync.RWMutex       // Protects SignatureCounter and LastSignature
}

// NewSignatureDeviceWithRSA creates a new SignatureDevice with RSA algorithm
func NewSignatureDeviceWithRSA(id string, label string, keyPair *crypto.RSAKeyPair) (*SignatureDevice, error) {
	if id == "" {
		return nil, errors.New("device ID cannot be empty")
	}
	if keyPair == nil || keyPair.Private == nil || keyPair.Public == nil {
		return nil, errors.New("RSA key pair cannot be nil")
	}

	// Initial last_signature is base64-encoded device ID (as per spec)
	lastSignature := base64.StdEncoding.EncodeToString([]byte(id))

	return &SignatureDevice{
		ID:               id,
		Algorithm:        AlgorithmRSA,
		Label:            label,
		SignatureCounter: 0,
		LastSignature:    lastSignature,
		rsaKeyPair:       keyPair,
		eccKeyPair:       nil,
	}, nil
}

// NewSignatureDeviceWithECC creates a new SignatureDevice with ECC algorithm
func NewSignatureDeviceWithECC(id string, label string, keyPair *crypto.ECCKeyPair) (*SignatureDevice, error) {
	if id == "" {
		return nil, errors.New("device ID cannot be empty")
	}
	if keyPair == nil || keyPair.Private == nil || keyPair.Public == nil {
		return nil, errors.New("ECC key pair cannot be nil")
	}

	// Initial last_signature is base64-encoded device ID (as per spec)
	lastSignature := base64.StdEncoding.EncodeToString([]byte(id))

	return &SignatureDevice{
		ID:               id,
		Algorithm:        AlgorithmECC,
		Label:            label,
		SignatureCounter: 0,
		LastSignature:    lastSignature,
		rsaKeyPair:       nil,
		eccKeyPair:       keyPair,
	}, nil
}

// GetRSAKeyPair returns the RSA key pair (thread-safe read)
func (d *SignatureDevice) GetRSAKeyPair() *crypto.RSAKeyPair {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.rsaKeyPair
}

// GetECCKeyPair returns the ECC key pair (thread-safe read)
func (d *SignatureDevice) GetECCKeyPair() *crypto.ECCKeyPair {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.eccKeyPair
}

// GetSignatureCounter returns the current signature counter (thread-safe)
func (d *SignatureDevice) GetSignatureCounter() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.SignatureCounter
}

// GetLastSignature returns the last signature (thread-safe)
func (d *SignatureDevice) GetLastSignature() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.LastSignature
}

// PrepareDataToSign creates the secured data string as per specification:
// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
// Design Decision: This method is thread-safe and prepares data atomically
func (d *SignatureDevice) PrepareDataToSign(data string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return fmt.Sprintf("%d_%s_%s", d.SignatureCounter, data, d.LastSignature)
}

// UpdateAfterSigning updates the device state after a successful signature
// Design Decision: This is a critical section protected by mutex to ensure
// the signature_counter is strictly monotonically increasing
func (d *SignatureDevice) UpdateAfterSigning(newSignature string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.LastSignature = newSignature
	d.SignatureCounter++
}

// ToPublicView returns a view of the device without sensitive information
// Design Decision: Separates internal state from API responses
func (d *SignatureDevice) ToPublicView() *SignatureDeviceView {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return &SignatureDeviceView{
		ID:               d.ID,
		Algorithm:        d.Algorithm,
		Label:            d.Label,
		SignatureCounter: d.SignatureCounter,
		LastSignature:    d.LastSignature,
	}
}

// SignatureDeviceView represents the public view of a device (without private keys)
type SignatureDeviceView struct {
	ID               string             `json:"id"`
	Algorithm        SignatureAlgorithm `json:"algorithm"`
	Label            string             `json:"label,omitempty"`
	SignatureCounter int                `json:"signature_counter"`
	LastSignature    string             `json:"last_signature"`
}
