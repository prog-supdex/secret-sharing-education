package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	MockId     = "e49c26c03203274857decfd0f4594881"
	MockSecret = "mock_secret"
)

type mockSecretManager struct{}

func (m mockSecretManager) CreateSecret(_ string) (string, error) {
	return MockId, nil
}

func (m mockSecretManager) GetSecret(id string) (string, error) {
	if id == MockId {
		return MockSecret, nil
	}
	return "", errors.New("value is not present")
}

func TestCreateSecretHandler(t *testing.T) {
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
}

func TestGetSecretHandler(t *testing.T) {
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
}
