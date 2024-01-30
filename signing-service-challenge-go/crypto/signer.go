package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type SignerFactoryImpl struct{}

func NewSignerFactoryImpl() *SignerFactoryImpl {
	return &SignerFactoryImpl{}
}

func (f *SignerFactoryImpl) CreateSigner(algorithm domain.CryptoAlgorithm) domain.Signer {
	switch algorithm {
	case domain.RSA:
		return NewRSASigner()
	case domain.ECC:
		return NewECCSigner()
	default:
		//todo: Handle unsupported algorithm
		return nil
	}
}

// implement RSA and ECDSA signing ...
type RSASigner struct{}

func NewRSASigner() *RSASigner {
	return &RSASigner{}
}
func (g *RSASigner) Sign(privateKey string, data string) (signeddata string, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

type ECCSigner struct{}

func NewECCSigner() *ECCSigner {
	return &ECCSigner{}
}
func (g *ECCSigner) Sign(privateKey string, data string) (signeddata string, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	ecPrivateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	hashed := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, ecPrivateKey, hashed[:])
	if err != nil {
		return "", err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return base64.StdEncoding.EncodeToString(signature), nil
}
