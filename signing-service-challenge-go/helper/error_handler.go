package helper

import (
	"errors"
	"net/http"
	"signing-service-challenge/crypto"
	"signing-service-challenge/service"
)

func HandleDeviceServiceError(err error) (int, string) {
	var keysGenerationError *service.KeysGenerationError
	var keysEncodingError *service.KeysEncodingError
	var invalidAlgorithmError *service.InvalidAlgorithmError
	var deviceNotFoundError *service.DeviceNotFoundError
	var signOperationError *crypto.SignOperationError
	var marshalError *crypto.MarshalError

	switch {
	case errors.As(err, &keysGenerationError):
		return http.StatusUnprocessableEntity, keysGenerationError.Error()
	case errors.As(err, &keysEncodingError):
		return http.StatusUnprocessableEntity, keysEncodingError.Error()
	case errors.As(err, &invalidAlgorithmError):
		return http.StatusUnprocessableEntity, invalidAlgorithmError.Error()
	case errors.As(err, &deviceNotFoundError):
		return http.StatusNotFound, deviceNotFoundError.Error()
	case errors.As(err, &deviceNotFoundError):
		return http.StatusBadRequest, deviceNotFoundError.Error()
	case errors.As(err, &signOperationError):
		return http.StatusUnprocessableEntity, signOperationError.Error()
	case errors.As(err, &marshalError):
		return http.StatusUnprocessableEntity, marshalError.Error()
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
