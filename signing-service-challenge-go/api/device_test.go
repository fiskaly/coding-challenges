package api_test

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
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
    "signature_device_id": "%s"
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
    "signature_device_id": "%s"
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

func TestSignTransaction(t *testing.T) {
	t.Run("returns not found when device with id does not exist", func(t *testing.T) {
		id := uuid.NewString()

		repository := persistence.NewInMemorySignatureDeviceRepository()
		signatureService := api.NewSignatureService(repository)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: "some-data"},
		)

		// check status code
		expectedStatusCode := http.StatusNotFound
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// check body
		body := readBody(t, response)
		expectedBody := `{"errors":["signature device not found"]}`
		diff := cmp.Diff(body, expectedBody)
		if diff != "" {
			t.Errorf("unexpected diff: %s", diff)
		}
	})

	t.Run("successfully signs data with device (algorithm: RSA, counter = 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		base64EncodedID := "NjRmZjc5NmUtZmNkZS00OTlhLWEwM2QtODJkZDFmODllOGU1"
		dataToSign := "some-data"
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.RSAGenerator{})
		if err != nil {
			t.Fatal(err)
		}

		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(repository)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.RSAKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		err = rsa.VerifyPSS(keyPair.Public, crypto.HashFunction, digest, decodedSignature, nil)
		if err != nil {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("0_%s_%s", dataToSign, base64EncodedID)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 1 {
			t.Errorf("device signature counter should be incremented to 1, got: %d", device.SignatureCounter)
		}
		if device.Base64EncodedLastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.Base64EncodedLastSignature)
		}
	})

	t.Run("successfully signs data (algorithm: RSA, counter > 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		dataToSign := "some-data"

		// create a device that has been used once
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.RSAGenerator{})
		if err != nil {
			t.Fatal(err)
		}
		device.SignatureCounter = 1
		device.Base64EncodedLastSignature = "last-signature-base-64-encoded"
		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(repository)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.RSAKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		err = rsa.VerifyPSS(keyPair.Public, crypto.HashFunction, digest, decodedSignature, nil)
		if err != nil {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("1_%s_%s", dataToSign, device.Base64EncodedLastSignature)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 2 {
			t.Errorf("device signature counter should be incremented to 2, got: %d", device.SignatureCounter)
		}
		if device.Base64EncodedLastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.Base64EncodedLastSignature)
		}
	})

	t.Run("successfully signs data with device (algorithm: ECC, counter = 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		base64EncodedID := "NjRmZjc5NmUtZmNkZS00OTlhLWEwM2QtODJkZDFmODllOGU1"
		dataToSign := "some-data"
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.ECCGenerator{})
		if err != nil {
			t.Fatal(err)
		}

		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(repository)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		keyPair := device.KeyPair.(*crypto.ECCKeyPair)
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		result := ecdsa.VerifyASN1(keyPair.Public, digest, decodedSignature)
		if !result {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("0_%s_%s", dataToSign, base64EncodedID)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 1 {
			t.Errorf("device signature counter should be incremented to 1, got: %d", device.SignatureCounter)
		}
		if device.Base64EncodedLastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.Base64EncodedLastSignature)
		}
	})

	t.Run("successfully signs data (algorithm: ECC, counter > 0)", func(t *testing.T) {
		id := "64ff796e-fcde-499a-a03d-82dd1f89e8e5"
		dataToSign := "some-data"

		// create a device that has been used once
		device, err := domain.BuildSignatureDevice(uuid.MustParse(id), crypto.ECCGenerator{})
		if err != nil {
			t.Fatal(err)
		}
		device.SignatureCounter = 1
		device.Base64EncodedLastSignature = "last-signature-base-64-encoded"
		repository := persistence.NewInMemorySignatureDeviceRepository()
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		signatureService := api.NewSignatureService(repository)
		testServer := httptest.NewServer(api.NewServer(":8888", signatureService).HTTPHandler())
		defer testServer.Close()

		response := sendJsonRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("%s/api/v0/signature_devices/%s/signatures", testServer.URL, id),
			api.SignTransactionRequest{DataToBeSigned: dataToSign},
		)

		// check status code
		expectedStatusCode := http.StatusOK
		if response.StatusCode != expectedStatusCode {
			t.Errorf("expected status code: %d, got: %d", expectedStatusCode, response.StatusCode)
		}

		// unmarshal body
		body := readBody(t, response)
		jsonBody := struct {
			Data api.SignTransactionResponse `json:"data"`
		}{}
		err = json.Unmarshal([]byte(body), &jsonBody)
		if err != nil {
			t.Errorf("unexpected response body format: %s", err)
		}

		// check signature is verifiable
		digest, err := crypto.ComputeHashDigest([]byte(jsonBody.Data.SignedData))
		if err != nil {
			t.Fatal(err)
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(jsonBody.Data.Signature)
		if err != nil {
			t.Fatal(err)
		}
		keyPair := device.KeyPair.(*crypto.ECCKeyPair)
		result := ecdsa.VerifyASN1(keyPair.Public, digest, decodedSignature)
		if !result {
			t.Errorf("verification of signed data and signature failed. err: %s, signed data: %s, signature: %s", err, jsonBody.Data.SignedData, jsonBody.Data.Signature)
		}

		// check signed_data is correct format
		expectedSignedData := fmt.Sprintf("1_%s_%s", dataToSign, device.Base64EncodedLastSignature)
		if jsonBody.Data.SignedData != expectedSignedData {
			t.Errorf("expected signed data: %s, got: %s", expectedSignedData, jsonBody.Data.SignedData)
		}

		// check persisted data
		device, ok, err := repository.Find(uuid.MustParse(id))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 2 {
			t.Errorf("device signature counter should be incremented to 2, got: %d", device.SignatureCounter)
		}
		if device.Base64EncodedLastSignature != jsonBody.Data.Signature {
			t.Errorf("device last signature should be updated to %s, got: %s", jsonBody.Data.Signature, device.Base64EncodedLastSignature)
		}
	})
}
