package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"log/slog"
	"net/http"
)

type RequestBodyPayload struct {
	PlainText string `json:"plain_text"`
}

type CreateSecretResponse struct {
	Id string `json:"id"`
}

type GetSecretResponse struct {
	Data string `json:"data"`
}

type Handler interface {
	CreateSecret(w http.ResponseWriter, r *http.Request)
	GetSecret(w http.ResponseWriter, r *http.Request)
	Routes() map[string]http.HandlerFunc
}

type handler struct {
	secretsManager secrets.Manager
}

func NewSecretHandler(s secrets.Manager) Handler {
	return handler{s}
}

func (h handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "GET" {
		server.RenderTemplate(w, "home", nil)

		return
	}

	secretText, err := extractSecretFromRequest(r)
	if err != nil {
		slog.Error("Error request: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	slog.Debug("Incoming request",
		"payload", r.Body,
		"headers", r.Header,
		"method", r.Method,
		"request_uri", r.RequestURI,
	)

	digest, err := h.secretsManager.CreateSecret(secretText)
	if err != nil {
		slog.Error("Failed to create secret:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("Write the response", "id", digest)

	if r.Header.Get("Content-Type") == "application/json" {
		resp := CreateSecretResponse{Id: digest}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error("Failed to write response:" + err.Error())
		}
	} else {
		data := map[string]interface{}{
			"SecretID": digest,
		}
		server.RenderTemplate(w, "home", data)
	}
}

func (h handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		slog.Debug("Invalid request method", "request method", r.Method, "expected method", "GET")
		w.WriteHeader(422)
		return
	}

	vars := mux.Vars(r)
	id := vars["hash"]

	if len(id) == 0 {
		slog.Debug("The Secret ID was not present")
		http.Error(w, "No Secret ID specified", http.StatusBadRequest)
		return
	}

	slog.Debug("Incoming request",
		"payload", r.Body,
		"headers", r.Header,
		"method", r.Method,
		"request_uri", r.RequestURI,
	)

	decryptedSecret, err := h.secretsManager.GetSecret(id)
	if err != nil {
		slog.Error("Failed to get secret:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(decryptedSecret) == 0 {
		slog.Debug("The Secret Data was not found")
		w.WriteHeader(http.StatusNotFound)
	}

	if r.Header.Get("Accept") == "application/json" {
		resp := GetSecretResponse{Data: decryptedSecret}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error("Failed to write response:" + err.Error())
		}
	} else {
		data := map[string]interface{}{
			"SecretData": decryptedSecret,
		}
		server.RenderTemplate(w, "secret", data)
	}
}

func (h handler) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/":                               h.CreateSecret,
		"/{hash:[0-9a-fA-F]{32}}":         h.GetSecret,
		"/secrets/{hash:[0-9a-fA-F]{32}}": h.GetSecret,
	}
}

func extractSecretFromRequest(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	switch {
	case contentType == "application/json":
		var requestBody RequestBodyPayload
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			return "", fmt.Errorf("failed to decode JSON body: %w", err)
		}
		return requestBody.PlainText, nil
	case contentType == "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return "", fmt.Errorf("failed to parse form data: %w", err)
		}
		return r.FormValue("secret"), nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}
