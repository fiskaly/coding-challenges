package signer

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

// ECDSA signing
type ECDSASigner struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewECDSASigner(privateKey *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{
		PrivateKey: privateKey,
	}
}

func (e *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	r, s, err := ecdsa.Sign(rand.Reader, e.PrivateKey, hashed[:])
	if err != nil {
		return nil, err
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	return signature, nil
}
