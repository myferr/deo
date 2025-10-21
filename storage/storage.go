package storage

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type DocumentFilter struct {
	Field string
	Value string
}

type DocumentSort struct {
	Field string
	Order string // "asc" or "desc"
}

type DocumentPagination struct {
	Limit  int
	Offset int
}

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

func ListDocuments(dbName, collectionName string, filters []DocumentFilter, sortParams *DocumentSort, pagination *DocumentPagination) ([]map[string]interface{}, error) {
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

	var allDocuments []map[string]interface{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".msgpack" {
			docID := entry.Name()[:len(entry.Name())-len(".msgpack")]
			doc, err := LoadDocument(dbName, collectionName, docID)
			if err != nil {
				// Skip corrupted or unreadable files
				continue
			}
			allDocuments = append(allDocuments, doc)
		}
	}

	// Apply filters
	filteredDocuments := []map[string]interface{}(nil)
	if len(filters) > 0 {
		for _, doc := range allDocuments {
			match := true
			for _, f := range filters {
				if val, ok := doc[f.Field]; !ok || val != f.Value {
					match = false
					break
				}
			}
			if match {
				filteredDocuments = append(filteredDocuments, doc)
			}
		}
	} else {
		filteredDocuments = allDocuments
	}

	// Apply sorting
	if sortParams != nil && sortParams.Field != "" {
		sortParams.Order = strings.ToLower(sortParams.Order)
		sortParams.Order = strings.TrimSpace(sortParams.Order)
		sortParams.Order = strings.ReplaceAll(sortParams.Order, " ", "")

		sort.Slice(filteredDocuments, func(i, j int) bool {
			valI, okI := filteredDocuments[i][sortParams.Field]
			valJ, okJ := filteredDocuments[j][sortParams.Field]

			if !okI && !okJ {
				return false // Both don't have the field, order doesn't matter
			}
			if !okI {
				return sortParams.Order == "desc" // If i doesn't have field, it comes after j for asc, before for desc
			}
			if !okJ {
				return sortParams.Order == "asc" // If j doesn't have field, it comes after i for desc, before for asc
			}

			switch vI := valI.(type) {
			case int:
				if vJ, ok := valJ.(int); ok {
					if sortParams.Order == "desc" {
						return vI > vJ
					}
					return vI < vJ
				}
			case float64:
				if vJ, ok := valJ.(float64); ok {
					if sortParams.Order == "desc" {
						return vI > vJ
					}
					return vI < vJ
				}
			case string:
				if vJ, ok := valJ.(string); ok {
					if sortParams.Order == "desc" {
						return vI > vJ
					}
					return vI < vJ
				}
			}
			return false // Fallback if types are not comparable or not handled
		})
	}

	// Apply pagination
	start := 0
	end := len(filteredDocuments)

	if pagination != nil {
		if pagination.Offset > 0 && pagination.Offset < end {
			start = pagination.Offset
		}
		if pagination.Limit > 0 {
			end = start + pagination.Limit
			if end > len(filteredDocuments) {
				end = len(filteredDocuments)
			}
		}
	}

	return filteredDocuments[start:end], nil
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
