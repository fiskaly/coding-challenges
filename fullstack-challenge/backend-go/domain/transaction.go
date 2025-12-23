package domain

import "time"

// Transaction represents a signed piece of data.
type Transaction struct {
	ID                string
	DeviceID          string
	Counter           int
	Timestamp         time.Time
	Data              string
	Signature         string // Base64 encoded
	SignedData        string // The full string that was signed
	PreviousSignature string // Base64 encoded or base64(device.id) for first transaction
}

// TODO: Add validation methods and business logic for Transaction
