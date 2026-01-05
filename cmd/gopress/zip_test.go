package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestUnzip(t *testing.T) {
	// 1. Create a temp directory for our test workspace
	workDir, err := os.MkdirTemp("", "gopress_zip_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(workDir)

	zipPath := filepath.Join(workDir, "test.zip")
	destDir := filepath.Join(workDir, "extracted")

	// 2. Create a ZIP file with some content
	err = createTestZip(zipPath)
	if err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	// 3. Run the unzip function
	if err := unzip(zipPath, destDir); err != nil {
		t.Fatalf("Unzip failed: %v", err)
	}

	// 4. Verify contents
	expectedFiles := []string{
		"root.txt",
		filepath.Join("subdir", "nested.txt"),
	}

	for _, f := range expectedFiles {
		fullPath := filepath.Join(destDir, f)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file not found: %s", fullPath)
		}
	}

	// Verify content of one file
	content, err := os.ReadFile(filepath.Join(destDir, "root.txt"))
	if err != nil {
		t.Errorf("Failed to read extracted file: %v", err)
	}
	if string(content) != "root content" {
		t.Errorf("Expected 'root content', got '%s'", string(content))
	}
}

// Helper to create a zip file
func createTestZip(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	// Add file 1: root.txt
	f1, err := w.Create("root.txt")
	if err != nil {
		return err
	}
	_, err = f1.Write([]byte("root content"))
	if err != nil {
		return err
	}

	// Add file 2: subdir/nested.txt
	f2, err := w.Create("subdir/nested.txt")
	if err != nil {
		return err
	}
	_, err = f2.Write([]byte("nested content"))
	if err != nil {
		return err
	}

	return nil
}
