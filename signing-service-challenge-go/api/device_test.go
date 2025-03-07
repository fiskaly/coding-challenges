package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type SignatureDomainStub struct {
	CreateSignatureDeviceFunc     func(domain.Device) error
	ListSignatureDevicesFunc      func() ([]domain.ID, error)
	SignTransactionFunc           func(domain.ID, domain.Data) (*domain.CreatedSignature, error)
	GetSignatureDeviceDetailsFunc func(domain.ID) (domain.Device, error)
}

func (s *SignatureDomainStub) CreateSignatureDevice(d domain.Device) error {
	return s.CreateSignatureDeviceFunc(d)
}

func (s *SignatureDomainStub) ListSignatureDevices() ([]domain.ID, error) {
	return s.ListSignatureDevicesFunc()
}

func (s *SignatureDomainStub) SignTransaction(id domain.ID, data domain.Data) (*domain.CreatedSignature, error) {
	return s.SignTransactionFunc(id, data)
}

func (s *SignatureDomainStub) GetSignatureDeviceDetails(id domain.ID) (domain.Device, error) {
	return s.GetSignatureDeviceDetailsFunc(id)
}

func TestCreateSignatureDevice(t *testing.T) {
	d := domain.SignatureDomain(
		&SignatureDomainStub{
			CreateSignatureDeviceFunc: func(d domain.Device) error {
				return nil
			},
		})
	mux := SignatureDeviceRoutes(&d)
	server := httptest.NewServer(mux)
	defer server.Close()

	reqBody := CreateSignatureDeviceRequest{
		SignatureDeviceDetails: SignatureDeviceDetails{
			ID:        "dev01",
			Algorithm: "RSA",
			Label:     "Test Device",
		},
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(server.URL+"/signature-devices", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %v, got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestSignTransaction(t *testing.T) {
	d := domain.SignatureDomain(&SignatureDomainStub{
		SignTransactionFunc: func(id domain.ID, data domain.Data) (*domain.CreatedSignature, error) {
			return &domain.CreatedSignature{
				Signature:  "signature",
				SignedData: "1_" + string(data) + "last_signature",
			}, nil
		},
	})
	mux := SignatureDeviceRoutes(&d)
	server := httptest.NewServer(mux)
	defer server.Close()

	// First, create a signature device
	reqBody := CreateSignatureDeviceRequest{
		SignatureDeviceDetails: SignatureDeviceDetails{
			Algorithm: "RSA",
			Label:     "Test Device",
		},
	}
	body, _ := json.Marshal(reqBody)
	http.Post(server.URL+"/signature-devices", "application/json", bytes.NewBuffer(body))

	// Now, sign a transaction
	signReqBody := SignTransactionRequest{
		Data: domain.Data("testdata"),
	}
	signBody, _ := json.Marshal(signReqBody)

	resp, err := http.Post(server.URL+"/signature-devices/1/sign-transaction", "application/json", bytes.NewBuffer(signBody))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	var signRespBody SignatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&signRespBody); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedSignature := "signature"
	if signRespBody.Signature != expectedSignature {
		t.Fatalf("expected %v, got %v", expectedSignature, signRespBody.Signature)
	}
}

func TestListSignatureDevices(t *testing.T) {
	d := domain.SignatureDomain(&SignatureDomainStub{
		ListSignatureDevicesFunc: func() ([]domain.ID, error) {
			return []domain.ID{"1"}, nil
		},
	})
	mux := SignatureDeviceRoutes(&d)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a signature device
	reqBody := CreateSignatureDeviceRequest{
		SignatureDeviceDetails: SignatureDeviceDetails{
			Algorithm: "RSA",
			Label:     "Test Device",
		},
	}
	body, _ := json.Marshal(reqBody)
	http.Post(server.URL+"/signature-devices", "application/json", bytes.NewBuffer(body))

	// List signature devices
	resp, err := http.Get(server.URL + "/signature-devices")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	var listRespBody ListSignatureDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&listRespBody); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(listRespBody.IDs) != 1 {
		t.Fatalf("expected 1 device, got %v", len(listRespBody.IDs))
	}
}

func TestGetSignatureDeviceDetails(t *testing.T) {
	d := domain.SignatureDomain(&SignatureDomainStub{
		GetSignatureDeviceDetailsFunc: func(domain.ID) (domain.Device, error) {
			return domain.Device{
				Algorithm: "RSA",
				Label:     "Test Device",
			}, nil
		},
	})
	mux := SignatureDeviceRoutes(&d)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a signature device
	reqBody := CreateSignatureDeviceRequest{
		SignatureDeviceDetails: SignatureDeviceDetails{
			Algorithm: "RSA",
			Label:     "Test Device",
		},
	}
	body, _ := json.Marshal(reqBody)
	http.Post(server.URL+"/signature-devices", "application/json", bytes.NewBuffer(body))

	// Get signature device details
	resp, err := http.Get(server.URL + "/signature-devices/1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	var getRespBody GetSignatureDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&getRespBody); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if getRespBody.Algorithm != "RSA" || getRespBody.Label != "Test Device" {
		t.Fatalf("expected Algorithm: RSA, Label: Test Device, got Algorithm: %v, Label: %v", getRespBody.Algorithm, getRespBody.Label)
	}
}
