package api

import (
	"encoding/json"
	"fmt"
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

type RequestData struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

func (dh *DeviceHandler) CreateSignatureDevie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Fprintf(w, "Received data: %+v", requestData)

	sigDevice, err := dh.servie.CreateSignatureDevice(requestData.ID, requestData.Algorithm, requestData.Label)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating the device: %v", err), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "device: %+v", sigDevice)

	defer r.Body.Close()
}
