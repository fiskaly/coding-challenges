package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// ErrNilKeyPair is returned when a signer is constructed with a nil key pair.
var ErrNilKeyPair = errors.New("crypto: key pair must not be nil")

// ErrNilPrivateKey is returned when the key pair does not contain a private key.
var ErrNilPrivateKey = errors.New("crypto: private key must not be nil")

// RSASigner implements Signer using RSA PKCS#1 v1.5 with SHA-256.
type RSASigner struct {
	privateKey *rsa.PrivateKey
}

// NewRSASigner constructs an RSASigner for the provided key pair.
func NewRSASigner(keyPair *RSAKeyPair) (*RSASigner, error) {
	if keyPair == nil {
		return nil, ErrNilKeyPair
	}

	if keyPair.Private == nil {
		return nil, ErrNilPrivateKey
	}

	return &RSASigner{
		privateKey: keyPair.Private,
	}, nil
}

// Sign signs the provided data using RSA PKCS#1 v1.5 and returns the raw signature bytes.
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	if s == nil || s.privateKey == nil {
		return nil, ErrNilPrivateKey
	}

	hashed := sha256.Sum256(dataToBeSigned)
	return rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
}

// ECDSASigner implements Signer using ECDSA with SHA-384 digest.
type ECDSASigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewECDSASigner constructs an ECDSASigner for the provided key pair.
func NewECDSASigner(keyPair *ECCKeyPair) (*ECDSASigner, error) {
	if keyPair == nil {
		return nil, ErrNilKeyPair
	}

	if keyPair.Private == nil {
		return nil, ErrNilPrivateKey
	}

	return &ECDSASigner{
		privateKey: keyPair.Private,
	}, nil
}

// Sign signs the provided data using ECDSA and returns the ASN.1 encoded signature.
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	if s == nil || s.privateKey == nil {
		return nil, ErrNilPrivateKey
	}

	hasher := sha512.New384()
	if _, err := hasher.Write(dataToBeSigned); err != nil {
		return nil, fmt.Errorf("crypto: hashing data: %w", err)
	}

	digest := hasher.Sum(nil)
	return ecdsa.SignASN1(rand.Reader, s.privateKey, digest)
}
