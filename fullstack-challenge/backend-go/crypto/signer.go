package crypto

// Signer is the interface that must be implemented by all signature algorithms.
type Signer interface {
	Sign(data []byte) ([]byte, error)
}

// PEMEncoder is the interface for encoding keys to PEM format.
type PEMEncoder interface {
	EncodePEM() (publicPEM, privatePEM string, err error)
}
