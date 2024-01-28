package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	listenAddress    string
	signatureService SignatureService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(
	listenAddress string,
	signatureService SignatureService,
) *Server {
	return &Server{
		listenAddress:    listenAddress,
		signatureService: signatureService,
	}
}

// Register all HandlerFuncs for routes
func (s *Server) HTTPHandler() http.Handler {
	mux := chi.NewMux()
	mux.Get("/api/v0/health", http.HandlerFunc(s.Health))
	mux.Post("/api/v0/signature_devices", http.HandlerFunc(s.signatureService.CreateSignatureDevice))
	return mux
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.listenAddress, s.HTTPHandler())
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
