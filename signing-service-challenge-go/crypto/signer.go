package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"
)

// Signer defines a contract for different types of signing implementations.
// For future implementation of other algorithms, one has only to make sure
// that the new signing algorithm adheres to the SignerI interface
type SignerI interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// structure used to model the ECC Signer
type ECCSigner struct {
	Marshaler    *ECCMarshaler
	KeyPair      *ECCKeyPair
	KeyGenerator *ECCGenerator
}

// ctor of a new ECC Signer
func NewECCSigner() (*ECCSigner, error) {

	// create base stuct
	signerEcc := &ECCSigner{
		Marshaler:    NewECCMarshaler(),
		KeyGenerator: &ECCGenerator{},
	}

	// generate keys using the KeyGenerator and use them to fill the ECC Signer
	var err error
	signerEcc.KeyPair, err = signerEcc.KeyGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("could not generate key pair associated to the ECC Signer: %s", err)
	}

	log.Debug("New ECCSigner created")

	return signerEcc, nil
}

// implementation of the Signing function for the ECC Signer
func (s ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {

	// check input
	if s.KeyPair == nil || s.KeyPair.Private == nil || s.KeyPair.Public == nil {
		return nil, fmt.Errorf("ECC signer is missing some key")
	}
	if len(dataToBeSigned) == 0 {
		return nil, fmt.Errorf("zero bytes passed to the ECC signer")
	}

	// hash data
	hashedDataToBeSigned := sha256.Sum256(dataToBeSigned)

	// sign
	rValue, sValue, err := ecdsa.Sign(rand.Reader, s.KeyPair.Private, hashedDataToBeSigned[:])
	if err != nil {
		return nil, fmt.Errorf("error signing with ECC : %s", err)
	}

	// pack in a struct and marshal
	signedData, err := asn1.Marshal(struct {
		RValue *big.Int
		SValue *big.Int
	}{
		RValue: rValue,
		SValue: sValue,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling with ECC : %s", err)
	}

	log.Debug("ECCSigner succesfully signed the message")

	return signedData, nil
}

// structure used to model the RSA Signer
type RSASigner struct {
	Marshaler    *RSAMarshaler
	KeyPair      *RSAKeyPair
	KeyGenerator *RSAGenerator
}

// ctor of a new RSA Signer
func NewRSASigner() (*RSASigner, error) {

	// create base struct
	signerRsa := &RSASigner{
		Marshaler:    NewRSAMarshaler(),
		KeyGenerator: &RSAGenerator{},
	}

	// generate keys using the RSA generator and fill the RSA signer with them
	var err error
	signerRsa.KeyPair, err = signerRsa.KeyGenerator.Generate()
	if err != nil {
		return nil, fmt.Errorf("could not generate key pair associated to the RSA Signer: %s", err)
	}

	log.Debug("New RSASigner created")

	return signerRsa, nil
}

// implementation of the Signing function for the RSA Signer
func (s RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {

	// check input
	if s.KeyPair == nil || s.KeyPair.Private == nil || s.KeyPair.Public == nil {
		return nil, fmt.Errorf("RSA signer is missing some key")
	}
	if len(dataToBeSigned) == 0 {
		return nil, fmt.Errorf("zero bytes passed to the RSA signer")
	}

	// hash data
	hashedDataToBeSigned := sha256.Sum256(dataToBeSigned)

	// sign
	signedData, err := rsa.SignPSS(rand.Reader, s.KeyPair.Private, crypto.SHA256, hashedDataToBeSigned[:], nil)
	if err != nil {
		return nil, fmt.Errorf("error signing with RSA: %s", err)
	}

	log.Debug("RSASigner succesfully signed the message")

	return signedData, nil

}
