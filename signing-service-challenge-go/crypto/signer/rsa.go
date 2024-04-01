package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// RSA signing
type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		PrivateKey: privateKey,
	}
}

func (r *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}
