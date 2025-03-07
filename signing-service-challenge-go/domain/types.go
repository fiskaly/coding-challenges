package domain

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

var ErrAlreadyExists = persistence.ErrAlreadyExists
var ErrUnsupportedAlgorithm = crypto.ErrUnsupportedAlgorithm

type SignatureDomain interface {
	CreateSignatureDevice(Device) error
	ListSignatureDevices() ([]ID, error)
	SignTransaction(ID, Data) (*CreatedSignature, error)
	GetSignatureDeviceDetails(ID) (Device, error)
}

type ID string

type Data string

type Device struct {
	ID        ID
	Algorithm string
	Label     string
}

type CreatedSignature struct {
	SignedData string
	Signature  string
}
