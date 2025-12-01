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
	".bmp":  true,
	".tiff": true,
	".tif":  true,
}

// FindImages przeszukuje rekursywnie katalog rootDir i zwraca listę ścieżek do obrazków.
func FindImages(rootDir string) ([]string, int64, error) {
	var images []string
	var totalSize int64

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileExtension := strings.ToLower(filepath.Ext(path))
			if AllowedExtensions[fileExtension] {
				images = append(images, path)
				info, err := d.Info()
				if err == nil {
					totalSize += info.Size()
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return images, totalSize, nil
}
