package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEccSign(t *testing.T) {

	// test signer without keys
	eccSigner, err := NewECCSigner()
	assert.Nil(t, err)
	eccSigner.KeyPair = nil
	_, err = eccSigner.Sign([]byte("asd"))
	assert.NotNil(t, err)

	// test in case of zero bytes
	eccSigner, err = NewECCSigner()
	assert.Nil(t, err)
	_, err = eccSigner.Sign([]byte{})
	assert.NotNil(t, err)

	// verify the signature
	signature, err := eccSigner.Sign([]byte("asd"))
	assert.Nil(t, err)
	eccHashedData := sha256.Sum256([]byte("asd"))
	assert.True(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], signature))
}

func TestRsaSign(t *testing.T) {

	// test signer without keys
	rsaSigner, err := NewRSASigner()
	assert.Nil(t, err)
	rsaSigner.KeyPair = nil
	_, err = rsaSigner.Sign([]byte("asd"))
	assert.NotNil(t, err)

	// test in case of zero bytes
	rsaSigner, err = NewRSASigner()
	assert.Nil(t, err)
	_, err = rsaSigner.Sign([]byte{})
	assert.NotNil(t, err)

	// verify the signature
	signature, err := rsaSigner.Sign([]byte("asd"))
	assert.Nil(t, err)
	hashedData := sha256.Sum256([]byte("asd"))
	assert.Nil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], signature, nil))
}
