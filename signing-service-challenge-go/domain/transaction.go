package domain

import "time"

// Transaction
type Transaction struct {
	ID        string
	DeviceID  string
	Data      string
	Signature []byte
	Timestamp time.Time
}

type Signature struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
