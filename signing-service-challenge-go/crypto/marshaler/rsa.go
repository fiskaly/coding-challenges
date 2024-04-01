package marshaler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	h "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto/helpers"
)

type RSAMarshaler struct{}

func NewRSAMarshaler() *RSAMarshaler {
	return &RSAMarshaler{}
}

func (m *RSAMarshaler) Marshal(keyPair c.KeyPair) ([]byte, []byte, error) {
	rsaKeyPair, ok := keyPair.Private.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid key type for RSA")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(rsaKeyPair)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&rsaKeyPair.PublicKey)

	return h.EncodePEM("RSA_PRIVATE_KEY", privateKeyBytes), h.EncodePEM("RSA_PUBLIC_KEY", publicKeyBytes), nil
}

func (m *RSAMarshaler) Unmarshal(privateKeyBytes []byte) (*c.KeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &c.KeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
