package domain

import (
	"context"
	"errors"
)

// Repository errors to allow storage-agnostic error handling.
var (
	ErrDeviceAlreadyExists = errors.New("repository: signature device already exists")
	ErrDeviceNotFound      = errors.New("repository: signature device not found")
)

// SignatureDeviceRepository defines the storage contract for signature devices.
type SignatureDeviceRepository interface {
	// Save persists a brand-new signature device. It fails with ErrDeviceAlreadyExists when the id is already taken.
	Save(ctx context.Context, device *SignatureDevice) error
	// Update persists changes to an existing device. It fails with ErrDeviceNotFound when the device does not exist.
	Update(ctx context.Context, device *SignatureDevice) error
	// Get retrieves a signature device by id.
	Get(ctx context.Context, id string) (*SignatureDevice, error)
	// List returns all known signature devices.
	List(ctx context.Context) ([]*SignatureDevice, error)
}
