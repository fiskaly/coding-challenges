package crypto

import "crypto"

// KeyPair holds a generic private/public key.
type KeyPair struct {
	Public  crypto.PublicKey
	Private crypto.PrivateKey
}
