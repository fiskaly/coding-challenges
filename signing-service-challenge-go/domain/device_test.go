package domain_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

func TestSignWithRSA(t *testing.T) {
	id := uuid.NewString()
	device, err := domain.NewSignatureDevice(id, nil, crypto.AlgorithmECC)
	if err != nil {
		t.Fatal("failed to create signature device")
	}

	data := "data_to_be_signed"
	actual, err := device.Sign([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	lastSignature := base64.StdEncoding.EncodeToString(actual.Signature)

	expectedSignedData := fmt.Sprintf("%d_%s_%s", 0, data, base64.StdEncoding.EncodeToString([]byte(id.String())))
	t.Logf("expected - %s", expectedSignedData)
	if expectedSignedData != string(actual.SignedData) {
		t.Fatalf("expected - %s | actual - %s", expectedSignedData, actual.SignedData)
	}

	actual, err = device.Sign([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	expectedSignedData = fmt.Sprintf("%d_%s_%s", 1, data, lastSignature)
	t.Logf("expected - %s", expectedSignedData)
	if expectedSignedData != string(actual.SignedData) {
		t.Fatalf("expected - %s | actual - %s", expectedSignedData, actual.SignedData)
	}
}
