package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func sendJsonRequest(
	t *testing.T,
	httpMethod string,
	url string,
	serializableData ...any,
) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if len(serializableData) > 0 {
		jsonBytes, err := json.Marshal(serializableData[0])
		if err != nil {
			t.Fatal(fmt.Sprintf("json.Marshal failed: err"))
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	request, err := http.NewRequest(httpMethod, url, bodyReader)
	if err != nil {
		t.Fatal(fmt.Sprintf("json.Marshal failed: err"))
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	return response
}

func readBody(t *testing.T, response *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	return string(body)
}
