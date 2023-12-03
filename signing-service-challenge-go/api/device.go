package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/util"
)

type CreateSignatureDeviceResponse struct {
	Id string `json:"id"`
}

type ListSignatureDevicesResponse struct {
	Devices []SignatureDeviceDTO `json:"signature_devices"`
}

type SignatureDeviceDTO struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

type SignatureResponse struct {
	Signature  string
	SignedData string
}

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	id := request.Header.Get("device_id")
	label := request.Header.Get("label")
	algorithm := request.Header.Get("algorithm")
	device, err := s.signingService.CreateSignatureDevice(id, label, algorithm)
	if err != nil {
		if apiError, ok := err.(domain.ApiError); ok {
			code := apiError.Code
			if code == domain.InvalidDeviceId || code == domain.InvalidSigningAlgorithm {
				WriteErrorResponse(response, http.StatusBadRequest, []string{err.Error()})
			} else if code == domain.DeviceAlreadyExists {
				WriteErrorResponse(response, http.StatusConflict, []string{err.Error()})
			} else {
				WriteErrorResponse(response, http.StatusInternalServerError, []string{http.StatusText(http.StatusInternalServerError)})
			}
		}
	} else {
		WriteAPIResponse(response, http.StatusCreated, CreateSignatureDeviceResponse{Id: device.Id})
	}
}

func (s *Server) DeviceActions(response http.ResponseWriter, request *http.Request) {
	pathParts := strings.Split(request.URL.Path, "/")[1:]
	pathPartsLen := len(pathParts)

	switch pathPartsLen {
	case 4:
		deviceId := pathParts[3]
		if deviceId == "" {
			s.ListSignatureDevices(response, request)
			return
		}
		s.GetSignatureDevice(response, request)
	case 5:
		deviceId := pathParts[3]
		action := pathParts[4]
		if action == "sign" {
			s.SignTransaction(response, request, deviceId)
			return
		}
		WriteInvalidEndpointResponse(response)
	default:
		WriteInvalidEndpointResponse(response)
	}
}

// TODO: Add pagination and filtering
func (s *Server) ListSignatureDevices(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	devices, err := s.signingService.ListSignatureDevices()
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusCreated, ListSignatureDevicesResponse{
			Devices: util.Map(devices, func(device domain.SignatureDevice) SignatureDeviceDTO {
				return SignatureDeviceDTO{
					Id:        device.Id,
					Label:     device.Label,
					Algorithm: device.Algorithm,
				}
			}),
		})
	}
}

// TODO: Add pagination and filtering
func (s *Server) GetSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	id := fetchIdFromPath(request)
	if id == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{"unspecified device id"})
		return
	}

	device, err := s.signingService.GetSignatureDeviceById(id)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusOK, SignatureDeviceDTO{
			Id:        device.Id,
			Label:     device.Label,
			Algorithm: device.Algorithm,
		})
	}
}

func fetchIdFromPath(request *http.Request) string {
	parts := strings.Split(request.URL.Path, "/")
	id := parts[len(parts)-1]
	return id
}

func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request, deviceId string) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	data, err := readData(request, response)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	}

	result, err := s.signingService.SignTransaction(deviceId, data)
	if err != nil {
		// handle error
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, http.StatusOK, SignatureResponse{
		Signature:  util.EncodeToBase64String(result.Signature),
		SignedData: string(result.SignedData),
	})
}

func readData(request *http.Request, response http.ResponseWriter) ([]byte, error) {
	var data []byte

	dataHeader := request.Header.Get("data_to_be_signed")
	if dataHeader != "" {
		data = []byte(dataHeader)
	} else {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			err2 := fmt.Errorf("failed to read request body")
			return nil, errors.Join(err2, err)
		}
		data = body
	}
	return data, nil
}
