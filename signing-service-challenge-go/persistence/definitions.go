package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type DeviceRepository interface {
	FindDeviceById(id string) (domain.SignatureDevice, error)
	NewDevice(device domain.SignatureDevice) error
	UpdateDevice(device domain.SignatureDevice) error
}
