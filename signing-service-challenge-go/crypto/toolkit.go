package crypto

import (
	"crypto"
)

// KeyPairGenerator - generating key pairs
type KeyPairGenerator interface {
	Generate() (*KeyPair, error)
}

// Signer - signing transactions
type Signer interface {
	Sign(pk crypto.PrivateKey, dataToBeSigned []byte) ([]byte, error)
}

// KeyMarshaler - encoding and decoding key pairs
type KeyMarshaler interface {
	Marshal(keyPair KeyPair) ([]byte, []byte, error)
	Unmarshal(privateKeyBytes []byte) (*KeyPair, error)
}

type CryptoToolkit interface {
	KeyPairGenerator
	Signer
	KeyMarshaler
}
