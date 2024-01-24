package filestore

import (
	"encoding/json"
	"errors"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"os"
	"sync"
)

type fileStore struct {
	Mu           sync.Mutex
	dataFilePath string
	Store        map[string]string
}

func New(dataFilePath string) (secrets.Storage, error) {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		_, err := os.Create(dataFilePath)
		if err != nil {
			return nil, err
		}
	}

	fs := &fileStore{Mu: sync.Mutex{}, Store: make(map[string]string), dataFilePath: dataFilePath}

	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *fileStore) Write(data secrets.SecretData) error {
	fs.Mu.Lock()
	defer fs.Mu.Unlock()

	fs.Store[data.Id] = data.Secret

	return fs.save()
}

func (fs *fileStore) Read(id string) (string, error) {
	fs.Mu.Lock()
	defer fs.Mu.Unlock()

	secret, exists := fs.Store[id]
	if !exists {
		return "", errors.New("value is not present")
	}

	return secret, nil
}

func (fs *fileStore) save() error {
	byteValue, err := json.Marshal(fs.Store)
	if err != nil {
		return err
	}

	return os.WriteFile(fs.dataFilePath, byteValue, 0664)
}

func (fs *fileStore) load() error {
	byteValue, err := os.ReadFile(fs.dataFilePath)
	if err != nil {
		return err
	}

	if len(byteValue) == 0 {
		fs.Store = make(map[string]string)
		return nil
	}

	if err := json.Unmarshal(byteValue, &fs.Store); err != nil {
		return err
	}

	return nil
}
