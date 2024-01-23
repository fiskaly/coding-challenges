package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/v0/health", nil)
	responseRecorder := httptest.NewRecorder()

	server := api.NewServer("")
	server.Health(responseRecorder, request)

	expectedStatusCode := http.StatusOK
	if responseRecorder.Code != expectedStatusCode {
		t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
	}

	body := responseRecorder.Body.String()
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
