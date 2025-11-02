package api

import (
	"encoding/json"
	"net/http"
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
	deviceHandler *DeviceHandler
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceHandler *DeviceHandler) *Server {
	return &Server{
		listenAddress: listenAddress,
		deviceHandler: deviceHandler,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	if s.deviceHandler != nil {
		mux.Handle("/api/v0/devices", http.HandlerFunc(s.deviceHandler.HandleCollection))
		mux.Handle("/api/v0/devices/", http.HandlerFunc(s.deviceHandler.HandleResource))
	}

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// Best effort JSON error payload; ignore encoding errors because there is no reasonable recovery.
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Errors: []string{http.StatusText(http.StatusInternalServerError)},
	})
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	body, err := json.Marshal(ErrorResponse{
		Errors: errors,
	})
	if err != nil {
		WriteInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	body, err := json.MarshalIndent(Response{
		Data: data,
	}, "", "  ")
	if err != nil {
		WriteInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}
