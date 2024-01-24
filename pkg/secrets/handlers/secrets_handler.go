package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
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
	RegisterHandler(r *mux.Router)
}

type handler struct {
	secretsManager secrets.Manager
}

func NewSecretHandler(s secrets.Manager) Handler {
	return handler{s}
}

func (h handler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/", h.CreateSecret).Methods("POST")
	r.HandleFunc("/{hash:[0-9a-fA-F]{32}}", h.GetSecret).Methods("GET")
}

func (h handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
	w.Write(jd)
}

func (h handler) GetSecret(w http.ResponseWriter, r *http.Request) {
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
	w.Write(jd)
}
