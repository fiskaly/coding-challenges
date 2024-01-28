package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SignatureService struct {
	signatureDeviceRepository domain.SignatureDeviceRepository
}

func NewSignatureService(repository domain.SignatureDeviceRepository) SignatureService {
	return SignatureService{
		signatureDeviceRepository: repository,
	}
}

type CreateSignatureDeviceResponse struct {
	ID string `json:"signature_device_id"`
}

type CreateSignatureDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"` // optional
}

func (s *SignatureService) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var requestBody CreateSignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"invalid json",
		})
		return
	}

	id, err := uuid.Parse(requestBody.ID)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"id is not a valid uuid",
		})
		return
	}

	_, ok, err := s.signatureDeviceRepository.Find(id)
	if err != nil {
		WriteInternalError(response)
		return
	}
	if ok {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"duplicate id",
		})
		return
	}

	generator, found := crypto.FindKeyPairGenerator(requestBody.Algorithm)
	if !found {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"algorithm is not supported",
		})
		return
	}

	device, err := domain.BuildSignatureDevice(id, generator, requestBody.Label)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	err = s.signatureDeviceRepository.Create(device)
	if err != nil {
		// In a real application, this error would be logged and sent to an error notification service
		WriteInternalError(response)
		return
	}

	responseBody := CreateSignatureDeviceResponse{
		ID: device.ID.String(),
	}
	WriteAPIResponse(response, http.StatusCreated, responseBody)
}

type SignTransactionRequest struct {
	Data string `json:"data"`
}

type SignTransactionResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *SignatureService) SignTransaction(response http.ResponseWriter, request *http.Request) {
	deviceIDString := chi.URLParam(request, "deviceID")
	deviceID, err := uuid.Parse(deviceIDString)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"id is not a valid uuid",
		})
		return
	}

	device, ok, err := s.signatureDeviceRepository.Find(deviceID)
	if err != nil {
		WriteInternalError(response)
		return
	}
	if !ok {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			"signature device not found",
		})
		return
	}

	var requestBody SignTransactionRequest
	err = json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"invalid json",
		})
		return
	}

	encodedSignature, signedData, err := domain.SignTransaction(
		device,
		s.signatureDeviceRepository,
		requestBody.Data,
	)
	if err != nil {
		WriteInternalError(response)
		return
	}

	WriteAPIResponse(
		response,
		http.StatusOK,
		SignTransactionResponse{
			Signature:  encodedSignature,
			SignedData: signedData,
		},
	)
}
