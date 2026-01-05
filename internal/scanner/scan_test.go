package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindImages(t *testing.T) {
	// 1. Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 2. Define files to create
	filesToCreate := []struct {
		path    string
		content []byte
		isValid bool
	}{
		{"image1.jpg", []byte("fake-image"), true},
		{"doc.txt", []byte("text"), false},
		{"nested/image2.png", []byte("fake-image"), true},
		{"nested/deep/image3.HEIC", []byte("fake-image"), true}, // Uppercase check
		{"nested/script.sh", []byte("echo hello"), false},
		{"image4.webp", []byte("fake-image"), true},
	}

	// 3. Create the files on disk
	for _, f := range filesToCreate {
		fullPath := filepath.Join(tmpDir, f.path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create subdir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, f.content, 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	// 4. Run the function under test
	images, totalSize, err := FindImages(tmpDir)
	if err != nil {
		t.Fatalf("FindImages returned error: %v", err)
	}

	// 5. Assertions
	expectedCount := 4 // jpg, png, HEIC, webp
	if len(images) != expectedCount {
		t.Errorf("Expected %d images, got %d", expectedCount, len(images))
	}

	// Verify total size calculation (4 files * 10 bytes "fake-image")
	// Note: "fake-image" is 10 bytes.
	// image1.jpg (10) + image2.png (10) + image3.HEIC (10) + image4.webp (10) = 40 bytes
	expectedSize := int64(40)
	if totalSize != expectedSize {
		t.Errorf("Expected total size %d, got %d", expectedSize, totalSize)
	}

	// Verify specific file existence in results
	foundMap := make(map[string]bool)
	for _, img := range images {
		foundMap[filepath.Base(img)] = true
	}

	if !foundMap["image1.jpg"] {
		t.Error("image1.jpg not found in results")
	}
	if !foundMap["image3.HEIC"] {
		t.Error("image3.HEIC (uppercase) not found in results")
	}
}
