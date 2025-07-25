package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// ---------------------
	// test incorrect method
	// ---------------------

	url := "http://localhost" + srv.listenAddress + "/api/v0/device"

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
	// test invalid payload
	// ---------------------

	// invalid label
	body = []byte(`{
		"label":2,
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// no signingAlgorithm
	body = []byte(`{
		"label":"asd",
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// invalid signingAlgorithm_1
	body = []byte(`{
		"signingAlgorithm":2,
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// invalid signingAlgorithm_2
	body = []byte(`{
		"signingAlgorithm":"2",
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// invalid signingAlgorithm_3
	body = []byte(`{
		"signingAlgorithm":"NOPE",
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// ----------------------------------
	// working example and check response
	// ----------------------------------

	// working example 1
	body = []byte(`{
		"label":"yes",
		"signingAlgorithm":"ECC"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resBody []byte
	resp := respTypeDevice{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, "yes", resp.Data.Label)
	assert.Nil(t, uuid.Validate(resp.Data.Uuid))

	// working example 2
	resp = respTypeDevice{}
	body = []byte(`{
		"label":"asd",
		"signingAlgorithm":"RSA"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, "asd", resp.Data.Label)
	assert.Nil(t, uuid.Validate(resp.Data.Uuid))

	// working example 3
	resp = respTypeDevice{}
	body = []byte(`{
		"signingAlgorithm":"ECC"
	}`)
	postRequest, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, resp.Data.Uuid, resp.Data.Label)
	assert.Nil(t, uuid.Validate(resp.Data.Uuid))

}

func TestGetDevice(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// ---------------------
	// test incorrect method
	// ---------------------

	// add a device
	label := "asd"
	dev, err := persistence.AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)

	url := "http://localhost" + srv.listenAddress + "/api/v0/device/"

	// make requests
	postRequest, err := http.NewRequest("POST", url+dev.UUID, nil)
	assert.Nil(t, err)
	patchRequest, err := http.NewRequest("PATCH", url+dev.UUID, nil)
	assert.Nil(t, err)
	deleteRequest, err := http.NewRequest("DELETE", url+dev.UUID, nil)
	assert.Nil(t, err)

	// create client
	client := &http.Client{}

	// make requests
	res, err := client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	res, err = client.Do(patchRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	res, err = client.Do(deleteRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

	// ----------------------
	// test not existing uuid
	// ----------------------

	reqNoUUID, err := http.NewRequest("GET", url+"asd", nil)
	assert.Nil(t, err)
	res, err = client.Do(reqNoUUID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	reqNoUUID, err = http.NewRequest("GET", url+"2", nil)
	assert.Nil(t, err)
	res, err = client.Do(reqNoUUID)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// -------------------------------------
	// test valid exmaple and check response
	// -------------------------------------

	goodReq, err := http.NewRequest("GET", url+dev.UUID, nil)
	assert.Nil(t, err)
	res, err = client.Do(goodReq)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp := respTypeDevice{}
	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, dev.Label, resp.Data.Label)
	assert.Equal(t, dev.UUID, resp.Data.Uuid)
}

func TestGetDevices(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// add a couple of devices
	label := "dev1"
	dev1, err := persistence.AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)
	label = "dev2"
	dev2, err := persistence.AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)
	label = "dev1"
	dev3, err := persistence.AddDevice(crypto.ECC, &label)
	assert.Nil(t, err)

	// create client
	client := &http.Client{}

	// make request and check reqsponse
	url := "http://localhost" + srv.listenAddress + "/api/v0/devices"
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp := respTypeDeviceArr{}
	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Len(t, resp.Data, 3)

	// make sure all three devices are listed
	counter := 0
	for _, dev := range resp.Data {
		if dev.Uuid == dev1.UUID || dev.Uuid == dev2.UUID || dev.Uuid == dev3.UUID {
			counter++
		}
	}
	assert.Equal(t, 3, counter)

}
