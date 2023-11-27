package api

import (
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type CreateSignatureDeviceResponse struct {
	Id string `json:"id"`
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
	algorithm := request.Header.Get("alorithm")

	device, err := domain.NewSignatureDevice(id, &label, algorithm)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{err.Error()})
	}

	WriteAPIResponse(response, http.StatusCreated, CreateSignatureDeviceResponse{Id: device.Id})
}
