package storage

import (
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack/v5"
)

func GetDocPath(dbName, collectionName, docID string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".deo", dbName, collectionName, docID+".msgpack"), nil
}

func SaveDocument(dbName, collectionName, docID string, data map[string]interface{}) error {
	docPath, err := GetDocPath(dbName, collectionName, docID)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(docPath), 0755); err != nil {
		return err
	}

	packedData, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(docPath, packedData, 0644)
}

func LoadDocument(dbName, collectionName, docID string) (map[string]interface{}, error) {
	docPath, err := GetDocPath(dbName, collectionName, docID)
	if err != nil {
		return nil, err
	}

	packedData, err := os.ReadFile(docPath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = msgpack.Unmarshal(packedData, &data)
	return data, err
}
