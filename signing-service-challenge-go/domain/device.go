package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type SignatureDevice struct {
	Id     uuid.UUID
	Label  *string
	signer crypto.Signer
}

func NewSignatureDevice(id uuid.UUID, label *string, algorithm string) (*SignatureDevice, error) {
	signer, err := crypto.NewSigner(algorithm)
	if err != nil {
		return nil, err
	}

	return &SignatureDevice{
		Id:     id,
		Label:  label,
		signer: signer,
	}, nil
}
