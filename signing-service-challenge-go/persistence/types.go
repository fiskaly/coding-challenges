package persistence

import "fmt"

type ID string

type Algorithm string

type Signature string

type Device struct {
	ID             ID
	Algorithm      Algorithm
	Label          *string
	LastSignature  Signature
	SignatureCount int
	PrivateKey     []byte
}

var ErrAlreadyExists = fmt.Errorf("device already exists")
var ErrRace = fmt.Errorf("compare failed in compare and swap")

type DevicePersister interface {
	AddDevice(Device) error
	ListDevices() ([]ID, error)
	GetDevice(ID) (*Device, error)
	CompareAndSwapSignature(ID, Signature, Signature) error
}
