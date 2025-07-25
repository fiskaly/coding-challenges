package api

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	// ---------------------
	// test incorrect method
	// ---------------------

	url := "http://localhost" + srv.listenAddress + "/api/v0/health"

	// create requests
	postRequest, err := http.NewRequest("POST", url, nil)
	assert.Nil(t, err)
	patchRequest, err := http.NewRequest("PATCH", url, nil)
	assert.Nil(t, err)
	deleteRequest, err := http.NewRequest("DELETE", url, nil)
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

	// make good request
	getRequest, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	res, err = client.Do(getRequest)
	assert.Nil(t, err)
	resp := respTypeHealth{}
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Equal(t, "pass", resp.Data.Status)
	assert.Equal(t, "v0", resp.Data.Version)
}
