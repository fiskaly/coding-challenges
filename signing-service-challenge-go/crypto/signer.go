package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	marshaller      RSAMarshaler
	privateKeyBytes []byte
}

func NewRSASigner(privateKey []byte) RSASigner {
	return RSASigner{
		marshaller:      NewRSAMarshaler(),
		privateKeyBytes: privateKey,
	}
}

func (signer RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash, err := getHashSum(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	keyPair, err := signer.marshaller.Unmarshal(signer.privateKeyBytes)
	if err != nil {
		return nil, err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, keyPair.Private, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hash, signature)
	if err == nil {
		return nil, err
	}
	return signature, nil
}

type ECCSigner struct {
	marshaller      ECCMarshaler
	privateKeyBytes []byte
}

func NewECCSigner(privateKey []byte) ECCSigner {
	return ECCSigner{
		marshaller:      NewECCMarshaler(),
		privateKeyBytes: privateKey,
	}
}

func (signer ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash, err := getHashSum(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	keyPair, err := signer.marshaller.Unmarshal(signer.privateKeyBytes)
	if err != nil {
		return nil, err
	}
	signature, err := ecdsa.SignASN1(rand.Reader, keyPair.Private, hash[:])
	if err != nil {
		return nil, err
	}
	if !ecdsa.VerifyASN1(keyPair.Public, hash, signature) {
		return nil, errors.New("Failed to verify ASN1 signature")
	}
	return signature, nil
}

func getHashSum(dataToBeSigned []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(dataToBeSigned)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash sum: %w", err)
	}
	return msgHash.Sum(nil), nil
}
