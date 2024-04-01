package generator

import (
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func init() {
	registerGenerator("RSA", &RSAKeyGenerator{})
	registerGenerator("ECC", &ECCKeyGenerator{})
}

// KeyPairGenerator interface for generating key pairs.
type KeyPairGenerator interface {
	Generate() (*crypto.KeyPair, error)
}

// generators is used as a registry map
var generators = make(map[string]KeyPairGenerator)

func registerGenerator(algorithm string, gen KeyPairGenerator) {
	generators[algorithm] = gen
}

func GetGenerator(algorithm string) (KeyPairGenerator, error) {
	gen, exists := generators[algorithm]
	if !exists {
		return nil, fmt.Errorf("no generator exists for algorithm: %s", algorithm)
	}

	return gen, nil
}
