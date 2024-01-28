package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreateSignatureDeviceResponse(t *testing.T) {
	t.Run("fails when method is not POST", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/v0/signature_devices", nil)
		responseRecorder := httptest.NewRecorder()

		service := api.NewSignatureService(persistence.NewInMemorySignatureDeviceRepository())
		service.CreateSignatureDevice(responseRecorder, request)

		expectedStatusCode := http.StatusMethodNotAllowed
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["Method Not Allowed"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when uuid is invalid", func(t *testing.T) {
		id := "invalid-uuid"
		algorithmName := "RSA"
		request := newJsonRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id,
				Algorithm: algorithmName,
			},
		)
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["id is not a valid uuid"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when id already exists", func(t *testing.T) {
		id := uuid.New()
		generator := crypto.RSAGenerator{}
		request := newJsonRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: generator.AlgorithmName(),
			},
		)
		responseRecorder := httptest.NewRecorder()

		keyPair, err := generator.Generate()
		if err != nil {
			t.Fatal(err)
		}

		repository := persistence.NewInMemorySignatureDeviceRepository()
		// create existing device with the same id
		repository.Create(domain.SignatureDevice{
			ID:      id,
			KeyPair: keyPair,
		})
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["duplicate id"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when algorithm is invalid", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "ABC"
		request := newJsonRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := `{"errors":["algorithm is not supported"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("creates a SignatureDevice successfully", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "RSA"
		request := newJsonRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusCreated
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := fmt.Sprintf(`{
  "data": {
    "signatureDeviceId": "%s"
  }
}`, id)
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}

		// check persisted data
		device, found, err := repository.Find(id)
		if err != nil {
			t.Error(err)
		}
		if !found {
			t.Error("expected device with id to be found")
		}
		if device.ID != id {
			t.Errorf("id not persisted correctly. expected: %s, got: %s", id, device.ID)
		}
		if device.Label != "" {
			t.Errorf("label not persisted correctly. expected blank string, got: %s", device.Label)
		}
		_, ok := device.KeyPair.(*crypto.RSAKeyPair)
		if !ok {
			t.Errorf("key pair generation failed: %s", err)
		}
	})

	t.Run("creates a SignatureDevice with a label successfully", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "RSA"
		label := "my RSA key"
		request := newJsonRequest(
			http.MethodPost,
			"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
				Label:     label,
			},
		)
		responseRecorder := httptest.NewRecorder()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		service := api.NewSignatureService(repository)
		service.CreateSignatureDevice(responseRecorder, request)

		// check status code
		expectedStatusCode := http.StatusCreated
		if responseRecorder.Code != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, responseRecorder.Code)
		}

		// check body
		body := responseRecorder.Body.String()
		expectedBody := fmt.Sprintf(`{
  "data": {
    "signatureDeviceId": "%s"
  }
}`, id)
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}

		// check persisted data
		device, found, err := repository.Find(id)
		if err != nil {
			t.Error(err)
		}
		if !found {
			t.Error("expected device with id to be found")
		}
		if device.ID != id {
			t.Errorf("id not persisted correctly. expected: %s, got: %s", id, device.ID)
		}
		if device.Label != label {
			t.Errorf("label not persisted correctly. expected: %s, got: %s", label, device.Label)
		}
		_, ok := device.KeyPair.(*crypto.RSAKeyPair)
		if !ok {
			t.Errorf("key pair generation failed: %s", err)
		}
	})
}

func newJsonRequest(httpMethod string, path string, serializableData any) *http.Request {
	jsonBytes, err := json.Marshal(serializableData)
	if err != nil {
		panic(fmt.Sprintf("json.Marshal failed: err"))
	}

	request := httptest.NewRequest(
		httpMethod,
		path,
		bytes.NewReader(jsonBytes),
	)
	request.Header.Set("Content-Type", "application/json")
	return request
}
