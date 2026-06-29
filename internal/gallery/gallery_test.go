package gallery

import (
	"os"
	"path/filepath"
	"testing"
)

func fileExistsExactCase(parentDir string, expectedName string) bool {
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.Name() == expectedName {
			return true
		}
	}
	return false
}

func TestCopyAndMerge(t *testing.T) {
	outputDir, err := os.MkdirTemp("", "gopress_output_test")
	if err != nil {
		t.Fatalf("failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	galleryDir, err := os.MkdirTemp("", "gopress_gallery_test")
	if err != nil {
		t.Fatalf("failed to create temp gallery dir: %v", err)
	}
	defer os.RemoveAll(galleryDir)

	filesToCreate := []struct {
		relPath string
		content string
	}{
		{"Malowanie/foto1.webp", "data1"},
		{"Malowanie/foto2.webp", "data2"},
		{"Instalacje/foto3.webp", "data3"},
		{"ogrodzenia/płot.webp", "data_plot"},
	}

	var convertedFiles []string
	for _, tc := range filesToCreate {
		fullPath := filepath.Join(outputDir, tc.relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(tc.content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		convertedFiles = append(convertedFiles, fullPath)
	}

	existingInGallery := []struct {
		relPath string
		content string
	}{
		{"malowanie/FOTO1.webp", "existing_data1"},
		{"malowanie/existing.webp", "existing_data2"},
		{"Ogrodzenia/other.webp", "other"},
	}

	for _, eg := range existingInGallery {
		fullPath := filepath.Join(galleryDir, eg.relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(eg.content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	if err := CopyAndMerge(convertedFiles, outputDir, galleryDir); err != nil {
		t.Fatalf("CopyAndMerge failed: %v", err)
	}

	// 1. Check folder 'malowanie' exists in its original lowercase form
	if !fileExistsExactCase(galleryDir, "malowanie") {
		t.Errorf("expected folder 'malowanie' to exist with exact case")
	}

	// 2. Check that uppercase 'Malowanie' is NOT a separate folder in gallery
	// Note: on Windows it might look like it exists if we query case-insensitively,
	// but fileExistsExactCase only checks exact name matches in ReadDir.
	if fileExistsExactCase(galleryDir, "Malowanie") {
		t.Errorf("did not expect a separate folder 'Malowanie' to exist (should have merged with 'malowanie')")
	}

	malowaniePath := filepath.Join(galleryDir, "malowanie")

	// 3. FOTO1.webp must still exist in malowanie (original case)
	if !fileExistsExactCase(malowaniePath, "FOTO1.webp") {
		t.Errorf("expected FOTO1.webp to exist in 'malowanie'")
	}

	// 4. foto1.webp (lowercase) must NOT have been copied/created (avoid duplicate)
	if fileExistsExactCase(malowaniePath, "foto1.webp") {
		t.Errorf("did not expect duplicate file 'foto1.webp' to be created")
	}

	// 5. foto2.webp must exist in malowanie
	if !fileExistsExactCase(malowaniePath, "foto2.webp") {
		t.Errorf("expected foto2.webp to exist in 'malowanie'")
	}

	// 6. Check 'Ogrodzenia' exists with its exact case
	if !fileExistsExactCase(galleryDir, "Ogrodzenia") {
		t.Errorf("expected folder 'Ogrodzenia' to exist with exact case")
	}

	// 7. Check that lowercase 'ogrodzenia' is NOT a separate folder in gallery
	if fileExistsExactCase(galleryDir, "ogrodzenia") {
		t.Errorf("did not expect a separate folder 'ogrodzenia' to exist")
	}

	ogrodzeniaPath := filepath.Join(galleryDir, "Ogrodzenia")

	// 8. płot.webp must exist in 'Ogrodzenia'
	if !fileExistsExactCase(ogrodzeniaPath, "płot.webp") {
		t.Errorf("expected 'płot.webp' to exist in 'Ogrodzenia'")
	}

	// 9. 'Instalacje' must exist with its exact case
	if !fileExistsExactCase(galleryDir, "Instalacje") {
		t.Errorf("expected folder 'Instalacje' to exist")
	}
}
