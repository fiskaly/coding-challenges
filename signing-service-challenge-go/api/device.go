package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"

	"github.com/google/uuid"
)

type CreateSignatureDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type CreateSignatureDeviceResponse struct {
	Message string `json:"message"`
}

type SignTransactionRequest struct {
	DeviceID       string `json:"deviceId"`
	DataToBeSigned string `json:"data"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type SignatureService struct {
	mu            sync.RWMutex
	devices       map[uuid.UUID]*domain.SignatureDevice
	signatureAlgo map[string]string
}

type CheckDeviceExistsRequest struct {
	DeviceID string `json:"deviceId"`
}

type CheckDeviceExistsResponse struct {
	Message string `json:"message"`
}

func NewSignatureService() *SignatureService {
	return &SignatureService{
		devices:       make(map[uuid.UUID]*domain.SignatureDevice),
		signatureAlgo: map[string]string{"ECC": "ECDSA", "RSA": "RSA"},
	}
}

func (s *SignatureService) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method Not Allowed"})
		return
	}

	var request CreateSignatureDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON request"})
		return
	}

	deviceID, err := uuid.Parse(request.ID)
	if HandleError(w, http.StatusBadRequest, err, "Invalid UUID") {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.devices[deviceID]; exists {
		http.Error(w, "Device with the given ID already exists", http.StatusConflict)
		return
	}

	algo := strings.ToUpper(request.Algorithm)
	if algo != "ECC" && algo != "RSA" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid algorithm"})
		return
	}

	device := &domain.SignatureDevice{
		ID:        deviceID,
		Algorithm: algo,
		Label:     request.Label,
	}

	err = device.AddDevice()
	if HandleError(w, http.StatusBadRequest, err, "Device does not exist") {
		return
	}

	s.devices[deviceID] = device
	time.Sleep(2 * time.Second)
	response := CreateSignatureDeviceResponse{Message: fmt.Sprintf("The device %s has been successfully created", deviceID)}

	WriteAPIResponse(w, http.StatusCreated, response)
}

func (s *SignatureService) SignTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method Not Allowed"})
		return
	}
	var request SignTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON request"})
		return
	}
	var err error

	deviceID, err := uuid.Parse(request.DeviceID)
	if HandleError(w, http.StatusBadRequest, err, "Invalid device_id") {
		return
	}

	device, exists := s.GetSignatureDevice(deviceID)
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if device.SignatureCounter == 0 {
		device.LastSignature = base64.StdEncoding.EncodeToString([]byte(device.ID.String()))
	}
	unsignedData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, request.DataToBeSigned, device.LastSignature)
	signature, err := crypto.SignData(device.Algorithm, device.PrivateKey, []byte(unsignedData))
	if HandleError(w, http.StatusInternalServerError, err, "Failed to sign the data") {
		return
	}

	signatureEncoded := base64.StdEncoding.EncodeToString([]byte(signature))
	device.SignatureCounter++
	device.LastSignature = string(signatureEncoded)
	err = persistence.SaveDevice(*device, fmt.Sprintf("%s.json", request.DeviceID))
	if HandleError(w, http.StatusInternalServerError, err, "Failed to Save the Device in Json") {
		return
	}
	encodeJSONResponse(w, http.StatusOK, SignatureResponse{
		Signature:  signatureEncoded,
		SignedData: unsignedData,
	})

}

func (s *SignatureService) CheckDeviceExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method Not Allowed"})
		return
	}
	var request SignTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON request"})
		return
	}

	deviceID, err := uuid.Parse(request.DeviceID)

	if HandleError(w, http.StatusBadRequest, err, "Invalid device_id") {
		return
	}

	device, exists := s.GetSignatureDevice(deviceID)
	if !exists {
		http.Error(w, "Device not found. Please use CreateSignatureDevice API to add the deivce", http.StatusNotFound)
		return
	}
	response := CreateSignatureDeviceResponse{Message: fmt.Sprintf("The device %s exists", device.ID)}
	WriteAPIResponse(w, http.StatusCreated, response)

}

func (s *SignatureService) AddSignatureDevice(device *domain.SignatureDevice) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.devices[device.ID] = device
}

func (s *SignatureService) GetSignatureDevice(deviceID uuid.UUID) (*domain.SignatureDevice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, exists := s.devices[deviceID]
	return device, exists
}

func encodeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func HandleError(w http.ResponseWriter, status int, err error, messages ...string) bool {
	if err != nil {
		WriteErrorResponse(w, status, messages)
		return true
	}
	return false
}
