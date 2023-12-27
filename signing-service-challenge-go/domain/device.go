package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/crypto"
)

type SignatureDevice struct {
	Id               string
	PrivateKeyBytes  []byte
	PublicKey        string
	Algorithm        crypto.SignatureAlgorithm
	SignatureCounter int64
	LastSignature    []byte
	Alias            string
}

func NewSignatureDevice(id string, privateKeyBytes []byte, publicKey string, algorithm crypto.SignatureAlgorithm, alias string) *SignatureDevice {
	return &SignatureDevice{
		Id:               id,
		PrivateKeyBytes:  privateKeyBytes,
		PublicKey:        publicKey,
		Algorithm:        algorithm,
		SignatureCounter: 0,
		LastSignature:    nil,
		Alias:            alias,
	}
}

func (device SignatureDevice) GetCreSignatureDeviceResponse() *api.CreateSignatureDeviceResponse {
	return &api.CreateSignatureDeviceResponse{
		DeviceId:  device.Id,
		PublicKey: device.PublicKey,
		Algorithm: device.Algorithm,
		Alias:     device.Id,
	}
}

func (device SignatureDevice) GetSignatureResponse() *api.SignatureResponse {
	return &api.SignatureResponse{
		DeviceId:          device.Id,
		Signature:         string(device.LastSignature),
		PublicKey:         device.PublicKey,
		SignaturesCreated: device.SignatureCounter,
		Algorithm:         device.Algorithm,
		Alias:             device.Id,
	}
}

func (device SignatureDevice) GetSignatureDeviceInfoResponse() *api.SignatureDeviceInfoResponse {
	return &api.SignatureDeviceInfoResponse{
		DeviceId:          device.Id,
		LastSignature:     string(device.LastSignature),
		PublicKey:         device.PublicKey,
		SignaturesCreated: device.SignatureCounter,
		Algorithm:         device.Algorithm,
		Alias:             device.Id,
	}
}
