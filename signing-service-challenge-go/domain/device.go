package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// Algorithm represents a supported signing algorithm.
type Algorithm string

const (
	AlgorithmRSA Algorithm = "RSA"
	AlgorithmECC Algorithm = "ECC"
)

// ErrInvalidDeviceID is returned when a device is created with an empty ID.
var ErrInvalidDeviceID = errors.New("domain: signature device id must be non-empty")

// ErrUnsupportedAlgorithm is returned when the provided algorithm is not supported.
var ErrUnsupportedAlgorithm = errors.New("domain: unsupported signing algorithm")

// ErrEmptyPayload is returned when the provided payload is empty.
var ErrEmptyPayload = errors.New("domain: data to be signed must be non-empty")

// ErrEmptySignature is returned when a signature update is attempted with empty data.
var ErrEmptySignature = errors.New("domain: signature cannot be empty")

// ParseAlgorithm converts a string into an Algorithm and validates it.
func ParseAlgorithm(raw string) (Algorithm, error) {
	algo := Algorithm(strings.ToUpper(strings.TrimSpace(raw)))
	if algo.IsSupported() {
		return algo, nil
	}

	return "", ErrUnsupportedAlgorithm
}

// IsSupported reports whether the algorithm is supported.
func (a Algorithm) IsSupported() bool {
	switch a {
	case AlgorithmRSA, AlgorithmECC:
		return true
	default:
		return false
	}
}

// String returns the underlying string representation.
func (a Algorithm) String() string {
	return string(a)
}

// SignatureDevice models the domain behaviour of a device capable of producing signatures.
type SignatureDevice struct {
	id               string
	label            string
	algorithm        Algorithm
	signatureCounter uint64
	lastSignature    string // always base64 encoded
}

// NewSignatureDevice constructs a new signature device with the supplied metadata.
// A newly created device starts with counter 0 and lastSignature = base64(deviceID).
func NewSignatureDevice(id string, algorithm Algorithm, label string) (*SignatureDevice, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidDeviceID
	}

	if !algorithm.IsSupported() {
		return nil, ErrUnsupportedAlgorithm
	}

	return &SignatureDevice{
		id:               id,
		label:            label,
		algorithm:        algorithm,
		signatureCounter: 0,
		lastSignature:    base64.StdEncoding.EncodeToString([]byte(id)),
	}, nil
}

// RestoreSignatureDevice recreates a SignatureDevice from persisted state.
func RestoreSignatureDevice(id string, algorithm Algorithm, label string, signatureCounter uint64, lastSignature string) (*SignatureDevice, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidDeviceID
	}

	if !algorithm.IsSupported() {
		return nil, ErrUnsupportedAlgorithm
	}

	if strings.TrimSpace(lastSignature) == "" {
		lastSignature = base64.StdEncoding.EncodeToString([]byte(id))
	}

	return &SignatureDevice{
		id:               id,
		label:            label,
		algorithm:        algorithm,
		signatureCounter: signatureCounter,
		lastSignature:    lastSignature,
	}, nil
}

// ID returns the unique identifier of the device.
func (d *SignatureDevice) ID() string {
	return d.id
}

// Label returns the user-friendly label of the device.
func (d *SignatureDevice) Label() string {
	return d.label
}

// SetLabel updates the device label.
func (d *SignatureDevice) SetLabel(label string) {
	d.label = label
}

// Algorithm returns the configured signing algorithm.
func (d *SignatureDevice) Algorithm() Algorithm {
	return d.algorithm
}

// SignatureCounter returns the current signature counter.
func (d *SignatureDevice) SignatureCounter() uint64 {
	return d.signatureCounter
}

// LastSignature returns the base64 encoded representation of the last signature.
func (d *SignatureDevice) LastSignature() string {
	return d.lastSignature
}

// SecuredPayload assembles the payload that has to be signed following the format
// "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>".
func (d *SignatureDevice) SecuredPayload(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", ErrEmptyPayload
	}

	return fmt.Sprintf("%d_%s_%s", d.signatureCounter, data, d.lastSignature), nil
}

// RecordSignature updates the internal state after a successful signing operation.
// It encodes the signature to base64, stores it as the new lastSignature, and increments the counter.
func (d *SignatureDevice) RecordSignature(signature []byte) (string, error) {
	if len(signature) == 0 {
		return "", ErrEmptySignature
	}

	encoded := base64.StdEncoding.EncodeToString(signature)
	d.signatureCounter++
	d.lastSignature = encoded

	return encoded, nil
}
