package api_test

import (
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
	t.Run("fails when uuid is invalid", func(t *testing.T) {
		id := "invalid-uuid"
		algorithmName := crypto.RSAGenerator{}.AlgorithmName()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(repository)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id,
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["id is not a valid uuid"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when id already exists", func(t *testing.T) {
		id := uuid.New()

		// create existing device with the id
		generator := crypto.RSAGenerator{}
		keyPair, err := generator.Generate()
		if err != nil {
			t.Fatal(err)
		}
		repository := persistence.NewInMemorySignatureDeviceRepository()
		repository.Create(domain.SignatureDevice{
			ID:      id,
			KeyPair: keyPair,
		})

		signatureService := api.NewSignatureService(repository)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: generator.AlgorithmName(),
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["duplicate id"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("fails when algorithm is invalid", func(t *testing.T) {
		id := uuid.New()
		algorithmName := "ABC"

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(repository)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusBadRequest
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["algorithm is not supported"]}`
		if body != expectedBody {
			t.Errorf("expected: %s, got: %s", expectedBody, body)
		}
	})

	t.Run("creates a SignatureDevice successfully", func(t *testing.T) {
		id := uuid.New()
		algorithmName := crypto.RSAGenerator{}.AlgorithmName()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(repository)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
			},
		)

		// check status code
		expectedStatusCode := http.StatusCreated
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
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

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(repository)
		server := httptest.NewServer(api.NewServer("", signatureService).HTTPHandler())
		defer server.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			server.URL+"/api/v0/signature_devices",
			api.CreateSignatureDeviceRequest{
				ID:        id.String(),
				Algorithm: algorithmName,
				Label:     label,
			},
		)

		// check status code
		expectedStatusCode := http.StatusCreated
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
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
