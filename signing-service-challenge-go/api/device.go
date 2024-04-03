package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

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
	ID               string    `json:"id"`
	Algorithm        string    `json:"algorithm"`
	Label            string    `json:"label"`
	SignatureCounter int       `json:"signature_counter"`
	CreatedAt        time.Time `json:"created_at"`
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
	defer r.Body.Close()
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
		ID:               sigDevice.ID,
		Algorithm:        sigDevice.Algorithm,
		Label:            sigDevice.Label,
		SignatureCounter: sigDevice.SignatureCounter,
		CreatedAt:        sigDevice.CreatedAt,
	}
	WriteAPIResponse(w, http.StatusCreated, resp)
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
	defer r.Body.Close()
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
}

func (dh *DeviceHandler) GetSignatureDeviceById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	deviceID := r.URL.Query().Get("id")
	if deviceID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"device id is not valid"})
		return
	}

	sigDevice, err := dh.servie.GetSignatureDeviceByID(deviceID)
	if err != nil {
		if deviceID == "" {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{fmt.Sprintf("error fetching devince with ID: %s", deviceID)})
			return
		}
	}

	resp := ResponseCreateDevice{
		ID:               sigDevice.ID,
		Algorithm:        sigDevice.Algorithm,
		Label:            sigDevice.Label,
		SignatureCounter: sigDevice.SignatureCounter,
	}
	WriteAPIResponse(w, http.StatusCreated, resp)
}

func (dh *DeviceHandler) GetTransactionsByServiceId(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"device id is not valid"})
		return
	}

	transactions, err := dh.servie.GetTransactionsByDeviceID(deviceID)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusCreated, transactions)
}
