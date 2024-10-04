package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type ECCGenerator interface {
	Generate() (*ECCKeyPair, error)
}

// ECCGenerator generates an ECC key pair.
type DefaultECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *DefaultECCGenerator) Generate() (*ECCKeyPair, error) {
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
