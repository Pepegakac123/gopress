package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/jdeng/goheif"
	"github.com/rwcarlsen/goexif/exif"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

// ConvertFile bierze plik wejściowy, zmienia rozmiar i zapisuje jako WebP w outputDir.
// To jest funkcja SYNCHRONICZNA (blokująca).
func ConvertFile(inputPath string, outputDir string, quality int, maxWidth int) (int64, string, error) {
	var src image.Image
	var err error

	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext == ".heic" || ext == ".heif" {
		file, err := os.Open(inputPath)
		if err != nil {
			return 0, "", fmt.Errorf("nie udało się otworzyć pliku: %w", err)
		}
		defer file.Close()

		src, err = goheif.Decode(file)
		if err != nil {
			return 0, "", fmt.Errorf("błąd dekodowania HEIC: %w", err)
		}

		// Obsługa rotacji EXIF dla HEIC
		if _, err := file.Seek(0, 0); err == nil {
			if exifData, err := goheif.ExtractExif(file); err == nil {
				if x, err := exif.Decode(bytes.NewReader(exifData)); err == nil {
					if tag, err := x.Get(exif.Orientation); err == nil {
						if orient, err := tag.Int(0); err == nil {
							switch orient {
							case 2:
								src = imaging.FlipH(src)
							case 3:
								src = imaging.Rotate180(src)
							case 4:
								src = imaging.FlipV(src)
							case 5:
								src = imaging.Transpose(src)
							case 6:
								src = imaging.Rotate270(src)
							case 7:
								src = imaging.Transverse(src)
							case 8:
								src = imaging.Rotate90(src)
							}
						}
					}
				}
			}
		}
	} else {
		// Używamy imaging.Open, który automatycznie obsługuje orientację EXIF
		src, err = imaging.Open(inputPath)
		if err != nil {
			return 0, "", fmt.Errorf("nieznany format obrazu: %w", err)
		}
	}

	var dst image.Image

	if src.Bounds().Dx() > maxWidth {
		dst = imaging.Resize(src, maxWidth, 0, imaging.Lanczos)
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
		Quality:  float32(quality),
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
