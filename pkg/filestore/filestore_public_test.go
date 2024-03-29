package filestore_test

import (
	"bytes"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/logger"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"os"
	"strings"
	"testing"
)

const (
	DataFilePath    = "../../testdata/temp_data.json"
	ExistentId      = "d0b2fa28502506a01bb074879eaa083d"
	ExistentValue   = "my_secret_value"
	NonExistentID   = "334c4a4c42fdb79d7ebc3e73b517e6f8"
	PatternTempFile = "filestore_test_*.json"
)

func TestRead(t *testing.T) {
	var logBuffer bytes.Buffer
	config := logger.Config{LogLevel: "DEBUG"}

	logger.InitLogger(config, &logBuffer)

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

	t.Run("Rejects when reading not existing id", func(t *testing.T) {
		_, err := fs.Read(NonExistentID)

		if err.Error() != "value is not present" {
			t.Errorf("expected not found message, got %s", err.Error())
		}

		if !strings.Contains(logBuffer.String(), "The id: "+NonExistentID+" was not found") {
			t.Errorf("expected to find not existent error message, got: %s", logBuffer.String())
		}
	})
}

func TestWrite(t *testing.T) {
	var logBuffer bytes.Buffer
	config := logger.Config{LogLevel: "DEBUG"}

	logger.InitLogger(config, &logBuffer)

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

	if !strings.Contains(logBuffer.String(), "Writing the secret with ID:"+newID) {
		t.Errorf("expected to find message with secret id, got: %s", logBuffer.String())
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
