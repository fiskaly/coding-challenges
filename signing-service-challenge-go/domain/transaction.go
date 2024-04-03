package domain

import "time"

// Transaction
type Transaction struct {
	ID        string
	DeviceID  string
	Data      string
	Signature string
	Timestamp time.Time
}

type Signature struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type TransactionResp struct {
	ID        string    `json:"transaction_id"`
	Signature string    `json:"signature"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionsByDeviceResp struct {
	ID           string             `json:"device_id"`
	Algorithm    string             `json:"algorithm"`
	Label        string             `json:"label"`
	CreatedAt    time.Time          `json:"created_at"`
	Transactions []*TransactionResp `json:"transactions"`
}
