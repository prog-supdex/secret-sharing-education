package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/logger"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	MockId        = "e49c26c03203274857decfd0f4594881"
	NonExistentID = "334c4a4c42fdb79d7ebc3e73b517e6f8"
	MockSecret    = "mock_secret"
)

type mockSecretManager struct{}

func (m mockSecretManager) CreateSecret(_ string) (string, error) {
	return MockId, nil
}

func (m mockSecretManager) GetSecret(id string) (string, error) {
	if id == MockId {
		return MockSecret, nil
	}
	return "", nil
}

func TestCreateSecretHandler(t *testing.T) {
	var logBuffer bytes.Buffer
	logger.InitLogger("DEBUG", &logBuffer)

	mockManager := mockSecretManager{}
	h := handlers.NewSecretHandler(mockManager)

	router := mux.NewRouter()
	for path, handler := range h.Routes() {
		router.Handle(path, handler)
	}

	payload := handlers.RequestBodyPayload{
		PlainText: "mock_id",
	}
	payloadBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(payloadBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}

	var resp handlers.CreateSecretResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp.Id != MockId {
		t.Errorf("expected mock_id, got %v", resp.Id)
	}

	if !strings.Contains(logBuffer.String(), "Write the response") {
		t.Errorf("expected to find specific log content for CreateSecret, got: %s", logBuffer.String())
	}

	t.Run("Rejects non-POST requests", func(t *testing.T) {
		logBuffer = bytes.Buffer{}

		req, err := http.NewRequest("GET", "/", strings.NewReader(string(payloadBytes)))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if !strings.Contains(logBuffer.String(), "Invalid request method") {
			t.Errorf("expected to find invalid_request log for CreateSecret, got: %s", logBuffer.String())
		}

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected to get 422 response, got: %d", w.Code)
		}
	})
}

func TestGetSecretHandler(t *testing.T) {
	var logBuffer bytes.Buffer
	logger.InitLogger("DEBUG", &logBuffer)

	mockManager := mockSecretManager{}
	h := handlers.NewSecretHandler(mockManager)

	router := mux.NewRouter()
	for path, handler := range h.Routes() {
		router.Handle(path, handler)
	}

	req := httptest.NewRequest("GET", "/"+MockId, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %v", w.Code)
	}

	var resp handlers.GetSecretResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp.Data != MockSecret {
		t.Errorf("expected mock_secret, got %v", resp.Data)
	}

	t.Run("Rejects non-GET requests", func(t *testing.T) {
		logBuffer = bytes.Buffer{}
		logger.InitLogger("DEBUG", &logBuffer)

		req := httptest.NewRequest("POST", "/"+MockId, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if !strings.Contains(logBuffer.String(), "Invalid request method") {
			t.Errorf("expected to find invalid_request log for CreateSecret, got: %s", logBuffer.String())
		}

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected to get 422 response, got: %d", w.Code)
		}
	})

	t.Run("Rejects when Secret Data is not found", func(t *testing.T) {
		logBuffer = bytes.Buffer{}
		logger.InitLogger("DEBUG", &logBuffer)

		req := httptest.NewRequest("GET", "/"+NonExistentID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if !strings.Contains(logBuffer.String(), "The Secret Data was not found") {
			t.Errorf("expected to find specific log content for CreateSecret, got: %s", logBuffer.String())
		}

		if w.Code != http.StatusNotFound {
			t.Errorf("expected to get 404 response, got: %d", w.Code)
		}
	})
}
