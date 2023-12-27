package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
)

type SignatureDevice struct {
	Id               string
	PrivateKeyBytes  []byte
	PublicKey        string
	Algorithm        crypto.SignatureAlgorithm
	SignatureCounter int64
	LastSignature    string
	Alias            string
}

func NewSignatureDevice(id string, privateKeyBytes []byte, publicKey string, algorithm crypto.SignatureAlgorithm, alias string) *SignatureDevice {
	return &SignatureDevice{
		Id:               id,
		PrivateKeyBytes:  privateKeyBytes,
		PublicKey:        publicKey,
		Algorithm:        algorithm,
		SignatureCounter: 0,
		LastSignature:    "",
		Alias:            alias,
	}
}

type CreateSignatureDeviceResponse struct {
	Id        string
	PublicKey string
	Algorithm crypto.SignatureAlgorithm
	Alias     string
}

func (device SignatureDevice) GetCreSignatureDeviceResponse() *CreateSignatureDeviceResponse {
	return &CreateSignatureDeviceResponse{
		Id:        device.Id,
		PublicKey: device.PublicKey,
		Algorithm: device.Algorithm,
		Alias:     device.Id,
	}
}
