package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	listenAddress string
	httpServer    *http.Server
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
		httpServer:    &http.Server{Addr: listenAddress},
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := mux.NewRouter()

	// - Doubt:
	// 		At first I assumed/understood that the Device UUIDs should have been generated automatically through the
	// 		package provided in the go.mod (I could not find it btw). Then, by looking at the TypeScript function
	// 		"CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse"
	// 		(assuming they act as the client), I could see that the ID was a required parameter to be passed to the
	// 		function. Therefore I could created two API Points, one that allows the client to explicitly specify the UUID
	// 		and another that returns the UUID automatically generated. This seems odd, UUID usually are managed/assigned
	// 		by a central service (this one). Maybe this is specific case, or it is only a oversight. I am assumin the latter
	// - Assumption:
	// 		The client does not specify the device UUID while creating a new device
	router.HandleFunc("/api/v0/health", s.Health).Methods("GET")
	router.HandleFunc("/api/v0/device", s.CreateSignatureDevice).Methods("POST")
	router.HandleFunc("/api/v0/device/{uuid}", s.GetDevice).Methods("GET")
	router.HandleFunc("/api/v0/devices", s.GetDevices).Methods("GET")
	router.HandleFunc("/api/v0/sign", s.SignTransaction).Methods("POST")

	s.httpServer.Handler = router

	log.Infof("Server listening on %s\n", s.listenAddress)
	return s.httpServer.ListenAndServe()
}

// Close the Server.
func (s *Server) Close() error {
	log.Info("Server closing\n")
	return s.httpServer.Close()
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(r *http.Request, w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	// log the error
	log.Errorf("[%d] '%s' '%s'\n", code, r.Method, r.URL)

	// in case the marshaling of the errorResponse fails
	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(r *http.Request, w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	// log the info
	log.Infof("[%d] '%s' '%s' Successfully served\n", code, r.Method, r.URL)

	// in case the marshaling of the errorResponse fails
	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
