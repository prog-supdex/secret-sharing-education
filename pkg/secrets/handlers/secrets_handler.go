package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
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
	Routes() map[string]http.Handler
}

type handler struct {
	secretsManager secrets.Manager
}

func NewSecretHandler(s secrets.Manager) Handler {
	return handler{s}
}

func (h handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		slog.Debug("Invalid request method", "request method", r.Method, "expected method", "POST")
		w.WriteHeader(422)
		return
	}

	requestBody := RequestBodyPayload{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		slog.Error("Failed to decode the request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("Incoming request",
		"payload", r.Body,
		"headers", r.Header,
		"method", r.Method,
		"request_uri", r.RequestURI,
	)

	digest, err := h.secretsManager.CreateSecret(requestBody.PlainText)
	if err != nil {
		slog.Error("Failed to create secret:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateSecretResponse{Id: digest}

	jd, err := json.Marshal(&resp)
	if err != nil {
		slog.Error("Failed to marshal response:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	slog.Debug("Write the response", "id", digest)

	_, err = w.Write(jd)
	if err != nil {
		slog.Error("Failed to write response:" + err.Error())
	}
}

func (h handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		slog.Debug("Invalid request method", "request method", r.Method, "expected method", "GET")
		w.WriteHeader(422)
		return
	}

	slog.Debug("Incoming request",
		"payload", r.Body,
		"headers", r.Header,
		"method", r.Method,
		"request_uri", r.RequestURI,
	)

	vars := mux.Vars(r)
	id := vars["hash"]

	if len(id) == 0 {
		slog.Debug("The Secret ID was not present")
		http.Error(w, "No Secret ID specified", http.StatusBadRequest)
		return
	}

	decryptedSecret, err := h.secretsManager.GetSecret(id)

	if err != nil {
		slog.Error("Failed to get secret:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetSecretResponse{Data: decryptedSecret}
	jd, err := json.Marshal(&resp)
	if err != nil {
		slog.Error("Failed to marshal response:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(resp.Data) == 0 {
		slog.Debug("The Secret Data was not found")
		w.WriteHeader(http.StatusNotFound)
	}

	_, err = w.Write(jd)
	if err != nil {
		slog.Error("Failed to write response:" + err.Error())
	}
}

func (h handler) Routes() map[string]http.Handler {
	return map[string]http.Handler{
		"/":                       http.HandlerFunc(h.CreateSecret),
		"/{hash:[0-9a-fA-F]{32}}": http.HandlerFunc(h.GetSecret),
	}
}
