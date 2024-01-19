package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(healthCheckHandler))
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL)
	if err != nil {
		t.Errorf("Get error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response code is not 200: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content type is not json: %s", resp.Header.Get("Content-Type"))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll error: %v", err)
	}

	if string(data) != "{\"status\":\"OK\"}" {
		t.Error("response body does not equal to Hello World")
	}
}
