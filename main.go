package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/prog-supdex/mini-project/milestone-code/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/handlers"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	port := os.Getenv("LISTEN_ADDR")

	if port == "" {
		port = ":8080"
	}

	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	dataFilePath := os.Getenv("DATA_FILE_PATH")

	if dataFilePath == "" {
		log.Fatal(errors.New("filepath is not present"))
	}

	err = filestore.Init(dataFilePath)
	err = http.ListenAndServe(port, mux)
}
