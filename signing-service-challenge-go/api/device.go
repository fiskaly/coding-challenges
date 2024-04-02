package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse
// SignTransaction(deviceId: string, data: string): SignatureResponse
type DeviceHandler struct {
	servie *service.SignatureDeviceService
}

func NewDeviceHandler(srv *service.SignatureDeviceService) *DeviceHandler {
	return &DeviceHandler{servie: srv}
}

type RequesCreateDevice struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type ResponseCreateDevice struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
}

type RequestSignTransaction struct {
	DeviceID string `json:"deviceId"`
	Data     string `json:"data"`
}

// SignatureResponse defines the structure for responses from

func (dh *DeviceHandler) CreateSignatureDevie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestData RequesCreateDevice
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sigDevice, err := dh.servie.CreateSignatureDevice(requestData.ID, requestData.Algorithm, requestData.Label)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	resp := ResponseCreateDevice{
		ID:        sigDevice.ID,
		Algorithm: sigDevice.Algorithm,
	}
	WriteAPIResponse(w, http.StatusCreated, resp)

	defer r.Body.Close()
}

func (dh *DeviceHandler) SignTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestData RequestSignTransaction
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := dh.servie.SignTransaction(requestData.DeviceID, requestData.Data)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{
			"error signing a transaction"})
		return
	}

	WriteAPIResponse(w, http.StatusCreated, signature)

	defer r.Body.Close()
}
