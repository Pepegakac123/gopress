package gallery

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyAndMerge(convertedFiles []string, outputRoot string, galleryRoot string) error {
	if _, err := os.Stat(galleryRoot); os.IsNotExist(err) {
		return fmt.Errorf("galeria pod adresem %s nie istnieje", galleryRoot)
	}

	existingFiles := make(map[string]map[string]bool)

	for _, file := range convertedFiles {
		relPath, err := filepath.Rel(outputRoot, file)
		if err != nil {
			return fmt.Errorf("nie udalo sie wyznaczyc sciezki wzglednej dla %s: %w", file, err)
		}

		relDir := filepath.Dir(relPath)
		fileName := filepath.Base(relPath)

		var targetDir string
		if relDir == "." || relDir == "" {
			targetDir = galleryRoot
		} else {
			resolvedRelDir, err := resolvePathCaseInsensitive(galleryRoot, relDir)
			if err != nil {
				return fmt.Errorf("blad dopasowania folderu %s: %w", relDir, err)
			}
			targetDir = filepath.Join(galleryRoot, resolvedRelDir)
		}

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("nie udalo sie utworzyc folderu %s: %w", targetDir, err)
		}

		filesMap, err := getExistingFilesForPath(targetDir, existingFiles)
		if err != nil {
			return fmt.Errorf("nie udalo sie pobrac plikow z %s: %w", targetDir, err)
		}

		if filesMap[strings.ToLower(fileName)] {
			continue
		}

		destPath := filepath.Join(targetDir, fileName)
		if err := copyFile(file, destPath); err != nil {
			return fmt.Errorf("nie udalo sie skopiowac pliku %s do %s: %w", file, destPath, err)
		}

		filesMap[strings.ToLower(fileName)] = true
	}

	return nil
}

func getExistingFilesForPath(dirPath string, cache map[string]map[string]bool) (map[string]bool, error) {
	key := strings.ToLower(filepath.Clean(dirPath))
	if filesMap, ok := cache[key]; ok {
		return filesMap, nil
	}

	filesMap := make(map[string]bool)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			cache[key] = filesMap
			return filesMap, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filesMap[strings.ToLower(entry.Name())] = true
		}
	}
	cache[key] = filesMap
	return filesMap, nil
}

func resolvePathCaseInsensitive(baseDir string, relPath string) (string, error) {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	currentRel := ""

	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}

		currentFullDir := baseDir
		if currentRel != "" {
			currentFullDir = filepath.Join(baseDir, currentRel)
		}

		entries, err := os.ReadDir(currentFullDir)
		if err != nil {
			if os.IsNotExist(err) {
				currentRel = filepath.Join(currentRel, part)
				continue
			}
			return "", err
		}

		found := false
		for _, entry := range entries {
			if entry.IsDir() && strings.EqualFold(entry.Name(), part) {
				currentRel = filepath.Join(currentRel, entry.Name())
				found = true
				break
			}
		}

		if !found {
			currentRel = filepath.Join(currentRel, part)
		}
	}

	return currentRel, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}
