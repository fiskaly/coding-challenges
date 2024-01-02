package domain

import (
	"errors"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"

	"github.com/google/uuid"
)

// TODO: signature device domain model ...
type SignatureDevice struct {
	ID               uuid.UUID
	Algorithm        string
	SignatureCounter int
	LastSignature    string
	PublicKey        interface{}
	PrivateKey       interface{}
	Label            string
}

func (signDevice *SignatureDevice) AddDevice() error {

	switch signDevice.Algorithm {
	case "RSA":
		rsaGenerator := &crypto.RSAGenerator{}

		rsaKeyPair, err := rsaGenerator.Generate()
		if err != nil {
			fmt.Println("Error generating RSA key pair:", err)
			return err
		}
		signDevice.PrivateKey = rsaKeyPair.Private
		signDevice.PublicKey = rsaKeyPair.Public

	case "ECC":
		eccGenerator := &crypto.ECCGenerator{}
		eccKeyPair, err := eccGenerator.Generate()
		if err != nil {
			fmt.Println("Error generating ECC key pair:", err)
			return err
		}
		signDevice.PrivateKey = eccKeyPair.Private
		signDevice.PublicKey = eccKeyPair.Public

	default:
		return errors.New("unsupported algorithm")
	}

	return nil
}
