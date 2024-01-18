package handlers

import (
	"encoding/json"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"status": "OK",
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
