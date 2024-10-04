package service

import "fmt"

// InvalidAlgorithmError definition
type InvalidAlgorithmError struct {
	Algorithm string
}

func (e *InvalidAlgorithmError) Error() string {
	return fmt.Sprintf("provided algorithm is not supported %s: ", e.Algorithm)
}

func NewInvalidAlgorithmError(algorithm string) *InvalidAlgorithmError {
	return &InvalidAlgorithmError{Algorithm: algorithm}
}

// KeysGenerationError definition
type KeysGenerationError struct {
	Algorithm string
	Err       error
}

func (e *KeysGenerationError) Error() string {
	return fmt.Sprintf("error while generating keys pair with algorithm %s: %s", e.Algorithm, e.Err)
}

func NewKeysGenerationError(algorithm string, err error) *KeysGenerationError {
	return &KeysGenerationError{Algorithm: algorithm, Err: err}
}

// KeysEncodingError definition
type KeysEncodingError struct {
	Algorithm string
	Err       error
}

func (e *KeysEncodingError) Error() string {
	return fmt.Sprintf("error while encoding keys with algorithm %s: %s", e.Algorithm, e.Err)
}

func NewKeysEncodingError(algorithm string, err error) *KeysEncodingError {
	return &KeysEncodingError{Algorithm: algorithm, Err: err}
}

// DeviceNotFoundError definition
type DeviceNotFoundError struct {
	DeviceId string
}

func (e *DeviceNotFoundError) Error() string {
	return fmt.Sprintf("device with id %s not found", e.DeviceId)
}

func NewDeviceNotFoundError(deviceId string) *DeviceNotFoundError {
	return &DeviceNotFoundError{DeviceId: deviceId}
}
