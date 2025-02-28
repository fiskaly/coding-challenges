package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

var uuidRegex = regexp.MustCompile("^[a-fA-F0-9-]{36}$")

type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

type CreateSignatureDeviceResponse struct {
	ID               string `json:"id"`
	Algorithm        string `json:"algorithm"`
	Label            string `json:"label,omitempty"`
	SignatureCounter int    `json:"signature_counter"`
}

type SignatureRequest struct {
	Data string `json:"data"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type DeviceHandler struct {
	store persistence.Storage
}

func NewDeviceHandler(store persistence.Storage) *DeviceHandler {
	return &DeviceHandler{store: store}
}

func (h *DeviceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v0/devices", h.HandleDevices)
	mux.HandleFunc("/api/v0/devices/", h.HandleDeviceActions)
}

func (h *DeviceHandler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateSignatureDevice(w, r)
	case http.MethodGet:
		h.ListDevices(w, r)
	default:
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
	}
}

func (h *DeviceHandler) HandleDeviceActions(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v0/devices/")
	id = strings.Split(id, "/")[0]

	if !uuidRegex.MatchString(id) && id != "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid device ID format"})
		return
	}

	switch {
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/sign"):
		h.SignTransaction(w, r, id)
	case r.Method == http.MethodGet && id != "":
		h.GetDeviceDetails(w, r, id)
	default:
		WriteErrorResponse(w, http.StatusNotFound, []string{"Endpoint not found"})
	}
}

func (h *DeviceHandler) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	var req CreateSignatureDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}
	if req.Algorithm != "RSA" && req.Algorithm != "ECC" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Unsupported algorithm"})
		return
	}

	req.Label = sanitizeLabel(req.Label)
	device, err := domain.CreateNewDevice(req.Algorithm, req.Label)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to create device"})
		return
	}
	if err := h.store.Save(device); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to save device"})
		return
	}
	response := CreateSignatureDeviceResponse{
		ID:               device.GetID(),
		Algorithm:        device.GetAlgorithm(),
		Label:            req.Label,
		SignatureCounter: 0,
	}
	WriteAPIResponse(w, http.StatusCreated, response)
}

func (h *DeviceHandler) SignTransaction(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	device, err := h.store.FindByID(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
		return
	}
	var req SignatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}
	if len(req.Data) == 0 || len(req.Data) > 5000 {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid data length"})
		return
	}

	signature, securedData, err := device.SignData(req.Data)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to sign data"})
		return
	}
	if err := h.store.Update(device); err != nil {
		fmt.Printf("Failed to update device: %v\n", err)
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to update device"})
		return
	}
	response := SignatureResponse{
		Signature:  signature,
		SignedData: securedData,
	}
	WriteAPIResponse(w, http.StatusOK, response)
}

func (h *DeviceHandler) GetDeviceDetails(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	device, err := h.store.FindByID(id)
	if err != nil {
		WriteErrorResponse(w, http.StatusNotFound, []string{"Device not found"})
		return
	}
	response := CreateSignatureDeviceResponse{
		ID:               device.GetID(),
		Algorithm:        device.GetAlgorithm(),
		Label:            device.GetLabel(),
		SignatureCounter: device.GetSignatureCounter(),
	}
	WriteAPIResponse(w, http.StatusOK, response)
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	devices, err := h.store.FindAll()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to retrieve devices"})
		return
	}

	responses := make([]CreateSignatureDeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = CreateSignatureDeviceResponse{
			ID:               device.GetID(),
			Algorithm:        device.GetAlgorithm(),
			Label:            device.GetLabel(),
			SignatureCounter: device.GetSignatureCounter(),
		}
	}
	WriteAPIResponse(w, http.StatusOK, responses)
}

func sanitizeLabel(label string) string {
	label = strings.TrimSpace(label)
	// Allow only letters, digits, spaces, underscores, and dashes.
	var sanitized strings.Builder
	for _, r := range label {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == ' ' || r == '_' || r == '-' {
			sanitized.WriteRune(r)
		}
	}
	return sanitized.String()
}
