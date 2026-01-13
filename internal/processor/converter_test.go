package processor

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create a sample JPEG input file
	inputPath := filepath.Join(tmpDir, "test_image.jpg")
	createSampleImage(t, inputPath, 100, 200) // Vertical image

	// Define output directory
	outputDir := filepath.Join(tmpDir, "output")
	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Run ConvertFile
	quality := 80
	maxWidth := 150
	_, outPath, err := ConvertFile(inputPath, outputDir, quality, maxWidth)

	if err != nil {
		t.Fatalf("ConvertFile failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created at %s", outPath)
	}

	// Verify output is WebP (by extension and check content if possible, but extension is a start)
	if filepath.Ext(outPath) != ".webp" {
		t.Errorf("Expected .webp extension, got %s", filepath.Ext(outPath))
	}

	// test resizing.
	inputPath2 := filepath.Join(tmpDir, "test_image_large.jpg")
	createSampleImage(t, inputPath2, 400, 800)

	size2, _, err2 := ConvertFile(inputPath2, outputDir, quality, maxWidth)
	if err2 != nil {
		t.Fatalf("ConvertFile (large) failed: %v", err2)
	}

	if size2 == 0 {
		t.Error("Returned size is 0")
	}
}

func createSampleImage(t *testing.T, path string, width, height int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 100, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create input image: %v", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatalf("Failed to encode JPEG: %v", err)
	}
}
