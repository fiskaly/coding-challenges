package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type ECCSigner struct {
	privateKey *ecdsa.PrivateKey
}

func NewECCSigner(privateKey *ecdsa.PrivateKey) *ECCSigner {
	return &ECCSigner{privateKey: privateKey}
}

// TODO: implement RSA and ECDSA signing ...
func (s *ECCSigner) Sign(data []byte) ([]byte, error) {
	hashedData := sha256.Sum256(data)
	r, sValue, err := ecdsa.Sign(rand.Reader, s.privateKey, hashedData[:])
	if err != nil {
		return nil, NewSignOperationError("ECDSA", err)
	}

	signedData, err := asn1.Marshal(struct{ R, S *big.Int }{R: r, S: sValue})
	if err != nil {
		return nil, NewMarshalError("ECDSA", err)
	}

	return signedData, nil
}

type RSASigner struct {
	privateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{privateKey: privateKey}
}

func (s *RSASigner) Sign(data []byte) ([]byte, error) {
	hashedData := sha256.Sum256(data)
	signedData, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashedData[:])
	if err != nil {
		return nil, NewSignOperationError("RSA", err)
	}

	return signedData, nil
}
