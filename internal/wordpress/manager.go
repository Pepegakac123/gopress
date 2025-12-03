package wordpress

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// FolderManager dba o odwzorowanie struktury katalogów i zapobiega duplikatom.
type FolderManager struct {
	client *Client
	rootID int            // Opcjonalne ID folderu startowego w WP
	cache  map[string]int // Cache: "2024/wakacje" -> ID 152
	mu     sync.Mutex     // Chroni cache przy wielowątkowości
}

func NewFolderManager(client *Client, rootID int) *FolderManager {
	return &FolderManager{
		client: client,
		rootID: rootID,
		cache:  make(map[string]int),
	}
}

// GetFolderID zwraca ID folderu w FileBird dla danej ścieżki względnej (np. "lato/morze").
// Jeśli folder nie istnieje, zostanie automatycznie utworzony (wraz z rodzicami).
func (fm *FolderManager) GetFolderID(relPath string) (int, error) {

	relPath = filepath.ToSlash(relPath)
	relPath = strings.TrimPrefix(relPath, "./")

	if relPath == "" || relPath == "." {
		return fm.rootID, nil
	}

	fm.mu.Lock()
	if id, exists := fm.cache[relPath]; exists {
		fm.mu.Unlock()
		return id, nil
	}
	fm.mu.Unlock()

	parentDir := filepath.Dir(relPath)
	folderName := filepath.Base(relPath)

	parentID := fm.rootID

	if parentDir != "." && parentDir != "/" && parentDir != "" {
		var err error
		parentID, err = fm.GetFolderID(parentDir)
		if err != nil {
			return 0, err
		}
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Double-check: może inny wątek stworzył ten folder w ułamku sekundy, gdy my szukaliśmy rodzica?
	if id, exists := fm.cache[relPath]; exists {
		return id, nil
	}

	// Wywołujemy API FileBirda
	// fmt.Printf("\n📂 Tworzę folder FileBird: '%s' ... ", relPath)
	newID, err := fm.client.CreateFolder(folderName, parentID)
	if err != nil {
		return 0, fmt.Errorf("nie udało się stworzyć folderu '%s': %w", folderName, err)
	}
	fm.cache[relPath] = newID

	return newID, nil
}
