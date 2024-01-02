package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...
func (es *ECCKeyPair) Sign(dataToBeSigned []byte) ([]byte, error) {
	if es.Private == nil {
		return nil, fmt.Errorf("private key is nil")
	}
	return ecdsa.SignASN1(rand.Reader, es.Private, dataToBeSigned)
}

func (rs *RSAKeyPair) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha512.New()
	hash.Write(dataToBeSigned)
	digest := hash.Sum(nil)
	if rs.Private == nil {
		return nil, fmt.Errorf("private key is nil")
	}
	return rsa.SignPKCS1v15(rand.Reader, rs.Private, crypto.SHA512, digest)
}

func SignData(algorithm string, privateKey interface{}, dataToBeSigned []byte) ([]byte, error) {
	switch algorithm {
	case "ECC":
		eccPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("ECC private key type is Invalid")
		}
		eccSigner := &ECCKeyPair{Private: eccPrivateKey}
		return eccSigner.Sign(dataToBeSigned)

	case "RSA":
		rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("invalid RSA private key type")
		}
		rsaSigner := &RSAKeyPair{Private: rsaPrivateKey}
		return rsaSigner.Sign(dataToBeSigned)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}
