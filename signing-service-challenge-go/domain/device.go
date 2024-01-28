package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type KeyPair interface {
	Sign(dataToBeSigned []byte) (signature []byte, err error)
}

type KeyPairGenerator interface {
	AlgorithmName() string
	Generate() (KeyPair, error)
}

type SignatureDevice struct {
	ID      uuid.UUID
	KeyPair KeyPair
	// (optional) user provided string to be displayed in the UI
	Label string
	// track the last signature created with this device
	Base64EncodedLastSignature string
	// track how many signatures have been created with this device
	SignatureCounter uint
}

func (device SignatureDevice) Sign(dataToBeSigned string) ([]byte, error) {
	return device.KeyPair.Sign([]byte(dataToBeSigned))
}

func BuildSignatureDevice(id uuid.UUID, generator KeyPairGenerator, label ...string) (SignatureDevice, error) {
	keyPair, err := generator.Generate()
	if err != nil {
		err = errors.New(fmt.Sprintf("key pair generation failed: %s", err.Error()))
		return SignatureDevice{}, err
	}

	device := SignatureDevice{
		ID:      id,
		KeyPair: keyPair,
	}

	if len(label) > 0 {
		device.Label = label[0]
	}

	return device, nil
}

type SignatureDeviceRepository interface {
	Create(device SignatureDevice) error
	Update(device SignatureDevice) error
	Find(id uuid.UUID) (SignatureDevice, bool, error)
}
