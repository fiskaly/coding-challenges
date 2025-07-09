package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		err := fmt.Errorf("error while generating the RSA key pair: %s", err)
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("RSA key pairs succesfully generated")

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		err := fmt.Errorf("error while generating the ECC key pair: %s", err)
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("ECC key pairs succesfully generated")

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}
