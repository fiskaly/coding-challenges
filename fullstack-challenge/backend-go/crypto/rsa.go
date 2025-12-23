package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// RSAKeyPair represents a RSA key pair.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// Sign signs the given data with the RSA private key.
func (kp *RSAKeyPair) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, kp.Private, crypto.SHA256, hash[:])
}

// EncodePEM encodes the RSA key pair to PEM format.
func (kp *RSAKeyPair) EncodePEM() (publicPEM, privatePEM string, err error) {
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
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(kp.Private),
	})

	return string(pubPEM), string(privPEM), nil
}

// DecodeRSAPEM decodes a RSA key pair from PEM format.
func DecodeRSAPEM(publicPEM, privatePEM string) (*RSAKeyPair, error) {
	// Decode public key
	pubBlock, _ := pem.Decode([]byte(publicPEM))
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Decode private key
	privBlock, _ := pem.Decode([]byte(privatePEM))
	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  pubKey.(*rsa.PublicKey),
		Private: privKey,
	}, nil
}
