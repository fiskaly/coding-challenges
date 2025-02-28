package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Signer interface {
	Sign(data []byte) (string, error)
}

type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

func NewRSASigner() (*RSASigner, error) {
	adapter, err := NewGeneratorAdapter("RSA")
	if err != nil {
		return nil, err
	}

	keyPair, err := adapter.Generate()
	if err != nil {
		return nil, err
	}

	privateKey, ok := keyPair.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid RSA private key type")
	}

	return &RSASigner{PrivateKey: privateKey}, nil
}

func (s *RSASigner) Sign(data []byte) (string, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(nil, s.PrivateKey, 0, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

type ECCSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewECCSigner() (*ECCSigner, error) {
	adapter, err := NewGeneratorAdapter("ECC")
	if err != nil {
		return nil, err
	}

	keyPair, err := adapter.Generate()
	if err != nil {
		return nil, err
	}

	privateKey, ok := keyPair.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid ECC private key type")
	}

	return &ECCSigner{PrivateKey: privateKey}, nil
}

func (s *ECCSigner) Sign(data []byte) (string, error) {
	hashed := sha256.Sum256(data)
	r, sSig, err := ecdsa.Sign(nil, s.PrivateKey, hashed[:])
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), sSig.Bytes()...)
	return base64.StdEncoding.EncodeToString(signature), nil
}

func NewSignerFactory(algorithm string) (Signer, error) {
	switch algorithm {
	case "RSA":
		return NewRSASigner()
	case "ECC":
		return NewECCSigner()
	default:
		return nil, errors.New("unsupported signing algorithm")
	}
}

func NewSignerWithKey(algorithm string, privateKey any) (Signer, error) {
	switch algorithm {
	case "RSA":
		rsaKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("invalid RSA private key")
		}
		return &RSASigner{PrivateKey: rsaKey}, nil
	case "ECC":
		eccKey, ok := privateKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("invalid ECC private key")
		}
		return &ECCSigner{PrivateKey: eccKey}, nil
	default:
		return nil, errors.New("unsupported signing algorithm")
	}
}
