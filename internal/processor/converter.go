package processor

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

// ConvertFile bierze plik wejściowy, zmienia rozmiar i zapisuje jako WebP w outputDir.
// To jest funkcja SYNCHRONICZNA (blokująca).
func ConvertFile(inputPath string, outputDir string) (int64, string, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return 0, "", fmt.Errorf("nie udało się otworzyć pliku: %w", err)
	}
	defer file.Close()

	var src image.Image

	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext == ".heic" || ext == ".heif" {

		src, err = goheif.Decode(file)
		if err != nil {
			return 0, "", fmt.Errorf("błąd dekodowania HEIC: %w", err)
		}
	} else {
		src, err = imaging.Decode(file)
		if err != nil {
			return 0, "", fmt.Errorf("nieznany format obrazu: %w", err)
		}
	}
	var dst image.Image

	if src.Bounds().Dx() > 2560 {
		// Jest za duży -> Skalujemy w dół do 2560px
		dst = imaging.Resize(src, 2560, 0, imaging.Lanczos)
	} else {
		dst = src
	}
	bounds := dst.Bounds()
	imgRGBA := image.NewRGBA(bounds)
	fileName := filepath.Base(inputPath)
	draw.Draw(imgRGBA, bounds, dst, bounds.Min, draw.Src)

	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	outPath := filepath.Join(outputDir, name+".webp")

	outFile, err := os.Create(outPath)
	if err != nil {
		return 0, "", fmt.Errorf("nie udało się utworzyć pliku wyjściowego: %w", err)
	}
	defer outFile.Close()

	err = webp.Encode(outFile, imgRGBA, &webp.Options{
		Lossless: false,
		Quality:  80,
	})
	if err != nil {
		os.Remove(outPath)
		return 0, "", fmt.Errorf("błąd kodowania WebP: %w", err)
	}
	stat, err := outFile.Stat()
	if err != nil {
		return 0, "", nil
	}

	return stat.Size(), outPath, nil
}
