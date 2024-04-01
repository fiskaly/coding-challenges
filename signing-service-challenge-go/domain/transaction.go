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
