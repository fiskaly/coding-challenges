package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// DeviceHandler bundles HTTP handlers for signature device operations.
type DeviceHandler struct {
	service *service.SignatureService
}

// NewDeviceHandler constructs a handler backed by the signature service.
func NewDeviceHandler(service *service.SignatureService) *DeviceHandler {
	return &DeviceHandler{
		service: service,
	}
}

// HandleCollection processes requests on /api/v0/devices.
func (h *DeviceHandler) HandleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createDevice(w, r)
	case http.MethodGet:
		h.listDevices(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleResource processes requests on /api/v0/devices/{id} and /api/v0/devices/{id}/sign.
func (h *DeviceHandler) HandleResource(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v0/devices/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	parts := strings.Split(path, "/")
	deviceID := parts[0]

	switch {
	case len(parts) == 1 && r.Method == http.MethodGet:
		h.getDevice(w, r, deviceID)
	case len(parts) == 2 && parts[1] == "sign" && r.Method == http.MethodPost:
		h.signDevice(w, r, deviceID)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

type createDeviceRequest struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

type createDeviceResponse struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter uint64 `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
	PublicKeyPEM     string `json:"public_key_pem"`
}

type deviceResponse struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter uint64 `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
}

type signDeviceRequest struct {
	Data string `json:"data"`
}

type signDeviceResponse struct {
	DeviceID         string `json:"device_id"`
	Signature        string `json:"signature"`
	SignedData       string `json:"signed_data"`
	SignatureCounter uint64 `json:"signature_counter"`
	LastSignature    string `json:"last_signature"`
}

func (h *DeviceHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var payload createDeviceRequest
	if err := decodeJSON(r.Body, &payload); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		return
	}

	result, err := h.service.CreateSignatureDevice(r.Context(), service.CreateSignatureDeviceRequest{
		ID:        payload.ID,
		Algorithm: payload.Algorithm,
		Label:     payload.Label,
	})
	if err != nil {
		status, messages := mapError(err)
		WriteErrorResponse(w, status, messages)
		return
	}

	response := createDeviceResponse{
		ID:               result.ID,
		Label:            result.Label,
		Algorithm:        result.Algorithm.String(),
		SignatureCounter: result.SignatureCounter,
		LastSignature:    result.LastSignature,
		PublicKeyPEM:     string(result.PublicKeyPEM),
	}

	WriteAPIResponse(w, http.StatusCreated, response)
}

func (h *DeviceHandler) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.service.ListSignatureDevices(r.Context())
	if err != nil {
		status, messages := mapError(err)
		WriteErrorResponse(w, status, messages)
		return
	}

	response := make([]deviceResponse, 0, len(devices))
	for _, device := range devices {
		response = append(response, deviceResponse{
			ID:               device.ID,
			Label:            device.Label,
			Algorithm:        device.Algorithm.String(),
			SignatureCounter: device.SignatureCounter,
			LastSignature:    device.LastSignature,
		})
	}

	WriteAPIResponse(w, http.StatusOK, response)
}

func (h *DeviceHandler) getDevice(w http.ResponseWriter, r *http.Request, deviceID string) {
	device, err := h.service.GetSignatureDevice(r.Context(), deviceID)
	if err != nil {
		status, messages := mapError(err)
		WriteErrorResponse(w, status, messages)
		return
	}

	response := deviceResponse{
		ID:               device.ID,
		Label:            device.Label,
		Algorithm:        device.Algorithm.String(),
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
	}

	WriteAPIResponse(w, http.StatusOK, response)
}

func (h *DeviceHandler) signDevice(w http.ResponseWriter, r *http.Request, deviceID string) {
	var payload signDeviceRequest
	if err := decodeJSON(r.Body, &payload); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		return
	}

	result, err := h.service.SignTransaction(r.Context(), service.SignTransactionRequest{
		DeviceID: deviceID,
		Data:     payload.Data,
	})
	if err != nil {
		status, messages := mapError(err)
		WriteErrorResponse(w, status, messages)
		return
	}

	response := signDeviceResponse{
		DeviceID:         result.DeviceID,
		Signature:        result.Signature,
		SignedData:       result.SignedData,
		SignatureCounter: result.SignatureCounter,
		LastSignature:    result.LastSignature,
	}

	WriteAPIResponse(w, http.StatusOK, response)
}

func decodeJSON(r io.ReadCloser, target interface{}) error {
	defer r.Close()

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func mapError(err error) (int, []string) {
	switch {
	case err == nil:
		return http.StatusOK, nil
	case errors.Is(err, service.ErrInvalidDeviceID),
		errors.Is(err, service.ErrInvalidData),
		errors.Is(err, domain.ErrUnsupportedAlgorithm),
		errors.Is(err, domain.ErrInvalidDeviceID):
		return http.StatusBadRequest, []string{err.Error()}
	case errors.Is(err, domain.ErrDeviceAlreadyExists):
		return http.StatusConflict, []string{err.Error()}
	case errors.Is(err, domain.ErrDeviceNotFound):
		return http.StatusNotFound, []string{err.Error()}
	case errors.Is(err, service.ErrMissingSigner):
		return http.StatusServiceUnavailable, []string{err.Error()}
	default:
		return http.StatusInternalServerError, []string{http.StatusText(http.StatusInternalServerError)}
	}
}
