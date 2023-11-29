package api

import (
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

func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	id := request.Header.Get("id")
	label := request.Header.Get("label")
	algorithm := request.Header.Get("algorithm")
	device, err := s.signingService.CreateSignatureDevice(id, label, algorithm)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusCreated, CreateSignatureDeviceResponse{Id: device.Id})
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
