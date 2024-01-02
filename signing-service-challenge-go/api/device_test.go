// main_test.go

package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSignatureDevice(t *testing.T) {
	service := NewSignatureService()

	reqBody := bytes.NewBufferString(`{"id": "2e8d4895-04e3-4fbf-9ecf-1daee23138ca", "algorithm": "RSA", "label": "TestDevice"}`)
	req, err := http.NewRequest("POST", "https://localhost:8080/api/v0/createdevice", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	service.CreateSignatureDevice(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response CreateSignatureDeviceResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

}

func TestSignTransaction(t *testing.T) {
	service := NewSignatureService()

	service.CreateSignatureDevice(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "https://localhost:8080/api/v0/createdevice", bytes.NewBufferString(`{"id": "2e8d4895-04e3-4fbf-9ecf-1daee23138ca","algorithm": "ECC", "label": "TestDevice"}`)),
	)

	reqBody := bytes.NewBufferString(`{"deviceId":"2e8d4895-04e3-4fbf-9ecf-1daee23138ca", "data": "TestTransactionData"}`)
	req, err := http.NewRequest("POST", "https://localhost:8080/api/v0/checkdevice", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	service.SignTransaction(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response SignatureResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Signature == "" {
		t.Error("Expected non-empty signature in response, got empty")
	}
	if response.SignedData == "" {
		t.Error("Expected non-empty signed data in response, got empty")
	}

	lastSignature, err := base64.StdEncoding.DecodeString(response.Signature)
	if err != nil {
		t.Errorf("Error decoding last signature: %v", err)
	}
	fmt.Printf("%s", lastSignature)
}
