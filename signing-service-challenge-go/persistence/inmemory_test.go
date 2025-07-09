package persistence

import (
	"encoding/base64"
	"testing"

	crypto "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetDevices(t *testing.T) {

	// add devices
	label := "dev1"
	dev1, err := AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)
	label = "dev2"
	dev2, err := AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)
	label = "dev3"
	dev3, err := AddDevice(crypto.RSA, &label)
	assert.Nil(t, err)

	// test that we can get them
	devs, err := GetDevices()
	assert.Nil(t, err)
	assert.Len(t, devs, 3)

	// make sure all three devices are listed
	counter := 0
	for _, dev := range devs {
		if dev.UUID == dev1.UUID || dev.UUID == dev2.UUID || dev.UUID == dev3.UUID {
			counter++
		}
	}
	assert.Equal(t, 3, counter)
}

func TestGetDevice(t *testing.T) {

	// add device
	label := "dev1"
	dev1, err := AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)

	// test unexisting uuid
	_, err = GetDevice("asd")
	assert.NotNil(t, err)

	// test working case
	dev, err := GetDevice(dev1.UUID)
	assert.Nil(t, err)
	assert.Equal(t, dev1.Label, dev.Label)
	assert.Equal(t, dev1.UUID, dev.UUID)
}

func TestAddDevice(t *testing.T) {

	// test with unexisting signing method
	var RUSSO crypto.SigningAlgorithm = 3
	label := "label"
	_, err := AddDevice(RUSSO, &label)
	assert.NotNil(t, err)

	// test workign case for ECC
	devECC, err := AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)
	assert.Equal(t, devECC.Label, label)
	assert.Nil(t, uuid.Validate(devECC.UUID))
	assert.Equal(t, devECC.LastSignatureBase64EncodedString, base64.StdEncoding.EncodeToString([]byte(string(devECC.UUID))))
	assert.Equal(t, uint(0), devECC.SignatureCounter)
	eccSigner, ok := devECC.Signer.(*crypto.ECCSigner)
	assert.True(t, ok)
	assert.NotNil(t, eccSigner.KeyGenerator)
	assert.NotNil(t, eccSigner.KeyPair)
	assert.NotNil(t, eccSigner.Marshaler)

	// test workign case for RSA
	devRSA, err := AddDevice(crypto.RSA, &label)
	assert.Nil(t, err)
	assert.Equal(t, devRSA.Label, label)
	assert.Nil(t, uuid.Validate(devRSA.UUID))
	assert.Equal(t, devRSA.LastSignatureBase64EncodedString, base64.StdEncoding.EncodeToString([]byte(string(devRSA.UUID))))
	assert.Equal(t, uint(0), devRSA.SignatureCounter)
	rsaSigner, ok := devRSA.Signer.(*crypto.RSASigner)
	assert.True(t, ok)
	assert.NotNil(t, rsaSigner.KeyGenerator)
	assert.NotNil(t, rsaSigner.KeyPair)
	assert.NotNil(t, rsaSigner.Marshaler)

	// test the label in case no one was provided and in case it was provided
	devWithLabel, err := AddDevice(crypto.RSA, &label)
	assert.Nil(t, err)
	assert.Equal(t, devWithLabel.Label, label)
	devWithoutLabel, err := AddDevice(crypto.RSA, nil)
	assert.Nil(t, err)
	assert.Equal(t, devWithoutLabel.Label, devWithoutLabel.UUID)
}

// Most of the verification related to signing has been done on
// the the tests of the API -> api/sign_test.go
func TestSignTransaction(t *testing.T) {

	// test uuid not found
	_, _, err := SignTransaction("asd", "data")
	assert.ErrorContains(t, err, "device with UUID 'asd' not found")

	// add ECC device
	eccDev, err := AddDevice(crypto.ECC, nil)
	assert.Nil(t, err)

	// sign three times and check that the counter increases
	_, _, err = SignTransaction(eccDev.UUID, "data")
	assert.Nil(t, err)
	_, _, err = SignTransaction(eccDev.UUID, "data")
	assert.Nil(t, err)
	_, _, err = SignTransaction(eccDev.UUID, "data")
	assert.Nil(t, err)
	retrievedDev, err := GetDevice(eccDev.UUID)
	assert.Nil(t, err)
	assert.Equal(t, retrievedDev.SignatureCounter, uint(3))

	// add RSA device
	rsaDev, err := AddDevice(crypto.RSA, nil)
	assert.Nil(t, err)

	// sign three times and check that the counter increases
	_, _, err = SignTransaction(rsaDev.UUID, "data")
	assert.Nil(t, err)
	_, _, err = SignTransaction(rsaDev.UUID, "data")
	assert.Nil(t, err)
	_, _, err = SignTransaction(rsaDev.UUID, "data")
	assert.Nil(t, err)
	retrievedDev, err = GetDevice(eccDev.UUID)
	assert.Nil(t, err)
	assert.Equal(t, retrievedDev.SignatureCounter, uint(3))

	// correct storage and upgrade of LastSignatureBase64EncodedString
	// has been tested on the API side -> api/sign_test.go
}
