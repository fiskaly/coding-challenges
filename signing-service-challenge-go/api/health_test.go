package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func TestHealth(t *testing.T) {
	repository := persistence.NewInMemorySignatureDeviceRepository()
	signatureService := api.NewSignatureService(repository)
	server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
	defer server.Close()

	response := sendJsonRequest(
		t,
		http.MethodGet,
		server.URL+"/api/v0/health",
	)

	expectedStatusCode := http.StatusOK
	if response.StatusCode != expectedStatusCode {
		t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
	}

	body := readBody(t, response)
	expectedBody := `{
  "data": {
    "status": "pass",
    "version": "v0"
  }
}`
	if body != expectedBody {
		t.Errorf("expected: %s, got: %s", expectedBody, body)
	}
}
