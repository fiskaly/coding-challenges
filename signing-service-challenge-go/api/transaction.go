package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"signing-service-challenge/helper"
	"signing-service-challenge/service"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

type SignTransactionRequest struct {
	Data string
}

type SignTransacitonResponse struct {
	Signature   string
	Signed_data string
}

func (s *TransactionHandler) SignTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid Payload."})
		return
	}

	if vars["deviceId"] == "" || req.Data == "" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Missing required field in payload: deviceId, data."})
		return
	}

	signature, dataToBeSigned, err := s.transactionService.SignTransaction(vars["deviceId"], req.Data)
	if err != nil {
		code, msg := helper.HandleDeviceServiceError(err)
		WriteErrorResponse(w, code, []string{msg})
		return
	}

	response := SignTransacitonResponse{
		Signature:   base64.RawStdEncoding.EncodeToString(signature),
		Signed_data: string(dataToBeSigned),
	}
	WriteAPIResponse(w, http.StatusOK, response)
}
