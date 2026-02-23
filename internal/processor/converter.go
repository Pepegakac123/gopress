package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"io"
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

	// Optymalizacja: Jeśli plik to WebP i jest mniejszy niż 400KB, kopiujemy go bez zmian
	ext := strings.ToLower(filepath.Ext(inputPath))
	info, err := os.Stat(inputPath)
	if err != nil {
		return 0, "", fmt.Errorf("nie udało się pobrać informacji o pliku: %w", err)
	}

	quality = adjustQuality(info.Size(), quality)

	if ext == ".webp" && info.Size() < 400*1024 {
		fileName := filepath.Base(inputPath)
		outPath := filepath.Join(outputDir, fileName)

		sourceFile, err := os.Open(inputPath)
		if err != nil {
			return 0, "", fmt.Errorf("nie udało się otworzyć pliku źródłowego: %w", err)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(outPath)
		if err != nil {
			return 0, "", fmt.Errorf("nie udało się utworzyć pliku docelowego: %w", err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, sourceFile); err != nil {
			return 0, "", fmt.Errorf("błąd kopiowania pliku: %w", err)
		}

		return info.Size(), outPath, nil
	}

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
				// Pliki HEIC w boksie Exif mają 4-bajtowy offset przed nagłówkiem TIFF (II lub MM).
				// Biblioteka goexif oczekuje nagłówka TIFF na samym początku.
				if len(exifData) > 4 {
					if (exifData[4] == 'I' && exifData[5] == 'I') || (exifData[4] == 'M' && exifData[5] == 'M') {
						exifData = exifData[4:]
					}
				}

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
		// Używamy imaging.Open z opcją AutoOrientation(true), aby automatycznie obrócić obraz na podstawie EXIF.
		// Domyślnie ta opcja jest wyłączona, co powodowało problemy ze zdjęciami pionowymi.
		src, err = imaging.Open(inputPath, imaging.AutoOrientation(true))
		if err != nil {
			return 0, "", fmt.Errorf("nie udało się otworzyć obrazu: %w", err)
		}
	}

	var dst image.Image

	if src.Bounds().Dx() > maxWidth {
		dst = imaging.Resize(src, maxWidth, 0, imaging.Lanczos)
	} else {
		dst = src
	}

	bounds := dst.Bounds()
	// Tworzymy nowy obraz RGBA o wymiarach dst, zaczynający się od (0,0).
	// To zapewnia "czysty" obraz bez metadanych i poprawną strukturę dla kodera WebP.
	imgRGBA := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(), dst, bounds.Min, draw.Src)

	fileName := filepath.Base(inputPath)
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

func adjustQuality(fileSize int64, quality int) int {
	const size4_5MB = 4718592 // 4.5 * 1024 * 1024
	const size10MB = 10485760 // 10 * 1024 * 1024

	if fileSize > size10MB {
		if quality > 60 {
			return 60
		}
	} else if fileSize > size4_5MB {
		if quality > 75 {
			return 75
		}
	}
	return quality
}
