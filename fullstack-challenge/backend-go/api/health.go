package api

import (
	"net/http"
)

// HealthResponse is the response for the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health is the HTTP handler for the health check endpoint.
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}

	WriteAPIResponse(w, http.StatusOK, response)
}
