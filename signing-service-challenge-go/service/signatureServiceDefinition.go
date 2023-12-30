package service

import "github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"

type SignatureService interface {
	CreateSignatureDevice(device *domain.SignatureDevice) (domain.CreateSignatureDeviceResponse, error)
	SignTransaction(deviceId string, data string) (domain.CreateSignatureResponse, error)
	GetDeviceInfo(deviceId string) (domain.SignatureDeviceInfoResponse, error)
	GetAllDevices() []string
}
