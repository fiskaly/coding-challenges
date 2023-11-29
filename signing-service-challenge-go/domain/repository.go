package domain

type SignatureDeviceRepository interface {
	StoreSignatureDevice(device *SignatureDevice) error
	ListSignatureDevices() ([]SignatureDevice, error)
	GetSignatureDeviceById(id string) *SignatureDevice
}
