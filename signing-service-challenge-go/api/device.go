package api

import (
	"encoding/json"
	"net/http"
	"signing-service-challenge/domain"
	"signing-service-challenge/helper"
	"signing-service-challenge/service"

	"github.com/gorilla/mux"
)

type DeviceHandler struct {
	deviceService service.DeviceService
}

func NewDeviceHandler(deviceService service.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService}
}

type CreateSignatureDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type CreateSignatureDeviceResponse struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

// TODO: REST endpoints ...
func (s *DeviceHandler) CreateSignatureDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateSignatureDeviceRequest
	var device *domain.Device
	w.Header().Set("Content-type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid Payload."})
	}

	if req.Algorithm == "" || req.Label == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required field in payload: algorithm, label."})
		return
	}

	device, err := s.deviceService.CreateDevice(req.Label, domain.AlgorithmType(req.Algorithm))
	if err != nil {
		code, msg := helper.HandleDeviceServiceError(err)
		WriteErrorResponse(w, code, []string{msg})
		return
	}

	response := CreateSignatureDeviceResponse{
		Label:     device.Label,
		Algorithm: string(device.Algorithm),
	}

	WriteAPIResponse(w, http.StatusCreated, response)

}

func (s *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.deviceService.ListDevices()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Failed to retrieve list of devices"})
		return
	}

	response := make([]DeviceListResponse, len(devices))
	for i, device := range devices {
		response[i] = DeviceListResponse{
			Id:        device.Id,
			Label:     device.Label,
			Algorithm: string(device.Algorithm),
		}
	}
	WriteAPIResponse(w, http.StatusOK, response)
}

func (s *DeviceHandler) GetDeviceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]

	if deviceId == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required field in parameters: deviceId."})
		return
	}

	device, err := s.deviceService.GetDeviceById(deviceId)
	if err != nil {
		code, msg := helper.HandleDeviceServiceError(err)
		WriteErrorResponse(w, code, []string{msg})
		return
	}

	response := DeviceByIdResponse{
		Id:               device.Id,
		Label:            device.Label,
		Algorithm:        string(device.Algorithm),
		SignatureCounter: device.SignatureCounter,
	}
	WriteAPIResponse(w, http.StatusOK, response)
}
