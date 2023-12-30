package domain

import (
	"strings"
)

type SignatureDevice struct {
	Id               string
	PrivateKeyBytes  []byte
	PublicKey        []byte
	Algorithm        string
	SignatureCounter int64
	LastSignature    []byte
	Alias            string
}

type CreateSignatureDeviceResponse struct {
	DeviceId  string `json:"device_id"`
	PublicKey []byte `json:"public_key"`
	Algorithm string `json:"signature_algorithm"`
	Alias     string `json:"alias"`
}

type CreateSignatureDeviceRequest struct {
	DeviceId  string `json:"device_id"`
	Algorithm string `json:"signature_algorithm"`
	Alias     string `json:"alias"`
}

type CreateSignatureResponse struct {
	DeviceId          string `json:"device_id"`
	Signature         []byte `json:"signature"`
	SignedData        string `json:"signed_data"`
	SignaturesCreated int64  `json:"signature_counter"`
	PublicKey         []byte `json:"public_key"`
	Algorithm         string `json:"signature_algorithm"`
	Alias             string `json:"alias"`
}

type SignatureDeviceInfoResponse struct {
	DeviceId          string `json:"device_id"`
	PublicKey         []byte `json:"public_key"`
	Algorithm         string `json:"signature_algorithm"`
	SignaturesCreated int64  `json:"signature_counter"`
	LastSignature     []byte `json:"last_generated_signature"`
	Alias             string `json:"alias"`
}

func NewSignatureDevice(id string, privateKeyBytes []byte, publicKey []byte, algorithm string, alias string) *SignatureDevice {
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

func GetSignatureDeviceFromRequest(request CreateSignatureDeviceRequest) *SignatureDevice {
	return &SignatureDevice{
		Id:        request.DeviceId,
		Algorithm: strings.ToUpper(request.Algorithm),
		Alias:     request.Alias,
	}
}

func (device SignatureDevice) GetCreateSignatureDeviceResponse() *CreateSignatureDeviceResponse {
	return &CreateSignatureDeviceResponse{
		DeviceId:  device.Id,
		PublicKey: device.PublicKey,
		Algorithm: device.Algorithm,
		Alias:     device.Id,
	}
}

func (device SignatureDevice) GetSignatureResponse(signedData string) *CreateSignatureResponse {
	return &CreateSignatureResponse{
		DeviceId:          device.Id,
		Signature:         device.LastSignature,
		SignedData:        signedData,
		PublicKey:         device.PublicKey,
		SignaturesCreated: device.SignatureCounter,
		Algorithm:         device.Algorithm,
		Alias:             device.Id,
	}
}

func (device SignatureDevice) GetSignatureDeviceInfoResponse() *SignatureDeviceInfoResponse {
	return &SignatureDeviceInfoResponse{
		DeviceId:          device.Id,
		LastSignature:     device.LastSignature,
		PublicKey:         device.PublicKey,
		SignaturesCreated: device.SignatureCounter,
		Algorithm:         device.Algorithm,
		Alias:             device.Id,
	}
}
