package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"

	"crypto/x509"
	"encoding/pem"
)

type KeyPairFactoryImpl struct{}

func NewKeyPairFactoryImpl() *KeyPairFactoryImpl {
	return &KeyPairFactoryImpl{}
}

func (f *KeyPairFactoryImpl) CreateKeyPairGenerator(algorithm domain.CryptoAlgorithm) domain.KeyPairGenerator {
	switch algorithm {
	case domain.RSA:
		return NewRSAKeyPairGenerator()
	case domain.ECC:
		return NewECDSAKeyPairGenerator()
	default:
		//todo: Handle unsupported algorithm
		return nil
	}
}

// // RSAGenerator generates a RSA key pair.
// type RSAGenerator struct{}

// // Generate generates a new RSAKeyPair.
// func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
// 	// Security has been ignored for the sake of simplicity.
// 	key, err := rsa.GenerateKey(rand.Reader, 512)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &RSAKeyPair{
// 		Public:  &key.PublicKey,
// 		Private: key,
// 	}, nil
// }

type RSAKeyPairGenerator struct{}

func NewRSAKeyPairGenerator() *RSAKeyPairGenerator {
	return &RSAKeyPairGenerator{}
}

func (g *RSAKeyPairGenerator) GenerateKeyPair() (publicKey string, privateKey string, err error) {

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), string(privKeyPEM), nil

}

// // ECCGenerator generates an ECC key pair.
// type ECCGenerator struct{}

// // Generate generates a new ECCKeyPair.
// func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
// 	// Security has been ignored for the sake of simplicity.
// 	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ECCKeyPair{
// 		Public:  &key.PublicKey,
// 		Private: key,
// 	}, nil
// }

type ECDSAKeyPairGenerator struct{}

func NewECDSAKeyPairGenerator() *ECDSAKeyPairGenerator {
	return &ECDSAKeyPairGenerator{}
}

func (g *ECDSAKeyPairGenerator) GenerateKeyPair() (publicKey string, privateKey string, err error) {
	// Generate ECDSA key pair
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// Encode private key to PEM format
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", "", err
	}
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	// Encode public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), string(privKeyPEM), nil
}
