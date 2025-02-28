package crypto

import (
	"errors"
)

type AdaptedKeyPair struct {
	PublicKey  any
	PrivateKey any
}

type GeneratorAdapter struct {
	algorithm string
}

func NewGeneratorAdapter(algorithm string) (*GeneratorAdapter, error) {
	switch algorithm {
	case "RSA":
		return &GeneratorAdapter{algorithm: "RSA"}, nil
	case "ECC":
		return &GeneratorAdapter{algorithm: "ECC"}, nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

func (ga *GeneratorAdapter) Generate() (AdaptedKeyPair, error) {
	switch ga.algorithm {
	case "RSA":
		rsaGen := &RSAGenerator{}
		keyPair, err := rsaGen.Generate()
		if err != nil {
			return AdaptedKeyPair{}, err
		}
		return AdaptedKeyPair{
			PublicKey:  keyPair.Public,
			PrivateKey: keyPair.Private,
		}, nil
	case "ECC":
		eccGen := &ECCGenerator{}
		keyPair, err := eccGen.Generate()
		if err != nil {
			return AdaptedKeyPair{}, err
		}
		return AdaptedKeyPair{
			PublicKey:  keyPair.Public,
			PrivateKey: keyPair.Private,
		}, nil
	default:
		return AdaptedKeyPair{}, errors.New("unsupported algorithm")
	}
}
