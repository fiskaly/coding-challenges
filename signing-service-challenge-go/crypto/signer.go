package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	privateKey rsa.PrivateKey
}

func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	return rsa.SignPKCS1v15(rand.Reader, &s.privateKey, crypto.SHA256, hashed[:])
}

type ECCSigner struct {
	privateKey ecdsa.PrivateKey
}

func (s *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	return ecdsa.SignASN1(rand.Reader, &s.privateKey, hashed[:])
}

func MakeSigner(algorithm string, privateKey []byte) (Signer, error) {
	switch algorithm {
	case "RSA":
		marshaler := RSAMarshaler{}
		kp, err := marshaler.Unmarshal(privateKey)
		if err != nil {
			return nil, fmt.Errorf("unmarshal for RSA private key failed: %w", err)
		}
		return &RSASigner{
			privateKey: *kp.Private,
		}, nil
	case "ECC":
		marshaler := NewECCMarshaler()
		kp, err := marshaler.Decode(privateKey)
		if err != nil {
			return nil, fmt.Errorf("unmarshal for ECC private key failed: %w", err)
		}
		return &ECCSigner{
			privateKey: *kp.Private,
		}, nil
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}
