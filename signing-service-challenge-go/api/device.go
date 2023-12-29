package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/service"
)

type CreateSignatureRequest struct {
	Data string `json:"data"`
}

func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		s.createSignatureDevice(response, request)
		return
	} else if request.Method == http.MethodGet {
		s.getSignatureDeviceInfo(response, request)
		return
	} else {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
}

func (s *Server) CreateSignature(response http.ResponseWriter, request *http.Request) {

	deviceId := request.URL.Query().Get("id")
	if deviceId == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			"Please supply the device ID in the \"id\" query parameter",
		})
		return
	}

	var createSignatureRequest CreateSignatureRequest
	err := json.NewDecoder(request.Body).Decode(&createSignatureRequest)
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{
			http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}
	signTransactionResponse, err := service.SignTransaction(deviceId, createSignatureRequest.Data)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			"Error occured while signing:",
			err.Error(),
		})
		return
	}
	WriteAPIResponse(response, http.StatusOK, signTransactionResponse)
}

func (s *Server) createSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var createDeviceRequest domain.CreateSignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&createDeviceRequest)
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{
			http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}
	device := domain.GetSignatureDeviceFromRequest(createDeviceRequest)
	fmt.Printf("Creating device with id [%s]...\n", device.Id)
	createSignatureDeviceResponse, err := service.CreateSignatureDevice(device)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			"Error occured while generating device:",
			err.Error(),
		})
		return
	}
	fmt.Printf("Device with id [%s] created.\n", device.Id)
	WriteAPIResponse(response, http.StatusOK, createSignatureDeviceResponse)
}

func (s *Server) getSignatureDeviceInfo(response http.ResponseWriter, request *http.Request) {

	deviceId := request.URL.Query().Get("id")
	if deviceId == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			"Please supply the device ID in the \"id\" query parameter",
		})
		return
	}

	getDeviceInfoResponse, err := service.GetDeviceInfo(deviceId)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			"Error occured while finding device:",
			err.Error(),
		})
		return
	}

	WriteAPIResponse(response, http.StatusOK, getDeviceInfoResponse)
}

func (s *Server) GetAllDevices(response http.ResponseWriter, request *http.Request) {
	repo := persistence.Get()
	WriteAPIResponse(response, http.StatusOK, repo.GetAllDevices())
}
