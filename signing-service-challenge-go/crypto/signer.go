package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

const (
	AlgorithmRSA = "RSA"
	AlgorithmECC = "ECC"
)

func NewSigner(algorithm string) (Signer, error) {
	switch algorithm {
	case AlgorithmRSA:
		return &RSASigner{}, nil
	case AlgorithmECC:
		return &ECCSigner{}, nil
	default:
		return nil, errors.New("invalid algorithm")
	}
}

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
}

func (signer *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// TODO: Key pair shall be provided to signer in constructor
	keyPair, err := generateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(dataToBeSigned)
	sig, err := rsa.SignPKCS1v15(nil, keyPair.Private, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	// TODO: Extract into a Verifier interface
	if err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hashed[:], sig); err != nil {
		return nil, err
	}

	return sig, nil
}

// TODO: inject generator in constructor
func generateRSAKeyPair() (*RSAKeyPair, error) {
	g := &RSAGenerator{}
	return g.Generate()
}

type ECCSigner struct {
}

func (signer *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// TODO: Key pair shall be provided to signer in constructor
	keyPair, err := generateECCKeyPair()
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(dataToBeSigned)
	sig, err := ecdsa.SignASN1(rand.Reader, keyPair.Private, hashed[:])
	if err != nil {
		return nil, err
	}

	// TODO: Extract into a Verifier interface
	if valid := ecdsa.VerifyASN1(keyPair.Public, hashed[:], sig); !valid {
		return nil, errors.New("invalid ECC signature")
	}

	return sig, nil
}

func generateECCKeyPair() (*ECCKeyPair, error) {
	g := &ECCGenerator{}
	return g.Generate()
}
