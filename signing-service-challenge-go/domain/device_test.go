package domain

import (
	"testing"

	"github.com/google/uuid"
)

type MockKeyPair struct{}

func (keyPair MockKeyPair) Sign(dataToBeSigned []byte) (signature []byte, err error) {
	return nil, nil
}

type MockKeyPairGenerator struct{}

func (g MockKeyPairGenerator) AlgorithmName() string {
	return ""
}

func (g MockKeyPairGenerator) Generate() (KeyPair, error) {
	return MockKeyPair{}, nil
}

func TestBuildSignatureDevice(t *testing.T) {
	t.Run("successfully builds signature device", func(t *testing.T) {
		id := uuid.New()
		device, err := BuildSignatureDevice(id, MockKeyPairGenerator{})

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.ID != id {
			t.Errorf("expected id: %s, got: %s", id, device.ID.String())
		}

		if device.KeyPair == nil {
			t.Error("expected key pair to be set, got nil")
		}

		if device.SignatureCounter != 0 {
			t.Errorf("expected initial signature counter value to be 0, got: %d", device.SignatureCounter)
		}

		if device.Base64EncodedLastSignature != "" {
			t.Errorf("expected initial last signature value to be blank, got: %s", device.Base64EncodedLastSignature)
		}

		if device.Label != "" {
			t.Errorf("expected label be blank when not provided, got: %s", device.Label)
		}
	})

	t.Run("sets label when provided", func(t *testing.T) {
		id := uuid.New()
		label := "some-label"
		device, err := BuildSignatureDevice(
			id,
			MockKeyPairGenerator{},
			"some-label",
		)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if device.Label != label {
			t.Errorf("expected label: %s, got: %s", label, device.Label)
		}
	})
}
