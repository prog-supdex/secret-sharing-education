package filestore_test

import (
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"os"
	"testing"
)

const (
	DataFilePath    = "../../testdata/temp_data.json"
	ExistentId      = "d0b2fa28502506a01bb074879eaa083d"
	ExistentValue   = "my_secret_value"
	PatternTempFile = "filestore_test_*.json"
)

func TestRead(t *testing.T) {
	tempFilePath, cleanup := createTempFileForTesting(t, DataFilePath)
	defer cleanup()

	fs, err := filestore.New(tempFilePath)
	if err != nil {
		t.Fatalf("failed to create fileStore: %v", err)
	}

	val, err := fs.Read(ExistentId)
	if err != nil {
		t.Fatalf("failed to read from fileStore: %v", err)
	}
	if val != ExistentValue {
		t.Errorf("expected %s, got %s", ExistentValue, val)
	}
}

func TestWrite(t *testing.T) {
	tempFilePath, cleanup := createTempFileForTesting(t, DataFilePath)
	defer cleanup()

	fs, err := filestore.New(tempFilePath)
	if err != nil {
		t.Fatalf("failed to create fileStore: %v", err)
	}

	newID := "new_id"
	newValue := "new_value"
	err = fs.Write(secrets.SecretData{Id: newID, Secret: newValue})
	if err != nil {
		t.Fatalf("failed to write to fileStore: %v", err)
	}

	val, err := fs.Read(newID)
	if err != nil {
		t.Fatalf("failed to read from fileStore: %v", err)
	}
	if val != newValue {
		t.Errorf("expected %s, got %s", newValue, val)
	}
}

func createTempFileForTesting(t *testing.T, originalFilePath string) (string, func()) {
	t.Helper()

	originalContent, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("failed to read original file: %v", err)
	}

	tempFile, err := os.CreateTemp("", PatternTempFile)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tempFile.Write(originalContent); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}
}
