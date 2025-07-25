package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// N clients, accessing the server attempting to simulate a heavy
// load on the server. Each client will create 3 signing devices,
// sign X messages with first one, Y messages with the second and
// Z with the third
func TestConcurrency(t *testing.T) {

	srv := setupServer()
	defer closeServer(srv)
	persistence.CleanMemory() // <- bleah

	const (
		nClients = 10 // number of clients
		X        = 10 // signing operations requested for device 1
		Y        = 15 // signing operations requested for device 2
		Z        = 18 // signing operations requested for device 3
	)

	wg := sync.WaitGroup{}
	ch := make(chan bool, 5)
	cn := make(chan int, 10)
	startTime := time.Now()

	// spawn nClients goroutine, each simulating a client
	for i := 0; i < nClients; i++ {
		wg.Add(1)
		go runClient(X, Y, Z, &wg, i, cn, t, srv.listenAddress)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-time.After(30 * time.Second):
		fmt.Print("Timeout before each client could finish its stuff\n")
		return
	case <-ch:
		fmt.Print("All clients are done\n")
		break
	}

	fmt.Printf("%d Devices have been registered by %d clients\n", nClients*3, nClients)
	fmt.Printf("%d Transactions have been signed in %v\n", nClients*(X+Y+Z), time.Since(startTime))
}

// function that will represent the full cycle of a client operativity
func runClient(X int, Y int, Z int, wg *sync.WaitGroup, clientNumber int, cn chan int, t *testing.T, listenAddr string) {
	defer wg.Done()

	// client A will do 3 + X + Y + Z operations
	// so the client will sleep (3 + X + Y + Z) times
	nSleepInterval := 3 + X + Y + Z
	sleepTime := make([]time.Duration, 0)
	counter := 0
	for i := 0; i < nSleepInterval; i++ {
		intervalString := fmt.Sprint(30+rand.Intn(50)) + "ms" // sleep period from 30 to 80 milliseconds
		interval, err := time.ParseDuration(intervalString)
		assert.Nil(t, err)
		sleepTime = append(sleepTime, interval)
	}

	// create client
	client := &http.Client{}
	deviceUrl := "http://localhost" + listenAddr + "/api/v0/device"
	signUrl := "http://localhost" + listenAddr + "/api/v0/sign"

	// create first device
	body := []byte(`{
			"label":"yes",
			"signingAlgorithm":"ECC"
		}`)
	postRequest, err := http.NewRequest("POST", deviceUrl, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err := client.Do(postRequest)
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
	devUUID1 := resp.Data.Uuid

	// sleep
	time.Sleep(sleepTime[counter])
	counter++

	// create second device
	body = []byte(`{
			"signingAlgorithm":"ECC"
		}`)
	postRequest, err = http.NewRequest("POST", deviceUrl, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp = respTypeDevice{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Nil(t, uuid.Validate(resp.Data.Uuid))
	devUUID2 := resp.Data.Uuid

	// sleep
	time.Sleep(sleepTime[counter])
	counter++

	// create third device
	body = []byte(`{
			"signingAlgorithm":"RSA"
		}`)
	postRequest, err = http.NewRequest("POST", deviceUrl, bytes.NewBuffer(body))
	assert.Nil(t, err)
	res, err = client.Do(postRequest)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp = respTypeDevice{}
	resBody, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(resBody, &resp)
	assert.Nil(t, err)
	assert.Nil(t, uuid.Validate(resp.Data.Uuid))
	devUUID3 := resp.Data.Uuid

	// sleep
	time.Sleep(sleepTime[counter])
	counter++

	// signing with first device Xs times
	for i := 0; i < X; i++ {
		payload := signPayload{
			DeviceId:       devUUID1,
			DataToBeSigned: "sign-this-message-" + fmt.Sprint(i),
		}
		body, err := json.Marshal(payload)
		assert.Nil(t, err)
		postRequest, err := http.NewRequest("POST", signUrl, bytes.NewBuffer(body))
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

		// sleep
		time.Sleep(sleepTime[counter])
		counter++
	}

	// signing with second device Ys times
	for i := 0; i < Y; i++ {
		payload := signPayload{
			DeviceId:       devUUID2,
			DataToBeSigned: "sign-this-message-" + fmt.Sprint(i),
		}
		body, err := json.Marshal(payload)
		assert.Nil(t, err)
		postRequest, err := http.NewRequest("POST", signUrl, bytes.NewBuffer(body))
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

		// sleep
		time.Sleep(sleepTime[counter])
		counter++
	}

	// signing with third device Zs times
	for i := 0; i < Z; i++ {
		payload := signPayload{
			DeviceId:       devUUID3,
			DataToBeSigned: "sign-this-message-" + fmt.Sprint(i),
		}
		body, err := json.Marshal(payload)
		assert.Nil(t, err)
		postRequest, err := http.NewRequest("POST", signUrl, bytes.NewBuffer(body))
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

		// sleep
		time.Sleep(sleepTime[counter])
		counter++
	}

	cn <- clientNumber
}
