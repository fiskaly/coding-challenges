package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress      string
	deviceHandler      *DeviceHandler
	transactionHandler *TransactionHandler
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceServiceHanlder *DeviceHandler, transactionHandler *TransactionHandler) *Server {

	return &Server{
		listenAddress:      listenAddress,
		transactionHandler: transactionHandler,
		deviceHandler:      deviceServiceHanlder,
		// TODO: add services / further dependencies here ...
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := mux.NewRouter()

	router.Handle("/api/v0/health", http.HandlerFunc(s.Health)).Methods("GET")

	// TODO: register further HandlerFuncs here ...
	//Device Handler
	router.Handle("/api/v0/devices/create-device", http.HandlerFunc(s.deviceHandler.CreateSignatureDevice)).Methods("POST")
	router.Handle("/api/v0/devices/list-devices", http.HandlerFunc(s.deviceHandler.ListDevices)).Methods("GET")
	router.Handle("/api/v0/devices/{deviceId}", http.HandlerFunc(s.deviceHandler.GetDeviceById)).Methods("GET")

	//Transaction Handler
	router.Handle("/api/v0/transactions/{deviceId}/sign", http.HandlerFunc(s.transactionHandler.SignTransactionHandler)).Methods("POST")

	return http.ListenAndServe(s.listenAddress, router)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
