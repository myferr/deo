package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/myferr/deo/storage"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

func main() {
	fmt.Println(ColorYellow + "--- Running Unit Tests ---" + ColorReset)

	runTest("TestGetDocPath", testGetDocPath)
	runTest("TestSaveAndLoadDocument", testSaveAndLoadDocument)

	fmt.Println(ColorYellow + "--- All Unit Tests Passed ---" + ColorReset)
}

func runTest(name string, testFunc func() error) {
	fmt.Printf("Running test: %s... ", name)
	if err := testFunc(); err != nil {
		fmt.Println(ColorRed + "FAILED" + ColorReset)
		fmt.Printf("  Error: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println(ColorGreen + "PASSED" + ColorReset)
	}
}

func testGetDocPath() error {
	path, err := storage.GetDocPath("testdb", "testcollection", "testid")
	if err != nil {
		return err
	}
	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, ".deo", "testdb", "testcollection", "testid.msgpack")
	if path != expectedPath {
		return fmt.Errorf("expected path %s, got %s", expectedPath, path)
	}
	return nil
}

func testSaveAndLoadDocument() error {
	db := "testdb_unit"
	coll := "testcoll_unit"
	id := "testid_unit"
	doc := map[string]interface{}{"hello": "world", "_id": id}

	err := storage.SaveDocument(db, coll, id, doc)
	if err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	loadedDoc, err := storage.LoadDocument(db, coll, id)
	if err != nil {
		return fmt.Errorf("failed to load document: %w", err)
	}

	if loadedDoc["hello"] != "world" {
		return fmt.Errorf("document content does not match")
	}

	// cleanup
	path, _ := storage.GetDocPath(db, coll, id)
	os.RemoveAll(filepath.Dir(path))

	return nil
}
