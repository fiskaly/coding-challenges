package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RSA struct{}

// Generate creates a new RSA key pair.
func (r *RSA) Generate() (*KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: key,
		Public:  &key.PublicKey,
	}, nil
}

func (r *RSA) Sign(pk crypto.PrivateKey, dataToBeSigned []byte) ([]byte, error) {
	// type assertion to convert pk to *rsa.PrivateKey
	rsaKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid key type for RSA")
	}
	hashed := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (r *RSA) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	// type assertion to convert pk to *rsa.PrivateKey
	rsaKeyPair, ok := keyPair.Private.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid key type for RSA")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(rsaKeyPair)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&rsaKeyPair.PublicKey)

	return EncodePEM("RSA_PRIVATE_KEY", privateKeyBytes), EncodePEM("RSA_PUBLIC_KEY", publicKeyBytes), nil
}

func (r *RSA) Unmarshal(privateKeyBytes []byte) (*KeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
