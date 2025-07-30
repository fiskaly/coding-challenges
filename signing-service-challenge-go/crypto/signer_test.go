package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewECCSigner(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)

	signer := NewECCSigner(privateKey)

	assert.NotNil(t, signer)
	assert.Equal(t, privateKey, signer.privateKey)
}

func TestECCSignerWithSignSuccess(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	signer := NewECCSigner(privateKey)

	dataToBeSigned := []byte("sign-test-data")

	signedData, err := signer.Sign(dataToBeSigned)

	assert.NoError(t, err)
	assert.NotNil(t, signedData)
	assert.Greater(t, len(signedData), 0, "sign should not be empty")
}

func TestNewRSASigner(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)

	signer := NewRSASigner(privateKey)

	assert.NotNil(t, signer)
	assert.Equal(t, privateKey, signer.privateKey)
}

func TestRSASignerWithSignSuccess(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	assert.NoError(t, err)

	signer := NewRSASigner(privateKey)

	dataToBeSigned := []byte("sign-test-data")

	signedData, err := signer.Sign(dataToBeSigned)

	assert.NoError(t, err, "expect no error when signing data")
	assert.NotNil(t, signedData, "expect signed data to not be nil")
	assert.Greater(t, len(signedData), 0, "expect signed data to not be empty")
}
