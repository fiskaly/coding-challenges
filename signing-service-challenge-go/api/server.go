package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
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
	store         persistence.Storage
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, readTimeout, writeTimeout, idleTimeout time.Duration) *Server {
	store := persistence.NewInMemoryStore()
	return &Server{
		listenAddress: listenAddress,
		store:         store,
		readTimeout:   readTimeout,
		writeTimeout:  writeTimeout,
		idleTimeout:   idleTimeout,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	deviceHandler := NewDeviceHandler(s.store)
	deviceHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         s.listenAddress,
		Handler:      mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}
	return server.ListenAndServe()
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
