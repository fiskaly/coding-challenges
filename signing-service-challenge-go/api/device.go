package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// CreateDeviceHandler handles POST /api/v0/devices
// Design Decision: Using POST with ID in the request body (client-provided UUID)
// This follows the specification: CreateSignatureDevice(id: string, ...)
func (s *Server) CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req service.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			"Invalid request body: " + err.Error(),
		})
		return
	}

	// Call service
	resp, err := s.deviceService.CreateDevice(req)
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if strings.Contains(errMsg, "already exists") {
			WriteErrorResponse(w, http.StatusConflict, []string{errMsg})
		} else if strings.Contains(errMsg, "unsupported algorithm") || strings.Contains(errMsg, "required") {
			WriteErrorResponse(w, http.StatusBadRequest, []string{errMsg})
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{errMsg})
		}
		return
	}

	WriteAPIResponse(w, http.StatusCreated, resp)
}

// GetDeviceHandler handles GET /api/v0/devices/{id}
func (s *Server) GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract device ID from path
	// Path format: /api/v0/devices/{id}
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Device ID is required"})
		return
	}
	deviceID := pathParts[3]

	device, err := s.deviceService.GetDevice(deviceID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		}
		return
	}

	WriteAPIResponse(w, http.StatusOK, device)
}

// ListDevicesHandler handles GET /api/v0/devices
func (s *Server) ListDevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	devices, err := s.deviceService.ListDevices()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusOK, devices)
}

// SignTransactionHandler handles POST /api/v0/signatures
// Design Decision: Using a separate resource endpoint for signatures
// This follows RESTful principles where signatures are a resource created from devices
func (s *Server) SignTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req service.SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			"Invalid request body: " + err.Error(),
		})
		return
	}

	resp, err := s.deviceService.SignTransaction(req)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			WriteErrorResponse(w, http.StatusNotFound, []string{errMsg})
		} else if strings.Contains(errMsg, "required") {
			WriteErrorResponse(w, http.StatusBadRequest, []string{errMsg})
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, []string{errMsg})
		}
		return
	}

	WriteAPIResponse(w, http.StatusCreated, resp)
}
