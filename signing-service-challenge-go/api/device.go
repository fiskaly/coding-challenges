package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/application"
)

type DeviceHTTPHandler struct {
	SignatureService application.SignatureDeviceService
	routes           map[string]http.HandlerFunc
}

// SetupRoutes sets up HTTP routes for the application
func (h *DeviceHTTPHandler) setupRoutes() {
	h.routes = make(map[string]http.HandlerFunc)

	h.routes["/api/v0/devices/create"] = h.HandleCreateSignatureDeviceRequest
	h.routes["/api/v0/devices/{id}"] = h.HandleGetSignatureDevice
	h.routes["/api/v0/devices"] = h.HandleListSignatureDevices
	h.routes["/api/v0/devices/"] = h.HandleListSignatureDevices

}
func (h *DeviceHTTPHandler) GetRoutes() map[string]http.HandlerFunc {
	return h.routes
}

type CreateSignatureDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}
type CreateSignatureDeviceResponse struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
	PublicKey string `json:"publicKey"`
}

// NewHTTPHandler creates a new instance of the HTTP handler layer
func NewDeviceHTTPHandler(service application.SignatureDeviceService) *DeviceHTTPHandler {

	var handler = &DeviceHTTPHandler{
		SignatureService: service,
		routes:           make(map[string]http.HandlerFunc),
	}
	handler.setupRoutes()
	return handler
}

// createSignatureDevice handles creating a new signature device
func (h *DeviceHTTPHandler) HandleCreateSignatureDeviceRequest(w http.ResponseWriter, r *http.Request) {

	if !ValidateMethod(w, r, http.MethodPost) {
		return
	}

	var deviceRequest CreateSignatureDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&deviceRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.SignatureService.CreateSignatureDevice(deviceRequest.ID, deviceRequest.Algorithm, deviceRequest.Label)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := CreateSignatureDeviceResponse{
		ID:        device.ID.String(),
		Label:     device.Label,
		Algorithm: device.KeyPairAlgorithm.String(),
		PublicKey: device.PublicKey,
	}
	//json.NewEncoder(w).Encode(device)
	WriteAPIResponse(w, http.StatusOK, response)
}

// HandleGetSignatureDevice handles retrieving a signature device by ID
func (h *DeviceHTTPHandler) HandleGetSignatureDevice(w http.ResponseWriter, r *http.Request) {

	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	deviceID := r.URL.Query().Get("id")
	if deviceID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	//deviceID := r.URL.Path[len("/signature_device/"):]
	device, err := h.SignatureService.GetSignatureDevice(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//json.NewEncoder(w).Encode(device)
	WriteAPIResponse(w, http.StatusOK, device)
}

// HandleListSignatureDevices handles listing all signature devices
func (h *DeviceHTTPHandler) HandleListSignatureDevices(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	devices, err := h.SignatureService.ListSignatureDevices()
	if err != nil {
		//http.Error(w, err.Error(), http.StatusNotFound)
		WriteInternalError(w)
		return
	}
	//json.NewEncoder(w).Encode(devices)
	WriteAPIResponse(w, http.StatusOK, devices)
}
