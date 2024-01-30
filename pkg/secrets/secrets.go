package secrets

import (
	"crypto/md5"
	"encoding/hex"
	"log/slog"
)

type SecretData struct {
	Id     string
	Secret string
}

type Storage interface {
	Read(id string) (string, error)
	Write(data SecretData) error
}

type Manager interface {
	CreateSecret(plainText string) (string, error)
	GetSecret(id string) (string, error)
}
type manager struct {
	storage Storage
}

type GetSecretResponse struct {
	Data string `json:"data"`
}

func NewSecretManager(s Storage) Manager {
	return manager{s}
}

func (m manager) CreateSecret(plainText string) (string, error) {
	digest := getHash(plainText)
	data := SecretData{Id: digest, Secret: plainText}

	if err := m.storage.Write(data); err != nil {
		slog.Error("Failed to write data to the storage: " + err.Error())
		return "", err
	}

	return data.Id, nil
}

func (m manager) GetSecret(id string) (string, error) {
	v, err := m.storage.Read(id)
	if err != nil {
		slog.Error("Failed to read from the storage:" + err.Error())
		return "", err
	}

	return v, nil
}

func getHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
