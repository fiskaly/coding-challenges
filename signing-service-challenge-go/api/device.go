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
	//h.routes["/api/v0/devices/{id}"] = h.HandleGetSignatureDevice
	h.routes["/api/v0/devices"] = h.HandleListSignatureDevices
	h.routes["/api/v0/devices/get"] = h.HandleGetSignatureDevice
	h.routes["/api/v0/signatures/create"] = h.HandleSignTransaction
	h.routes["/api/v0/signatures/get"] = h.HandleListSignTransactions
	//h.routes["/api/v0/signatures"] = h.HandleListSignTransactions

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
type GetByIdRequest struct {
	ID string `json:"id"`
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
	WriteAPIResponse(w, http.StatusOK, response)
}

// HandleGetSignatureDevice handles retrieving a signature device by ID
func (h *DeviceHTTPHandler) HandleGetSignatureDevice(w http.ResponseWriter, r *http.Request) {

	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	var deviceIdRequest GetByIdRequest
	err := json.NewDecoder(r.Body).Decode(&deviceIdRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deviceID := deviceIdRequest.ID
	if deviceID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	device, err := h.SignatureService.GetSignatureDevice(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	WriteAPIResponse(w, http.StatusOK, device)
}

// HandleListSignatureDevices handles listing all signature devices
func (h *DeviceHTTPHandler) HandleListSignatureDevices(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	devices, err := h.SignatureService.ListSignatureDevices()
	if err != nil {
		WriteInternalError(w)
		return
	}
	WriteAPIResponse(w, http.StatusOK, devices)
}

type SignTransactionRequest struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type SignTransactionResponse struct {
	ID          string `json:"id"`
	Data        string `json:"data"`
	Signature   string `json:"signature"`
	CreatedTime string `json:"createdtime"`
}

func (h *DeviceHTTPHandler) HandleSignTransaction(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(w, r, http.MethodPost) {
		return
	}
	var signRequest SignTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&signRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sign_transaction, err := h.SignatureService.SignData(signRequest.ID, signRequest.Data)
	if err != nil {
		WriteInternalError(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response := SignTransactionResponse{
		ID:          sign_transaction.ID,
		Signature:   sign_transaction.Signature,
		Data:        sign_transaction.Data,
		CreatedTime: sign_transaction.CreatedTime,
	}
	WriteAPIResponse(w, http.StatusOK, response)
}

func (h *DeviceHTTPHandler) HandleListSignTransactions(w http.ResponseWriter, r *http.Request) {
	if !ValidateMethod(w, r, http.MethodGet) {
		return
	}
	var deviceIdRequest GetByIdRequest
	err := json.NewDecoder(r.Body).Decode(&deviceIdRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deviceID := deviceIdRequest.ID
	if deviceID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	devices, err := h.SignatureService.ListSignTransactions(deviceID)
	if err != nil {
		WriteInternalError(w)
		return
	}
	WriteAPIResponse(w, http.StatusOK, devices)
}
