package marshaler

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	h "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto/helpers"
)

type ECCMarshaler struct{}

func NewECCMarshaler() *ECCMarshaler {
	return &ECCMarshaler{}
}

func (m *ECCMarshaler) Marshal(keyPair c.KeyPair) ([]byte, []byte, error) {
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

	return h.EncodePEM("EC_PRIVATE_KEY", privateKeyBytes), h.EncodePEM("ECC_PUBLIC_KEY", publicKeyBytes), nil
}

func (m *ECCMarshaler) Unmarshal(privateKeyBytes []byte) (*c.KeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &c.KeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}
