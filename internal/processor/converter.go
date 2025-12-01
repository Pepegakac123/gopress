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
)

// ConvertFile bierze plik wejściowy, zmienia rozmiar i zapisuje jako WebP w outputDir.
// To jest funkcja SYNCHRONICZNA (blokująca).
func ConvertFile(inputPath string, outputDir string) error {

	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("nie udało się otworzyć pliku: %w", err)
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
		return fmt.Errorf("nie udało się utworzyć pliku wyjściowego: %w", err)
	}
	defer outFile.Close()

	err = webp.Encode(outFile, imgRGBA, &webp.Options{
		Lossless: false,
		Quality:  80,
	})
	if err != nil {
		os.Remove(outPath)
		return fmt.Errorf("błąd kodowania WebP: %w", err)
	}

	return nil
}
