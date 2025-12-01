package scanner

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".heic": true,
	".heif": true,
}

// FindImages przeszukuje rekursywnie katalog rootDir i zwraca listę ścieżek do obrazków.
func FindImages(rootDir string) ([]string, error) {
	var images []string

	// filepath.WalkDir to standardowy sposób na chodzenie po drzewie plików w Go (wydajniejszy niż Walk).
	// Przyjmuje funkcję callback, która jest wywoływana dla KAŻDEGO pliku/katalogu.
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		// 1. Obsługa błędu wejścia (jeśli nie możemy wejść do katalogu)
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileExtension := strings.ToLower(filepath.Ext(path))
			if AllowedExtensions[fileExtension] {
				images = append(images, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return images, nil
}
