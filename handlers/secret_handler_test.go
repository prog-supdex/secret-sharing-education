package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prog-supdex/mini-project/milestone-code/filestore"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	TempFilePath = "../testdata/temp_data.json"
	SecretId     = "d0b2fa28502506a01bb074879eaa083d"
)

func TestCreateSecret(t *testing.T) {
	defer func() {
		t.Cleanup(func() {
			os.Remove(TempFilePath)
		})
	}()

	filestore.Init(TempFilePath)

	testServer := httptest.NewServer(http.HandlerFunc(secretHandler))
	defer testServer.Close()

	jsonData := map[string]string{"plain_text": "my_secret"}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
		return
	}

	req, err := http.NewRequest("POST", testServer.URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Errorf("Request error: %v", err)
		return
	}

	err = validateResponse(resp)
	if err != nil {
		t.Errorf("Validation failed: %s", err)
		return
	}

	result, err := unmarshalResponse(resp)
	if err != nil {
		t.Errorf("Response error: %s", err)
		return
	}

	if result["id"] != getHash("my_secret") {
		t.Errorf("invalid response: %s", result)
	}
}

func TestGetSecret(t *testing.T) {
	filestore.Init("../testdata/data.json")

	testServer := httptest.NewServer(http.HandlerFunc(secretHandler))
	defer testServer.Close()

	resp, err := testServer.Client().Get(testServer.URL + "/" + SecretId)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	err = validateResponse(resp)
	if err != nil {
		t.Errorf("Validation failed: %s", err)
		return
	}

	result, err := unmarshalResponse(resp)
	if err != nil {
		t.Errorf("Response error: %s", err)
		return
	}

	if result["data"] != "my_secret_value" {
		t.Errorf("invalid response: %s", result)
	}
}

func validateResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is not 200: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("content type is not json: %s", resp.Header.Get("Content-Type"))
	}

	return nil
}

func unmarshalResponse(response *http.Response) (result map[string]string, err error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response JSON: %v", err)
	}

	return result, nil
}
