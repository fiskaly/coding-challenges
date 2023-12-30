package service

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/domain"
)

type testingInput struct {
	Name  string
	Input domain.SignatureDevice
	Want  string
}

var text = "data_to_be_signed"

func TestBuildSigningString(t *testing.T) {
	var testingInputs [2]testingInput

	testDevice1 := domain.SignatureDevice{
		Id:               "test_noSignatures",
		SignatureCounter: 0,
		LastSignature:    nil,
	}

	testDeivce1IdBase64 := base64.StdEncoding.EncodeToString([]byte(testDevice1.Id))

	testingInputs[0] = testingInput{
		Name:  "when no past signatures present, use base64 encoded device id",
		Input: testDevice1,
		Want:  fmt.Sprintf("%d_%s_%s", 0, text, testDeivce1IdBase64),
	}

	testDevice2 := domain.SignatureDevice{
		Id:               "test_noSignatures",
		SignatureCounter: 1,
		LastSignature:    []byte("sadasdsadasdasdasdsadsadsadasdsdasdsadsdadssdasdasdasdsdads"),
	}

	testingInputs[1] = testingInput{
		Name:  "when past signatures present, use last signature",
		Input: testDevice2,
		Want:  fmt.Sprintf("%d_%s_%s", 1, text, string(testDevice2.LastSignature)),
	}

	for _, tt := range testingInputs {
		t.Run(tt.Name, func(t *testing.T) {
			ans := buildSigningString(tt.Input, text)
			if ans != tt.Want {
				t.Errorf("got %s, want %s", ans, tt.Want)
			}
		})
	}
}
