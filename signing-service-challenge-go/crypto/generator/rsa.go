package generator

import (
	"crypto/rand"
	"crypto/rsa"

	c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// RSAKeyGeneratorD
type RSAKeyGenerator struct{}

// Generate creates a new RSA key pair.
func (g *RSAKeyGenerator) Generate() (*c.KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &c.KeyPair{
		Private: key,
		Public:  &key.PublicKey,
	}, nil
}
