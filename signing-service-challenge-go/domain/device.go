package domain

import (
    "fmt"
    "github.com/google/uuid"
)

type SignatureDevice struct {
	id uuid.uuid
	publicKey []byte
	privateKey []byte
	algorithm enum.signatureAlgorithm
	signatureCounter int64
	lastSignature string
}
