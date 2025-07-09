package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// Handler of the Sign API. The passed UUID is used to retrieve the device
// from persistence, and then the device Signing method (RSA/ECC) is used to
// sign the passed data. The API returns the signature plus a string containing
// the signed data
func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {

	// get payload
	body, err := io.ReadAll(request.Body)
	if err != nil {
		errStr := "error while reading the payload"
		WriteErrorResponse(request, response, http.StatusInternalServerError, []string{
			http.StatusText(http.StatusInternalServerError),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// handle empty body
	if len(body) == 0 {
		errStr := "received an empty payload"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			errStr,
		})
		errStr = fmt.Sprintf("%s\n", errStr)
		log.Errorf(errStr)
		return
	}

	// unmarshal the payload
	payload := signPayload{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		errStr := "error while Unmarshaling the payload"
		WriteErrorResponse(request, response, http.StatusInternalServerError, []string{
			http.StatusText(http.StatusInternalServerError),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// validate payload
	if err := payload.verify(); err != nil {
		errStr := "invalid payload"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// sign the transaction
	signature, signedData, err := persistence.SignTransaction(payload.DeviceId, payload.DataToBeSigned)
	if err != nil {
		errStr := "could not sign the transaction"
		WriteErrorResponse(request, response, http.StatusInternalServerError, []string{
			http.StatusText(http.StatusInternalServerError),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// marshal into response
	WriteAPIResponse(request, response, http.StatusOK, signatureOut{
		Signature:  signature,
		SignedData: signedData,
	})

	log.Infof("Message succesfully signed by the device UUID '%s'", payload.DeviceId)
}

// structure used to model the payload expected to be passed to the signing API
type signPayload struct {
	DeviceId       string `json:"deviceId" validate:"required"`
	DataToBeSigned string `json:"dataToBeSigned" validate:"required"`
}

// method used to verify the validity of the signPayload structure
func (p signPayload) verify() error {

	// create validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	// validate
	if err := validate.Struct(p); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	log.Debug("succesfully verified the 'createSignatureDevicePayload' struct ")

	return nil
}

// structure used to model the payload of the response of the signing
// API
type signatureOut struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}
