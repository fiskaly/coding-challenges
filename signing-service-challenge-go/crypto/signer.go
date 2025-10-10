package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"math/big"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// RSASigner implements the Signer interface for RSA algorithm
// Design Decision: Using SHA256 as the hash function for RSA-PSS
// PSS (Probabilistic Signature Scheme) is more secure than PKCS#1 v1.5
type RSASigner struct {
	keyPair *RSAKeyPair
}

// NewRSASigner creates a new RSA signer from an RSAKeyPair
func NewRSASigner(keyPair *RSAKeyPair) (*RSASigner, error) {
	if keyPair == nil || keyPair.Private == nil {
		return nil, errors.New("RSA key pair and private key cannot be nil")
	}
	return &RSASigner{
		keyPair: keyPair,
	}, nil
}

// Sign signs the data using RSA-PSS
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data
	hashed := sha256.Sum256(dataToBeSigned)

	// Sign using RSA-PSS
	signature, err := rsa.SignPSS(
		rand.Reader,
		s.keyPair.Private,
		crypto.SHA256,
		hashed[:],
		nil,
	)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// ECDSASigner implements the Signer interface for ECDSA algorithm
// Design Decision: Using SHA256 as the hash function for ECDSA
// The signature is encoded as r||s (concatenation of r and s values)
type ECDSASigner struct {
	keyPair *ECCKeyPair
}

// NewECDSASigner creates a new ECDSA signer from an ECCKeyPair
func NewECDSASigner(keyPair *ECCKeyPair) (*ECDSASigner, error) {
	if keyPair == nil || keyPair.Private == nil {
		return nil, errors.New("ECC key pair and private key cannot be nil")
	}
	return &ECDSASigner{
		keyPair: keyPair,
	}, nil
}

// Sign signs the data using ECDSA
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data
	hashed := sha256.Sum256(dataToBeSigned)

	// Sign using ECDSA
	r, sigS, err := ecdsa.Sign(rand.Reader, s.keyPair.Private, hashed[:])
	if err != nil {
		return nil, err
	}

	// Encode signature as r||s
	// Get the curve's byte size
	curveOrderByteSize := (s.keyPair.Private.Curve.Params().BitSize + 7) / 8

	// Allocate buffer for signature
	signature := make([]byte, 2*curveOrderByteSize)

	// Fill r and s into the buffer
	r.FillBytes(signature[0:curveOrderByteSize])
	sigS.FillBytes(signature[curveOrderByteSize:])

	return signature, nil
}

// VerifyRSA verifies an RSA signature (useful for testing)
func VerifyRSA(keyPair *RSAKeyPair, data, signature []byte) error {
	if keyPair == nil || keyPair.Public == nil {
		return errors.New("RSA key pair and public key cannot be nil")
	}
	hashed := sha256.Sum256(data)
	return rsa.VerifyPSS(keyPair.Public, crypto.SHA256, hashed[:], signature, nil)
}

// VerifyECDSA verifies an ECDSA signature (useful for testing)
func VerifyECDSA(keyPair *ECCKeyPair, data, signature []byte) error {
	if keyPair == nil || keyPair.Public == nil {
		return errors.New("ECC key pair and public key cannot be nil")
	}

	hashed := sha256.Sum256(data)

	// Decode r and s from signature
	curveOrderByteSize := (keyPair.Public.Curve.Params().BitSize + 7) / 8

	if len(signature) != 2*curveOrderByteSize {
		return errors.New("invalid signature length")
	}

	r := new(big.Int).SetBytes(signature[0:curveOrderByteSize])
	s := new(big.Int).SetBytes(signature[curveOrderByteSize:])

	if !ecdsa.Verify(keyPair.Public, hashed[:], r, s) {
		return errors.New("signature verification failed")
	}

	return nil
}
