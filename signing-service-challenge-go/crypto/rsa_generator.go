package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

type RSAGenerator interface {
	Generate() (*RSAKeyPair, error)
}

// RSAGenerator generates a RSA key pair.
type DefaultRSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *DefaultRSAGenerator) Generate() (*RSAKeyPair, error) {
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
