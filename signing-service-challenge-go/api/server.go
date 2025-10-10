package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
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
	deviceService *service.DeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService *service.DeviceService) *Server {
	return &Server{
		listenAddress: listenAddress,
		deviceService: deviceService,
	}
}

// Handler returns the HTTP handler with all routes configured
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	mux.HandleFunc("/api/v0/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.CreateDeviceHandler(w, r)
		} else if r.Method == http.MethodGet {
			s.ListDevicesHandler(w, r)
		} else {
			WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
				http.StatusText(http.StatusMethodNotAllowed),
			})
		}
	})

	mux.HandleFunc("/api/v0/devices/", s.GetDeviceHandler)

	mux.Handle("/api/v0/signatures", http.HandlerFunc(s.SignTransactionHandler))

	return mux
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
// Deprecated: Use Handler() method instead for better control over server lifecycle
func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddress, s.Handler())
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
