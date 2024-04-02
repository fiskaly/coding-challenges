package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type ECDSA struct{}

// Generate - creates a new ECC key pair.
func (ecc *ECDSA) Generate() (*KeyPair, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: key,
		Public:  &key.PublicKey,
	}, nil
}

func (e *ECDSA) Sign(pk crypto.PrivateKey, dataToBeSigned []byte) ([]byte, error) {
	// type assertion to convert pk to *ecdsa.PrivateKey
	ecdsaKey, ok := pk.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid key type for ECC")
	}

	hashed := sha256.Sum256(dataToBeSigned)
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey, hashed[:])
	if err != nil {
		return nil, err
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	return signature, nil
}

func (m *ECDSA) Marshal(keyPair KeyPair) ([]byte, []byte, error) {
	// type assertion to convert pk to *ecdsa.PrivateKey
	eccKeyPair, ok := keyPair.Private.(*ecdsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid key type for ECC")
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(eccKeyPair)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(eccKeyPair.Public())
	if err != nil {
		return nil, nil, err
	}

	return EncodePEM("EC_PRIVATE_KEY", privateKeyBytes), EncodePEM("ECC_PUBLIC_KEY", publicKeyBytes), nil
}

func (m *ECDSA) Unmarshal(privateKeyBytes []byte) (*KeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
