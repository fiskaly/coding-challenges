package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	crypto "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler of the CreateSignatureDevice API. It creates a new Signature Device
// with the provided label (optional) and the provided signing algorithm. It
// returns the UUID and the label of the created device
func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {

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
	payload := createSignatureDevicePayload{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		errStr := "error while Unmarshaling the payload"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
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

	// create new device
	dev, err := persistence.AddDevice(payload.SigningAlgorithm, payload.Label)
	if err != nil {
		errStr := "could not add a new device"
		WriteErrorResponse(request, response, http.StatusInternalServerError, []string{
			http.StatusText(http.StatusInternalServerError),
			errStr,
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// create struct to be passed back to the client
	devOut := deviceOut{
		Label: dev.Label,
		Uuid:  dev.UUID,
	}

	// marshal into response
	WriteAPIResponse(request, response, http.StatusOK, devOut)

	log.Infof("Create device with label '%s' and UUID '%s'", dev.Label, dev.UUID)

}

// Handler of the GetDevice API. It returns the device matching the provided
// UUID
func (s *Server) GetDevice(response http.ResponseWriter, request *http.Request) {

	// retrieve the id
	params := mux.Vars(request)
	uuid, found := params["uuid"]
	if !found || uuid == "" {
		errStr := "no device id provided in the request"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			errStr,
		})
		errStr = fmt.Sprintf("%s\n", errStr)
		log.Errorf(errStr)
		return
	}

	// get the device from persistence
	dev, err := persistence.GetDevice(uuid)
	if err != nil {
		// if not found
		if _, ok := err.(persistence.DeviceNotFoundError); ok {
			errStr := "requested device not found"
			WriteErrorResponse(request, response, http.StatusNotFound, []string{
				http.StatusText(http.StatusNotFound),
				errStr,
			})
			errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
			log.Errorf(errStr)
			return
		}

		// else
		errStr := "error while getting the requested device"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// create struct to be passed back to the client
	devOut := deviceOut{
		Label: dev.Label,
		Uuid:  dev.UUID,
	}

	// marshal into response
	WriteAPIResponse(request, response, http.StatusOK, devOut)

	log.Infof("Succesfully retrieved device with UUID '%s'", dev.UUID)

}

// Handler of the GetDevices API
func (s *Server) GetDevices(response http.ResponseWriter, request *http.Request) {

	// get the devices from persistence
	devs, err := persistence.GetDevices()
	if err != nil {
		errStr := "error while getting the devices"
		WriteErrorResponse(request, response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			errStr,
			err.Error(),
		})
		errStr = fmt.Sprintf("%s: %s\n", errStr, err.Error())
		log.Errorf(errStr)
		return
	}

	// create struct to be passed back to the client
	devsOut := make([]deviceOut, 0)
	for _, dev := range devs {
		devsOut = append(devsOut, deviceOut{
			Label: dev.Label,
			Uuid:  dev.UUID,
		})
	}

	// marshal into response
	WriteAPIResponse(request, response, http.StatusOK, devsOut)

	if len(devsOut) != 1 {
		log.Infof("Succesfully retrieved %d devices ", len(devsOut))
	} else {
		log.Infof("Succesfully retrieved %d device ", len(devsOut))
	}

}

// structure modelling the SignatureDevicePayload
type createSignatureDevicePayload struct {
	// See the assumption made in the server.go file
	// Id               string                           `json:"id" validate:"required"`
	SigningAlgorithm crypto.SigningAlgorithm `json:"signingAlgorithm" validate:"required,gte=1,lte=2"` // 0=ECC, 1=RSA
	Label            *string                 `json:"label,omitempty"`
}

// verify the validity of the createSignatureDevicePayload
func (p createSignatureDevicePayload) verify() error {

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

// support struct to allow decoupling between devices handled internally
// (by persistence/db) and what is exposed to the client. This should
// avoid exposing sensitive informaton
type deviceOut struct {
	Label string `json:"label"`
	Uuid  string `json:"uuid"`
}
