package filestore_test

import (
	"github.com/prog-supdex/mini-project/milestone-code/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/types"
	"os"
	"testing"
)

const (
	DataFilePath    = "../testdata/temp_data.json"
	SecretDataId    = "74e1d3f50df786aef9f602419dc88784"
	SecretDataValue = "SecretValue"
)

func TestInit(t *testing.T) {
	err := filestore.Init(DataFilePath)

	if err != nil {
		t.Error("file_store init function error:", err)
	}

	_, err = os.Stat(DataFilePath)

	if err != nil {
		t.Error("file was not created", err)
	}

	if filestore.FileStoreConfig.DataFilePath == "" {
		t.Error("FileStoreConfig.DataFilePath doesn`t contain a file path")
	}

	t.Cleanup(func() {
		os.Remove(DataFilePath)
	})
}

func TestWrite(t *testing.T) {
	defer func() {
		t.Cleanup(func() {
			os.Remove(DataFilePath)
		})
	}()

	err := writeTestFile()

	if err != nil {
		t.Error("Write function error:", err)
	}
}

func TestRead(t *testing.T) {
	defer func() {
		t.Cleanup(func() {
			os.Remove(DataFilePath)
		})
	}()

	err := writeTestFile()

	if err != nil {
		t.Error("Write function error:", err)
	}

	val, err := filestore.FileStoreConfig.Fs.Read(SecretDataId)

	if err != nil {
		t.Error("Read function error:", err)
	}

	if val != SecretDataValue {
		t.Error("incorrect value:", val, SecretDataValue)
	}
}

func writeTestFile() error {
	filestore.Init(DataFilePath)

	secretData := types.SecretData{
		Id:     SecretDataId,
		Secret: SecretDataValue,
	}

	return filestore.FileStoreConfig.Fs.Write(secretData)
}
