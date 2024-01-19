package filestore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/prog-supdex/mini-project/milestone-code/types"
)

type fileStore struct {
	Mu    sync.Mutex
	Store map[string]string
}

var FileStoreConfig struct {
	DataFilePath string
	Fs           fileStore
}

func Init(dataFilePath string) error {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		_, err := os.Create(dataFilePath)
		if err != nil {
			return err
		}
	}
	FileStoreConfig.Fs = fileStore{Mu: sync.Mutex{}, Store: make(map[string]string)}
	FileStoreConfig.DataFilePath = dataFilePath
	return nil
}

func (j *fileStore) Write(data types.SecretData) error {
	j.Mu.Lock()
	defer j.Mu.Unlock()

	fileContent, err := readJson(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}

	fileContent[data.Id] = data.Secret
	byteValue, err := json.Marshal(fileContent)

	if err != nil {
		return err
	}

	err = os.WriteFile(FileStoreConfig.DataFilePath, byteValue, 0664)
	if err != nil {
		return err
	}

	return nil
}

func (j *fileStore) Read(id string) (string, error) {
	j.Mu.Lock()
	defer j.Mu.Unlock()

	fileContent, err := readJson(FileStoreConfig.DataFilePath)

	if err != nil {
		log.Fatal(err)
	}

	decryptedValue, exists := fileContent[id]
	fmt.Println(fileContent, id)
	if !exists {
		return "", errors.New("value is not present")
	}

	return decryptedValue, nil
}

func readJson(path string) (map[string]string, error) {
	byteValue, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if len(byteValue) == 0 {
		return make(map[string]string), nil
	}

	var data map[string]string

	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
