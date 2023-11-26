package crypto

import (
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

// TODO: implement RSA and ECDSA signing ...
type RSASigner struct {
}

func (signer *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

type ECCSigner struct {
}

func (signer *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}
