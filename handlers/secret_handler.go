package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/prog-supdex/mini-project/milestone-code/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/types"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createSecret(w, r)
	case "GET":
		getSecret(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	p := types.CreateSecretPayload{}
	err = json.Unmarshal(bytes, &p)

	if err != nil {
		log.Fatal(err)
	}

	digest := getHash(p.PlainText)
	resp := types.CreateSecretResponse{Id: digest}

	s := types.SecretData{Id: resp.Id, Secret: p.PlainText}
	err = filestore.FileStoreConfig.Fs.Write(s)
	// todo error handling
	if err != nil {
		log.Fatal(err)
	}

	jd, err := json.Marshal(&resp)

	if err != nil {
		log.Printf("%v\n%s", err, debug.Stack())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jd)
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	id = strings.TrimPrefix(id, "/")
	if len(id) == 0 {
		http.Error(w, "No Secret ID specified", http.StatusBadRequest)
		return
	}
	resp := types.GetSecretResponse{}
	v, err := filestore.FileStoreConfig.Fs.Read(id)
	if err != nil {
		log.Fatal(err)
	}

	resp.Data = v
	jd, err := json.Marshal(&resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(resp.Data) == 0 {
		w.WriteHeader(404)
	}
	w.Write(jd)
}

func getHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
