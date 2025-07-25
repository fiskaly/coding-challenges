package domain

import (
	crypto "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// structure used to model the device internally
type Device struct {
	Label                            string
	UUID                             string
	Signer                           crypto.SignerI // <- This is where the device implements its own signing method
	LastSignatureBase64EncodedString string
	SignatureCounter                 uint
}
