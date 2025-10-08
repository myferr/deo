package storage

import (
	"os"
	"path/filepath"

	"github.com/vmihailenco/msgpack/v5"
)

func getDeoPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".deo"), nil
}

func ListDatabases() ([]string, error) {
	deoPath, err := getDeoPath()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(deoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var dbs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dbs = append(dbs, entry.Name())
		}
	}
	return dbs, nil
}

func CreateDatabase(dbName string) error {
	deoPath, err := getDeoPath()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(deoPath, dbName)
	return os.MkdirAll(dbPath, 0755)
}

func ListCollections(dbName string) ([]string, error) {
	deoPath, err := getDeoPath()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(deoPath, dbName)

	entries, err := os.ReadDir(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var collections []string
	for _, entry := range entries {
		if entry.IsDir() {
			collections = append(collections, entry.Name())
		}
	}
	return collections, nil
}

func ListDocuments(dbName, collectionName string) ([]map[string]interface{}, error) {
	deoPath, err := getDeoPath()
	if err != nil {
		return nil, err
	}
	collectionPath := filepath.Join(deoPath, dbName, collectionName)

	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]interface{}{}, nil
		}
		return nil, err
	}

	var documents []map[string]interface{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".msgpack" {
			docID := entry.Name()[:len(entry.Name())-len(".msgpack")]
			doc, err := LoadDocument(dbName, collectionName, docID)
			if err != nil {
				// Skip corrupted or unreadable files
				continue
			}
			documents = append(documents, doc)
		}
	}
	return documents, nil
}

func CreateCollection(dbName, collectionName string) error {
	deoPath, err := getDeoPath()
	if err != nil {
		return err
	}
	collectionPath := filepath.Join(deoPath, dbName, collectionName)
	return os.MkdirAll(collectionPath, 0755)
}

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

func DeleteDatabase(dbName string) error {
	deoPath, err := getDeoPath()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(deoPath, dbName)
	return os.RemoveAll(dbPath)
}

func DeleteCollection(dbName, collectionName string) error {
	deoPath, err := getDeoPath()
	if err != nil {
		return err
	}
	collectionPath := filepath.Join(deoPath, dbName, collectionName)
	return os.RemoveAll(collectionPath)
}
