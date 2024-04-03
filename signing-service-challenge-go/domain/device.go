package domain

import (
	"encoding/base64"
	"time"

	c "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

type SignatureDevice struct {
	ID               string
	Algorithm        string
	Label            string
	PrivateKey       []byte
	PublicKey        []byte
	SignatureCounter int
	LastSignature    string
	CreatedAt        time.Time
}

func NewSignatureDevice(id string, algorithm string, label string) (*SignatureDevice, error) {
	toolkit, err := c.GetToolkit(algorithm)
	if err != nil {
		return nil, err
	}

	keyPair, err := toolkit.Generate()
	if err != nil {
		return nil, err
	}

	privateKey, publicKey, err := toolkit.Marshal(*keyPair)
	if err != nil {
		return nil, err
	}

	return &SignatureDevice{
		ID:               id,
		Algorithm:        algorithm,
		Label:            label,
		PrivateKey:       privateKey,
		PublicKey:        publicKey,
		SignatureCounter: 0,
		LastSignature:    base64.StdEncoding.EncodeToString([]byte(id)),
		CreatedAt:        time.Now(),
	}, nil
}

func (sg *SignatureDevice) SignTransaction(data string) (*Transaction, error) {
	toolKit, err := c.GetToolkit(sg.Algorithm)

	if err != nil {
		return nil, err
	}

	keyPair, err := toolKit.Unmarshal(sg.PrivateKey)
	if err != nil {
		return nil, err
	}

	signature, err := toolKit.Sign(keyPair.Private, []byte(data))
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        uuid.New().String(),
		DeviceID:  sg.ID,
		Data:      data,
		Signature: base64.StdEncoding.EncodeToString(signature),
		Timestamp: time.Now(),
	}, nil
}
