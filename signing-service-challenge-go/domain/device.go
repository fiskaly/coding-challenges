package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Defines the algorithm related functions that `domain` package requires.
// These operations will be implemented by algorithm specific structs in the
// `crypto` package.
// e.g. `RSAAlgorithm`, `ECCAlgorithm`
type SignatureAlgorithm interface {
	Name() string
	GenerateEncodedPrivateKey() ([]byte, error)
}

type SignatureDevice struct {
	ID                uuid.UUID
	AlgorithmName     string
	EncodedPrivateKey []byte
	// (optional) user provided string to be displayed in the UI
	Label string
	// track the last signature created with this device
	LastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func BuildSignatureDevice(id uuid.UUID, algorithm SignatureAlgorithm, label ...string) (SignatureDevice, error) {
	encodedPrivateKey, err := algorithm.GenerateEncodedPrivateKey()
	if err != nil {
		err = errors.New(fmt.Sprintf("private key generation failed: %s", err.Error()))
		return SignatureDevice{}, err
	}

	device := SignatureDevice{
		ID:                id,
		AlgorithmName:     algorithm.Name(),
		EncodedPrivateKey: encodedPrivateKey,
	}

	if len(label) > 0 {
		device.Label = label[0]
	}

	return device, nil
}

type SignatureDeviceRepository interface {
	Create(device SignatureDevice) error
	Find(id uuid.UUID) (SignatureDevice, bool, error)
}
