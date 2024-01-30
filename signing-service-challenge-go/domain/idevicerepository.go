package domain

// SignatureDeviceRepository defines the interface for signature device storage
type SignatureDeviceRepository interface {
	AddDevice(device *SignatureDevice) error
	UpdateDevice(device *SignatureDevice) error
	GetDeviceByID(id string) (*SignatureDevice, error)
	ListDevices() ([]*SignatureDevice, error)
	SaveSignTransaction(device *SignatureDevice, sign_transaction *SignTransaction) error
	GetSignTransactionsForDevice(id string) ([]*SignTransaction, error)
}
