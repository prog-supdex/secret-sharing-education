package filestore

import (
	"encoding/json"
	"errors"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"log/slog"
	"os"
	"sync"
)

type FileStore struct {
	mu           sync.Mutex
	dataFilePath string
	Store        map[string]string
}

func New(dataFilePath string) (*FileStore, error) {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		slog.Debug("Creating filestore", "file_path", dataFilePath)
		_, err := os.Create(dataFilePath)
		if err != nil {
			slog.Error("Error creating file: "+err.Error(), "file_path", dataFilePath)
			return nil, err
		}
		slog.Debug("The file was created", "file_path", dataFilePath)
	}

	fs := &FileStore{mu: sync.Mutex{}, Store: make(map[string]string), dataFilePath: dataFilePath}

	if err := fs.load(); err != nil {
		slog.Error("Error loading file to FileStore: " + err.Error())
		return nil, err
	}

	return fs, nil
}

func (fs *FileStore) Write(data secrets.SecretData) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.Store[data.Id] = data.Secret

	slog.Debug("Writing the secret with ID:" + data.Id)

	return fs.save()
}

func (fs *FileStore) Read(id string) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	secret, exists := fs.Store[id]
	if !exists {
		slog.Debug("The id: " + id + " was not found")
		return "", errors.New("value is not present")
	}

	return secret, nil
}

func (fs *FileStore) save() error {
	byteValue, err := json.Marshal(fs.Store)
	if err != nil {
		slog.Error("The marshaling content was failed: " + err.Error())
		return err
	}

	return os.WriteFile(fs.dataFilePath, byteValue, 0664)
}

func (fs *FileStore) load() error {
	byteValue, err := os.ReadFile(fs.dataFilePath)
	if err != nil {
		slog.Error("The reading file was failed: "+err.Error(), "file_path", fs.dataFilePath)
		return err
	}

	if len(byteValue) == 0 {
		fs.Store = make(map[string]string)
		return nil
	}

	if err := json.Unmarshal(byteValue, &fs.Store); err != nil {
		slog.Error("The unmarshaling was failed: " + err.Error())
		return err
	}

	return nil
}
