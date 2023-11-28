package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/util"
)

type CreateSignatureDeviceResponse struct {
	Id string `json:"id"`
}

type ListSignatureDevicesResponse struct {
	Devices []SignatureDevice `json:"signature_devices"`
}

type SignatureDevice struct {
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
	device, err := createSignatureDevice(id, label, algorithm)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusCreated, CreateSignatureDeviceResponse{Id: device.Id})
	}
}

func createSignatureDevice(id string, label string, algorithm string) (*domain.SignatureDevice, error) {
	return &domain.SignatureDevice{
		Id:    id,
		Label: &label,
	}, nil
	// return domain.CreateSignatureDevice(id, &label, algorithm)
}

// TODO: Add pagination and filtering
func (s *Server) ListSignatureDevices(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	devices, err := listSignatureDevices()
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusCreated, ListSignatureDevicesResponse{
			Devices: util.Map(devices, func(device *domain.SignatureDevice) SignatureDevice {
				return SignatureDevice{
					Id:        device.Id,
					Label:     *device.Label,
					Algorithm: device.Algorithm,
				}
			}),
		})
	}
}

func listSignatureDevices() ([]*domain.SignatureDevice, error) {
	label1 := "Device 1"
	label2 := "Device 2"
	return []*domain.SignatureDevice{
		{
			Id:        "1",
			Label:     &label1,
			Algorithm: "RSA",
		},
		{
			Id:        "2",
			Label:     &label2,
			Algorithm: "ECC",
		},
	}, nil
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
	device, err := getSignatureDevice(id)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	} else {
		WriteAPIResponse(response, http.StatusOK, SignatureDevice{
			Id:        device.Id,
			Label:     *device.Label,
			Algorithm: device.Algorithm,
		})
	}
}

func fetchIdFromPath(request *http.Request) string {
	parts := strings.Split(request.URL.Path, "/")
	id := parts[len(parts)-1]
	return id
}

func getSignatureDevice(id string) (*domain.SignatureDevice, error) {
	if id == "" {
		return nil, errors.New("unspecified device id")
	}

	label := fmt.Sprintf("Device %s", id)
	return &domain.SignatureDevice{
		Id:        id,
		Label:     &label,
		Algorithm: "RSA",
	}, nil
}
