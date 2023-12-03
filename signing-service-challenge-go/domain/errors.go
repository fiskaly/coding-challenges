package domain

import "fmt"

type ApiError struct {
	Code    int
	message string
}

func (e ApiError) Error() string {
	return e.message
}

const (
	InvalidDeviceId = iota
	InvalidSigningAlgorithm
	DeviceNotFound
	DeviceAlreadyExists
)

func ErrorInvalidDeviceId() error {
	return ApiError{
		Code:    InvalidDeviceId,
		message: "invalid device id",
	}
}

func ErrorInvalidSigningAlgorithm() error {
	return ApiError{
		Code:    InvalidDeviceId,
		message: "invalid signing algorithm",
	}
}

func ErrorDeviceNotFound(deviceId string) error {
	return ApiError{
		Code:    DeviceNotFound,
		message: fmt.Sprintf("device with id %s does not exist", deviceId),
	}
}

func ErrorDeviceAlreadyExists(deviceId string) error {
	return ApiError{
		Code:    DeviceAlreadyExists,
		message: fmt.Sprintf("device with id %s already exists", deviceId),
	}
}
