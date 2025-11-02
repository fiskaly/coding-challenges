package domain

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
)

func TestNewSignatureDeviceInitialState(t *testing.T) {
	deviceID := "device-123"
	device, err := NewSignatureDevice(deviceID, AlgorithmRSA, "POS Device")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if device.ID() != deviceID {
		t.Fatalf("expected id %q, got %q", deviceID, device.ID())
	}

	if device.SignatureCounter() != 0 {
		t.Fatalf("expected signature counter 0, got %d", device.SignatureCounter())
	}

	expectedLast := base64.StdEncoding.EncodeToString([]byte(deviceID))
	if device.LastSignature() != expectedLast {
		t.Fatalf("expected last signature %q, got %q", expectedLast, device.LastSignature())
	}
}

func TestSignatureDeviceSecuredPayload(t *testing.T) {
	deviceID := "terminal-001"
	device, err := NewSignatureDevice(deviceID, AlgorithmECC, "Checkout terminal")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	payload, err := device.SecuredPayload("total=42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := fmt.Sprintf("0_total=42_%s", base64.StdEncoding.EncodeToString([]byte(deviceID)))
	if payload != expected {
		t.Fatalf("expected payload %q, got %q", expected, payload)
	}
}

func TestRecordSignatureUpdatesState(t *testing.T) {
	device, err := NewSignatureDevice("device-xyz", AlgorithmRSA, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rawSignature := []byte("signed bytes")
	encoded, err := device.RecordSignature(rawSignature)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedEncoded := base64.StdEncoding.EncodeToString(rawSignature)
	if encoded != expectedEncoded {
		t.Fatalf("expected encoded signature %q, got %q", expectedEncoded, encoded)
	}

	if device.SignatureCounter() != 1 {
		t.Fatalf("expected counter to be 1, got %d", device.SignatureCounter())
	}

	if device.LastSignature() != expectedEncoded {
		t.Fatalf("expected last signature to be %q, got %q", expectedEncoded, device.LastSignature())
	}

	nextPayload, err := device.SecuredPayload("ref=42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPayload := fmt.Sprintf("1_ref=42_%s", expectedEncoded)
	if nextPayload != expectedPayload {
		t.Fatalf("expected payload %q, got %q", expectedPayload, nextPayload)
	}
}

func TestSecuredPayloadRejectsEmptyData(t *testing.T) {
	device, err := NewSignatureDevice("device-xyz", AlgorithmRSA, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = device.SecuredPayload("  ")
	if !errors.Is(err, ErrEmptyPayload) {
		t.Fatalf("expected ErrEmptyPayload, got %v", err)
	}
}

func TestRecordSignatureRejectsEmptyInput(t *testing.T) {
	device, err := NewSignatureDevice("device-xyz", AlgorithmRSA, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = device.RecordSignature(nil)
	if !errors.Is(err, ErrEmptySignature) {
		t.Fatalf("expected ErrEmptySignature, got %v", err)
	}
}

func TestParseAlgorithm(t *testing.T) {
	cases := []struct {
		input string
		want  Algorithm
		err   error
	}{
		{"rsa", AlgorithmRSA, nil},
		{"ECC", AlgorithmECC, nil},
		{"  ecc  ", AlgorithmECC, nil},
		{"unknown", "", ErrUnsupportedAlgorithm},
	}

	for _, tc := range cases {
		got, err := ParseAlgorithm(tc.input)
		if tc.err != nil {
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected error %v for input %q, got %v", tc.err, tc.input, err)
			}
			continue
		}

		if err != nil {
			t.Fatalf("expected no error for input %q, got %v", tc.input, err)
		}

		if got != tc.want {
			t.Fatalf("expected algorithm %q, got %q", tc.want, got)
		}
	}
}

func TestRestoreSignatureDeviceDefaultsLastSignature(t *testing.T) {
	device, err := RestoreSignatureDevice("device-restore", AlgorithmRSA, "restored", 5, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if device.SignatureCounter() != 5 {
		t.Fatalf("expected counter 5, got %d", device.SignatureCounter())
	}

	expectedLast := base64.StdEncoding.EncodeToString([]byte("device-restore"))
	if device.LastSignature() != expectedLast {
		t.Fatalf("expected last signature %q, got %q", expectedLast, device.LastSignature())
	}
}
