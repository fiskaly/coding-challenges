package domain_test

import (
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
)

func TestSignTransaction(t *testing.T) {
	t.Run("successfully signs when device is being used for the first time", func(t *testing.T) {
		dataToBeSigned := "some-transaction-data"
		repository := persistence.NewInMemorySignatureDeviceRepository()
		deviceId := uuid.MustParse("121fe402-762a-411a-8eeb-9e6c3ca16886")
		device, err := domain.BuildSignatureDevice(deviceId, crypto.RSAGenerator{})
		if err != nil {
			t.Fatal(err)
		}
		err = repository.Create(device)
		if err != nil {
			t.Fatal(err)
		}

		encodedSignature, signedData, err := domain.SignTransaction(device, repository, dataToBeSigned)
		if err != nil {
			t.Fatal(err)
		}

		base64EncodedDeviceId := "MTIxZmU0MDItNzYyYS00MTFhLThlZWItOWU2YzNjYTE2ODg2"
		expectedSignedData := fmt.Sprintf("0_%s_%s", dataToBeSigned, base64EncodedDeviceId)
		if signedData != expectedSignedData {
			t.Errorf("expected signedData: %s, got: %s", expectedSignedData, signedData)
		}

		decodedSignature, err := base64.StdEncoding.DecodeString(encodedSignature)
		if err != nil {
			t.Errorf("failed to base64 decode the signature: %s", err)
		}

		// verify the decoded signature with the public key
		hash := sha512.New384()
		_, err = hash.Write([]byte(signedData))
		if err != nil {
			t.Fatal(err)
		}
		digest := hash.Sum(nil)

		rsaKeyPair := device.KeyPair.(*crypto.RSAKeyPair)
		err = rsa.VerifyPSS(
			rsaKeyPair.Public,
			crypto.HashFunction,
			digest,
			decodedSignature,
			nil,
		)
		if err != nil {
			t.Errorf("signature is not valid. err: %s", err)
		}

		// check updates to signature device
		// refetch the device from the repository to reflect updates
		device, ok, err := repository.Find(deviceId)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("device not found")
		}
		if device.SignatureCounter != 1 {
			t.Errorf("expected signature counter to be incremented to 1, got %d", device.SignatureCounter)
		}
		if device.Base64EncodedLastSignature != encodedSignature {
			t.Errorf("expected last signature to be updated to: %s, got %s", encodedSignature, device.Base64EncodedLastSignature)
		}
	})
}

func TestSecureDataToBeSigned(t *testing.T) {
	t.Run("concatenates data with counter and last signature when counter > 0", func(t *testing.T) {
		lastSignature := "last-signature"
		base64EncodedLastSignature := "bGFzdC1zaWduYXR1cmU="

		device := domain.SignatureDevice{
			Base64EncodedLastSignature: lastSignature,
			SignatureCounter:           1,
		}
		data := "some transaction data"

		got := domain.SecureDataToBeSigned(device, data)
		expected := fmt.Sprintf("1_%s_%s", data, base64EncodedLastSignature)

		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	})

	t.Run("concatenates data with counter and device id when counter == 0", func(t *testing.T) {
		id := uuid.MustParse("ed40597c-52b7-40bc-9e15-83e4741a102b")
		base64EncodedID := "ZWQ0MDU5N2MtNTJiNy00MGJjLTllMTUtODNlNDc0MWExMDJi"

		device := domain.SignatureDevice{
			ID:                         id,
			Base64EncodedLastSignature: "",
			SignatureCounter:           0,
		}
		data := "some transaction data"

		got := domain.SecureDataToBeSigned(device, data)
		expected := fmt.Sprintf("0_%s_%s", data, base64EncodedID)

		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	})
}
