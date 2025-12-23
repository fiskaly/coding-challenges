package domain

import "time"

// Device represents a signature device with cryptographic capabilities.
type Device struct {
	ID               string
	Label            string
	Algorithm        string // "RSA" or "ECC"
	PublicKey        string // PEM encoded
	PrivateKey       string // PEM encoded - should be stored securely
	SignatureCounter int
	Status           string // "active" or "deactivated"
	CreatedAt        time.Time
}

// TODO: Add validation methods and business logic for Device
