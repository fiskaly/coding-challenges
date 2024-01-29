package domain

import "errors"

type CryptoAlgorithm struct {
	value string
}

func (r CryptoAlgorithm) String() string {
	return r.value
}

var (
	Unknown = CryptoAlgorithm{""}
	ECC     = CryptoAlgorithm{"ECC"}
	RSA     = CryptoAlgorithm{"RSA"}
)

func FromString(s string) (CryptoAlgorithm, error) {
	switch s {
	case ECC.value:
		return ECC, nil
	case RSA.value:
		return RSA, nil
	default:
		return Unknown, errors.New("unknown algorithm: " + s)
	}

}

type KeyPairFactory interface {
	CreateKeyPairGenerator(algorithm CryptoAlgorithm) KeyPairGenerator
}

type KeyPairGenerator interface {
	GenerateKeyPair() (publicKey string, privateKey string, err error)
}
