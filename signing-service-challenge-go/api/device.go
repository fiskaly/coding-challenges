package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/service"
	"github.com/gorilla/mux"
)

type CreateSignatureRequest struct {
	Data string `json:"data"`
}

// Device is the handler for all .../devices calls
// POST .../devices gets routed to Server.createSignatureDevice
// GET .../devices gets routed to Server.getAllDevices
func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		s.createSignatureDevice(response, request)
		return
	} else if request.Method == http.MethodGet {
		s.getAllDevices(response, request)
		return
	} else {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
}

// Creates a signature device based on the supplied id, algorithm and (optinal) alias
// Writes the created devices information
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

// Fetches and writes all created device ID's
func (s *Server) getAllDevices(response http.ResponseWriter, request *http.Request) {
	repo := persistence.Get()
	WriteAPIResponse(response, http.StatusOK, repo.GetAllDevices())
}

// Creates a signature for the supplied data using the selected device.
// Writes signature, signed data string and device information
func (s *Server) CreateSignature(response http.ResponseWriter, request *http.Request) {

	deviceId := mux.Vars(request)["id"]

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

// Fetches the specified device and writes all the public info for the device
func (s *Server) GetSignatureDeviceInfo(response http.ResponseWriter, request *http.Request) {

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
