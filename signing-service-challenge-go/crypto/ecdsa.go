package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
)

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() ECCMarshaler {
	return ECCMarshaler{}
}

// Encode takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m ECCMarshaler) Encode(keyPair ECCKeyPair) (encodedPublicKey, encodedPrivateKey []byte, err error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(keyPair.Private)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Decode assembles an ECCKeyPair from an encoded private key.
func (m ECCMarshaler) Decode(privateKeyBytes []byte) (*ECCKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

// Implements domain.SignatureAlgorithm for RSA.
// Note that any actual logic is implemented in `ECCSigner`, `ECCMarshaller` and `ECCGenerator`,
// and this struct merely acts as a facade to make this logic easier to access in the
// `domain` package.
type ECCAlgorithm struct{}

func (ecc ECCAlgorithm) Name() string {
	return "ECC"
}

func (ecc ECCAlgorithm) GenerateEncodedPrivateKey() ([]byte, error) {
	generator := ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	marshaller := NewECCMarshaler()
	_, privateKey, err := marshaller.Encode(*keyPair)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (ecc ECCAlgorithm) SignTransaction(encodedPrivateKey []byte, dataToBeSigned []byte) ([]byte, error) {
	marshaller := NewECCMarshaler()
	keyPair, err := marshaller.Decode(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	signer := ECCSigner{keyPair: *keyPair}
	return signer.Sign(dataToBeSigned)
}
