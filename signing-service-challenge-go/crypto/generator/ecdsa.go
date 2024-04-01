package generator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// ECCKeyGenerator
type ECCKeyGenerator struct{}

// Generate creates a new ECC key pair.
func (g *ECCKeyGenerator) Generate() (*c.KeyPair, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader) // Using P-384 curve
	if err != nil {
		return nil, err
	}

	return &c.KeyPair{
		Private: key,
		Public:  &key.PublicKey,
	}, nil
}
