package api

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	crypto_fiscaly "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// ---------------------
	// test incorrect method
	// ---------------------

	url := "http://localhost" + srv.listenAddress + "/api/v0/sign"

	// make requests
	getRequest, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	patchRequest, err := http.NewRequest("PATCH", url, nil)
	assert.Nil(t, err)
	deleteRequest, err := http.NewRequest("DELETE", url, nil)
	assert.Nil(t, err)

	// create client
	client := &http.Client{}

	// make requests
	res, err := client.Do(getRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	res, err = client.Do(patchRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	res, err = client.Do(deleteRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

	// ---------------------
	// test no payload
	// ---------------------

	postRequest, err := http.NewRequest("POST", url, nil)
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// ---------------------
	// test empty payload
	// ---------------------

	body := []byte(``)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// ---------------------
	// test unmarhsable payload
	// ---------------------

	body = []byte(`{}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// --------------------
	// test invalid payload
	// --------------------

	// invalid device id
	body = []byte(`{
		"deviceId":2,
		"dataToBeSigned": "asd"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// invalid data to be signed
	body = []byte(`{
		"deviceId": "asd"
		"dataToBeSigned":2,
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// no device id
	body = []byte(`{
		"dataToBeSigned": "asd"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// no data to be signed
	body = []byte(`{
		"deviceId":"asd"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

}

func TestSignTransactionRSA(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// --------------------
	// device 1 (RSA)
	// --------------------

	// --------------------
	// test working example
	// --------------------

	url := "http://localhost" + srv.listenAddress + "/api/v0/sign"

	// add device
	label := "dev1"
	dev, err := persistence.AddDevice(crypto_fiscaly.RSA, &label)
	assert.Nil(t, err)

	// create client
	client := &http.Client{}

	// sign 1st time
	payload := signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-please",
	}
	body, err := json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err := client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resBody []byte
	resp1 := respTypeSign{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp1)
	assert.Nil(t, err)

	// sign 2nd time
	payload = signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-again-please",
	}
	body, err = json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp2 := respTypeSign{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp2)
	assert.Nil(t, err)

	// sign 3rd time
	payload = signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-last-time-please",
	}
	body, err = json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp3 := respTypeSign{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp3)
	assert.Nil(t, err)

	// check the counter returned in the response
	assert.True(t, strings.HasPrefix(resp1.Data.SignedData, "0"))
	assert.True(t, strings.HasPrefix(resp2.Data.SignedData, "1"))
	assert.True(t, strings.HasPrefix(resp3.Data.SignedData, "2"))

	// check that the last signature is correct
	lastSignature := make([]string, 3)
	for i, resp := range []signatureOut{resp1.Data, resp2.Data, resp3.Data} {
		signedData := strings.Split(resp.SignedData, "_")
		assert.Len(t, signedData, 3)
		lastSignature[i] = signedData[2]
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dev.UUID))
	assert.Equal(t, string(encoded), lastSignature[0])
	assert.Equal(t, resp1.Data.Signature, lastSignature[1])
	assert.Equal(t, resp2.Data.Signature, lastSignature[2])

	// check the siganture can be verified
	// resp1
	decoded, err := base64.StdEncoding.DecodeString(resp1.Data.Signature)
	assert.Nil(t, err)
	hashedData := sha256.Sum256([]byte(resp1.Data.SignedData))
	rsaSigner, ok := dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.Nil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], decoded, nil))
	// resp2
	decoded, err = base64.StdEncoding.DecodeString(resp2.Data.Signature)
	assert.Nil(t, err)
	hashedData = sha256.Sum256([]byte(resp2.Data.SignedData))
	rsaSigner, ok = dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.Nil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], decoded, nil))
	// resp3
	decoded, err = base64.StdEncoding.DecodeString(resp3.Data.Signature)
	assert.Nil(t, err)
	hashedData = sha256.Sum256([]byte(resp3.Data.SignedData))
	rsaSigner, ok = dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.Nil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], decoded, nil))

	// check the case in which the data gets modified that the signature is corrupted
	// resp1
	hashedData = sha256.Sum256([]byte("ALTERED"))
	rsaSigner, ok = dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.NotNil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], []byte(resp1.Data.Signature), nil))
	// resp2
	hashedData = sha256.Sum256([]byte("ALTERED"))
	rsaSigner, ok = dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.NotNil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], []byte(resp1.Data.Signature), nil))
	// resp3
	hashedData = sha256.Sum256([]byte("ALTERED"))
	rsaSigner, ok = dev.Signer.(*crypto_fiscaly.RSASigner)
	assert.True(t, ok)
	assert.NotNil(t, rsa.VerifyPSS(rsaSigner.KeyPair.Public, crypto.SHA256, hashedData[:], []byte(resp1.Data.Signature), nil))

}

func TestSignTransactionECC(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// --------------------
	// device 2 (ECC)
	// --------------------

	// --------------------
	// test working example
	// --------------------

	url := "http://localhost" + srv.listenAddress + "/api/v0/sign"

	// add device
	label := "dev2"
	dev, err := persistence.AddDevice(crypto_fiscaly.ECC, &label)
	assert.Nil(t, err)

	// create client
	client := &http.Client{}

	// sign 1st time
	payload := signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-please",
	}
	body, err := json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err := client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp1 := respTypeSign{}
	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp1)
	assert.Nil(t, err)

	// sign 2nd time
	payload = signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-again-please",
	}
	body, err = json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp2 := respTypeSign{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp2)
	assert.Nil(t, err)

	// sign 3rd time
	payload = signPayload{
		DeviceId:       dev.UUID,
		DataToBeSigned: "sign-me-last-time-please",
	}
	body, err = json.Marshal(payload)
	assert.Nil(t, err)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp3 := respTypeSign{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp3)
	assert.Nil(t, err)

	// check the counter returned in the response
	assert.True(t, strings.HasPrefix(resp1.Data.SignedData, "0"))
	assert.True(t, strings.HasPrefix(resp2.Data.SignedData, "1"))
	assert.True(t, strings.HasPrefix(resp3.Data.SignedData, "2"))

	// check that the last signature is correct
	lastSignature := make([]string, 3)
	for i, resp := range []signatureOut{resp1.Data, resp2.Data, resp3.Data} {
		signedData := strings.Split(resp.SignedData, "_")
		assert.Len(t, signedData, 3)
		lastSignature[i] = signedData[2]
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dev.UUID))
	assert.Nil(t, err)
	assert.Equal(t, string(encoded), lastSignature[0])
	assert.Equal(t, resp1.Data.Signature, lastSignature[1])
	assert.Equal(t, resp2.Data.Signature, lastSignature[2])

	// verify the signature
	// resp1
	decoded, err := base64.StdEncoding.DecodeString(resp1.Data.Signature)
	assert.Nil(t, err)
	eccHashedData := sha256.Sum256([]byte(resp1.Data.SignedData))
	eccSigner, ok := dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.True(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))
	// resp2
	decoded, err = base64.StdEncoding.DecodeString(resp2.Data.Signature)
	assert.Nil(t, err)
	eccHashedData = sha256.Sum256([]byte(resp2.Data.SignedData))
	eccSigner, ok = dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.True(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))
	// resp3
	decoded, err = base64.StdEncoding.DecodeString(resp3.Data.Signature)
	assert.Nil(t, err)
	eccHashedData = sha256.Sum256([]byte(resp3.Data.SignedData))
	eccSigner, ok = dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.True(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))

	// check the case in which the data gets modified
	// resp1
	decoded, err = base64.StdEncoding.DecodeString(resp1.Data.Signature)
	assert.Nil(t, err)
	eccHashedData = sha256.Sum256([]byte("ALTERED"))
	eccSigner, ok = dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.False(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))
	// resp2
	decoded, err = base64.StdEncoding.DecodeString(resp1.Data.Signature)
	assert.Nil(t, err)
	eccHashedData = sha256.Sum256([]byte("ALTERED"))
	eccSigner, ok = dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.False(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))
	// resp3
	decoded, err = base64.StdEncoding.DecodeString(resp1.Data.Signature)
	assert.Nil(t, err)
	eccHashedData = sha256.Sum256([]byte("ALTERED"))
	eccSigner, ok = dev.Signer.(*crypto_fiscaly.ECCSigner)
	assert.True(t, ok)
	assert.False(t, ecdsa.VerifyASN1(eccSigner.KeyPair.Public, eccHashedData[:], decoded))
}
