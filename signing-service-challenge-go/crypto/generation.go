package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

var ErrUnsupportedAlgorithm = fmt.Errorf("unsupported algorithm")

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

func GenerateAndEncode(algorithm string) ([]byte, []byte, error) {
	switch algorithm {
	case "RSA":
		generator := &RSAGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
		}
		marshaler := NewRSAMarshaler()
		return marshaler.Marshal(*keyPair)
	case "ECC":
		generator := &ECCGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate ECC key pair: %w", err)
		}
		marshaler := NewECCMarshaler()
		return marshaler.Encode(*keyPair)

	default:
		return nil, nil, ErrUnsupportedAlgorithm
	}
}
