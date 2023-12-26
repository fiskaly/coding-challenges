package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

type SignatureDevice struct {
	Id               string
	PublicKey        []byte
	Algorithm        crypto.SignatureAlgorithm
	SignatureCounter int64
	LastSignature    string
	Alias            string
}
