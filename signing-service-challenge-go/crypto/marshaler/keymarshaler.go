package marshaler

import c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"

// KeyMarshaler the interface for encoding and decoding key pairs.
type KeyMarshaler interface {
	Marshal(keyPair c.KeyPair) ([]byte, []byte, error)
	Unmarshal(privateKeyBytes []byte) (*c.KeyPair, error)
}
