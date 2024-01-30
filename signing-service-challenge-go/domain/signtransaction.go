package domain

import (
	"time"
)

type SignTransaction struct {
	ID          string
	Data        string
	Signature   string
	CreatedTime string
}

func NewSignTransaction(id, data, signature string) (*SignTransaction, error) {
	return &SignTransaction{
		ID:          id,
		Data:        data,
		Signature:   signature,
		CreatedTime: time.Now().String(),
	}, nil
}
