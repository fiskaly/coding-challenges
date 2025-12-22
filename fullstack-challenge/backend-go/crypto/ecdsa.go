package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// ECCKeyPair represents an ECC key pair.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// Sign signs the given data with the ECC private key.
func (kp *ECCKeyPair) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return ecdsa.SignASN1(rand.Reader, kp.Private, hash[:])
}

// EncodePEM encodes the ECC key pair to PEM format.
func (kp *ECCKeyPair) EncodePEM() (publicPEM, privatePEM string, err error) {
	// Encode public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(kp.Public)
	if err != nil {
		return "", "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	// Encode private key
	privKeyBytes, err := x509.MarshalECPrivateKey(kp.Private)
	if err != nil {
		return "", "", err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	return string(pubPEM), string(privPEM), nil
}

// DecodeECCPEM decodes an ECC key pair from PEM format.
func DecodeECCPEM(publicPEM, privatePEM string) (*ECCKeyPair, error) {
	// Decode public key
	pubBlock, _ := pem.Decode([]byte(publicPEM))
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Decode private key
	privBlock, _ := pem.Decode([]byte(privatePEM))
	privKey, err := x509.ParseECPrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  pubKey.(*ecdsa.PublicKey),
		Private: privKey,
	}, nil
}
