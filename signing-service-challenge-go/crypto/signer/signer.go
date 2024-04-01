package signer

import "fmt"

func init() {
	registerSigner("RSA", &RSASigner{})
	registerSigner("ECC", &ECDSASigner{})
}

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// generators is used as a registry map
var signers = make(map[string]Signer)

func registerSigner(algorithm string, singer Signer) {
	signers[algorithm] = singer
}

func GetSigner(algorithm string) (Signer, error) {
	signer, exists := signers[algorithm]
	if !exists {
		return nil, fmt.Errorf("no signer exists for algorithm: %s", algorithm)
	}

	return signer, nil
}
