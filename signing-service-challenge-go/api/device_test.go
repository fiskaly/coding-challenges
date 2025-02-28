package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
)

func setupTestServer() *httptest.Server {
	store := persistence.NewInMemoryStore()
	handler := NewDeviceHandler(store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

// Edge Case: Empty Payload
func TestCreateSignatureDevice_EmptyPayload(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", resp.StatusCode)
	}
}

// Edge Case: Invalid JSON
func TestCreateSignatureDevice_InvalidJSON(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payload := []byte(`{"algorithm": "RSA", "label": "Test Device"`) // Missing closing brace
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", resp.StatusCode)
	}
}

// Edge Case: Unsupported HTTP Method
func TestCreateSignatureDevice_UnsupportedMethod(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v0/devices", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code 405, got %d", resp.StatusCode)
	}
}

// Performance Test: High Volume Signing
func TestSignData_HighVolume(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payload := []byte(`{"algorithm": "RSA", "label": "High Volume Device"}`)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var wrappedResponse struct {
		Data CreateSignatureDeviceResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrappedResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	deviceID := wrappedResponse.Data.ID

	var wg sync.WaitGroup
	numRequests := 1000
	wg.Add(numRequests)

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()

			signPayload := []byte(`{"data": "test_data"}`)
			req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices/"+deviceID+"/sign", bytes.NewBuffer(signPayload))
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("Failed to send request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code 200, got %d", resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("High volume signing completed in %s", elapsed)
}

func TestCreateSignatureDevice_LabelSanitization(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payload := map[string]string{
		"algorithm": "RSA",
		"label":     "Test label'); DROP TABLE devices; --",
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var wrappedResponse struct {
		Data CreateSignatureDeviceResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrappedResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	deviceLabel := wrappedResponse.Data.Label
	expectedLabel := "Test label DROP TABLE devices --"
	if deviceLabel != expectedLabel {
		t.Errorf("Expected label to be '%s', got '%s'", expectedLabel, deviceLabel)
	}
}

func TestSignTransaction_InvalidDataLength(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create Device
	payload := []byte(`{"algorithm": "RSA", "label": "Invalid Data Length Device"}`)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var wrappedResponse struct {
		Data CreateSignatureDeviceResponse `json:"data"`
	}
	json.Unmarshal(body, &wrappedResponse)
	deviceID := wrappedResponse.Data.ID

	signPayload := []byte(`{"data": ""}`)
	req, err = http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices/"+deviceID+"/sign", bytes.NewBuffer(signPayload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", response.StatusCode)
	}
}

// GetDeviceDetails: Non-existent Device
func TestGetDeviceDetails_NotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	id := uuid.New().String()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v0/devices/"+id, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}
}

func TestListDevices_Empty(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v0/devices", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var wrappedResponse struct {
		Data []CreateSignatureDeviceResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrappedResponse); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(wrappedResponse.Data) != 0 {
		t.Errorf("Expected empty device list, got %d devices", len(wrappedResponse.Data))
	}

}

func TestListDevices_Populated(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payload := []byte(`{"algorithm": "RSA", "label": "Listed Device"}`)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v0/devices", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	http.DefaultClient.Do(req)

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/api/v0/devices", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
