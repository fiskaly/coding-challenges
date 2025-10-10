package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// setupTestServer creates a test server with in-memory storage
func setupTestServer() *Server {
	repo := persistence.NewInMemoryDeviceRepository()
	deviceService := service.NewDeviceService(repo)
	return NewServer(":8080", deviceService)
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v0/health", nil)
	w := httptest.NewRecorder()

	server.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
}

// TestCreateDeviceEndpoint tests device creation via API
func TestCreateDeviceEndpoint(t *testing.T) {
	server := setupTestServer()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "Valid RSA device",
			requestBody: map[string]interface{}{
				"algorithm": "RSA",
				"label":     "Test Device",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Valid ECC device",
			requestBody: map[string]interface{}{
				"algorithm": "ECC",
				"label":     "Test Device",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Device without label",
			requestBody: map[string]interface{}{
				"algorithm": "RSA",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid algorithm",
			requestBody: map[string]interface{}{
				"algorithm": "INVALID",
				"label":     "Test Device",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			server.CreateDeviceHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			// For successful creation, verify UUID was generated
			if tt.expectedStatus == http.StatusCreated {
				var response Response
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				deviceData, ok := response.Data.(map[string]interface{})
				if !ok {
					t.Fatal("Expected data to contain device object")
				}

				device, ok := deviceData["device"].(map[string]interface{})
				if !ok {
					t.Fatal("Expected device object in response")
				}

				deviceID, ok := device["id"].(string)
				if !ok || deviceID == "" {
					t.Error("Expected device ID to be generated and non-empty")
				}
			}
		})
	}
}

// TestListDevicesEndpoint tests listing devices via API
func TestListDevicesEndpoint(t *testing.T) {
	server := setupTestServer()

	// Create a few devices first (no ID in request anymore)
	for i := 0; i < 3; i++ {
		reqBody := map[string]interface{}{
			"algorithm": "RSA",
			"label":     "Test Device",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		server.CreateDeviceHandler(w, req)
	}

	// List devices
	req := httptest.NewRequest(http.MethodGet, "/api/v0/devices", nil)
	w := httptest.NewRecorder()

	server.ListDevicesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	devices, ok := response.Data.([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(devices) != 3 {
		t.Errorf("Expected 3 devices, got %d", len(devices))
	}
}

// TestGetDeviceEndpoint tests getting a specific device via API
func TestGetDeviceEndpoint(t *testing.T) {
	server := setupTestServer()

	// Create a device (ID is auto-generated)
	reqBody := map[string]interface{}{
		"algorithm": "RSA",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.CreateDeviceHandler(w, req)

	// Extract the generated device ID from response
	var createResponse Response
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	deviceData := createResponse.Data.(map[string]interface{})
	device := deviceData["device"].(map[string]interface{})
	deviceID := device["id"].(string)

	// Get the device
	req = httptest.NewRequest(http.MethodGet, "/api/v0/devices/"+deviceID, nil)
	w = httptest.NewRecorder()

	server.GetDeviceHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test non-existing device
	req = httptest.NewRequest(http.MethodGet, "/api/v0/devices/non-existing", nil)
	w = httptest.NewRecorder()

	server.GetDeviceHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestSignTransactionEndpoint tests signing via API
func TestSignTransactionEndpoint(t *testing.T) {
	server := setupTestServer()

	// Create a device first (ID is auto-generated)
	deviceReqBody := map[string]interface{}{
		"algorithm": "RSA",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(deviceReqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.CreateDeviceHandler(w, req)

	// Extract device ID
	var createResponse Response
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	deviceData := createResponse.Data.(map[string]interface{})
	device := deviceData["device"].(map[string]interface{})
	deviceID := device["id"].(string)

	// Sign a transaction
	signReqBody := map[string]interface{}{
		"device_id": deviceID,
		"data":      "test_transaction_data",
	}
	body, _ = json.Marshal(signReqBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v0/signatures", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	server.SignTransactionHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be an object")
	}

	if data["signature"] == nil || data["signature"].(string) == "" {
		t.Error("Signature should not be empty")
	}

	if data["signed_data"] == nil || data["signed_data"].(string) == "" {
		t.Error("Signed data should not be empty")
	}
}

// TestSignTransactionChaining tests that signatures are properly chained
func TestSignTransactionChaining(t *testing.T) {
	server := setupTestServer()

	// Create a device (ID is auto-generated)
	deviceReqBody := map[string]interface{}{
		"algorithm": "ECC",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(deviceReqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	server.CreateDeviceHandler(w, req)

	// Extract device ID
	var createResponse Response
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	deviceData := createResponse.Data.(map[string]interface{})
	device := deviceData["device"].(map[string]interface{})
	deviceID := device["id"].(string)

	// Sign first transaction
	signReq1 := map[string]interface{}{
		"device_id": deviceID,
		"data":      "transaction_1",
	}
	body, _ = json.Marshal(signReq1)
	req = httptest.NewRequest(http.MethodPost, "/api/v0/signatures", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.SignTransactionHandler(w, req)

	var response1 Response
	json.Unmarshal(w.Body.Bytes(), &response1)
	data1 := response1.Data.(map[string]interface{})
	signature1 := data1["signature"].(string)

	// Sign second transaction
	signReq2 := map[string]interface{}{
		"device_id": deviceID,
		"data":      "transaction_2",
	}
	body, _ = json.Marshal(signReq2)
	req = httptest.NewRequest(http.MethodPost, "/api/v0/signatures", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	server.SignTransactionHandler(w, req)

	var response2 Response
	json.Unmarshal(w.Body.Bytes(), &response2)
	data2 := response2.Data.(map[string]interface{})
	signedData2 := data2["signed_data"].(string)

	// Verify that second signed data contains first signature
	expectedPrefix := "1_transaction_2_" + signature1
	if signedData2 != expectedPrefix {
		t.Errorf("Expected signed_data to contain previous signature\ngot:  %s\nwant: %s", signedData2, expectedPrefix)
	}
}

// TestMethodNotAllowed tests that wrong HTTP methods are rejected
func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"POST health", http.MethodPost, "/api/v0/health"},
		{"DELETE devices", http.MethodDelete, "/api/v0/devices"},
		{"PUT signatures", http.MethodPut, "/api/v0/signatures"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			switch tt.path {
			case "/api/v0/health":
				server.Health(w, req)
			case "/api/v0/devices":
				server.ListDevicesHandler(w, req)
			case "/api/v0/signatures":
				server.SignTransactionHandler(w, req)
			}

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405, got %d", w.Code)
			}
		})
	}
}

// TestSignWithNonExistingDevice tests signing with non-existing device
func TestSignWithNonExistingDevice(t *testing.T) {
	server := setupTestServer()

	signReq := map[string]interface{}{
		"device_id": "non-existing-device",
		"data":      "test_data",
	}
	body, _ := json.Marshal(signReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v0/signatures", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	server.SignTransactionHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
