package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"log"
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
		w.WriteHeader(422)
		return
	}

	requestBody := RequestBodyPayload{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	digest, err := h.secretsManager.CreateSecret(requestBody.PlainText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateSecretResponse{Id: digest}

	jd, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jd)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func (h handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(422)
		return
	}

	vars := mux.Vars(r)
	id := vars["hash"]
	if len(id) == 0 {
		http.Error(w, "No Secret ID specified", http.StatusBadRequest)
		return
	}

	decryptedSecret, err := h.secretsManager.GetSecret(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetSecretResponse{Data: decryptedSecret}
	jd, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(resp.Data) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	_, err = w.Write(jd)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func (h handler) Routes() map[string]http.Handler {
	return map[string]http.Handler{
		"/":                       http.HandlerFunc(h.CreateSecret),
		"/{hash:[0-9a-fA-F]{32}}": http.HandlerFunc(h.GetSecret),
	}
}
