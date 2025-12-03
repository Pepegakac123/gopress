package uploader

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/Pepegakac123/gopress/internal/wordpress"
	"github.com/schollz/progressbar/v3"
)

// Run rozpoczyna współbieżny upload plików do WordPressa.
func Run(ctx context.Context, client *wordpress.Client, files []string, outputDir string, useFileBird bool, fbRootID int) {
	totalFiles := len(files)
	numWorkers := 4 // Na sztywno żeby uniknąc rate limitingu

	fmt.Printf("🚀 Wysyłanie %d plików do WordPressa (Upload Workers: %d)...\n", totalFiles, numWorkers)

	var folderMgr *wordpress.FolderManager
	if useFileBird {
		fmt.Println("📂 Obsługa folderów FileBird: AKTYWNA")
		folderMgr = wordpress.NewFolderManager(client, fbRootID)
	}

	bar := progressbar.Default(int64(totalFiles))

	// Kanały i liczniki
	jobs := make(chan string, totalFiles)
	var uploadErrors int64
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go worker(ctx, jobs, client, folderMgr, outputDir, &wg, &uploadErrors, bar)
	}

	// Wrzucanie zadań
	go func() {
		for _, filePath := range files {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- filePath:
			}
		}
		close(jobs)
	}()

	wg.Wait()

	// Raport końcowy
	fmt.Println("\n")
	if uploadErrors > 0 {
		fmt.Printf("⚠️  Zakończono z błędami uploadu: %d\n", uploadErrors)
	} else {
		fmt.Println("🎉 Sukces! Wszystkie pliki wysłane.")
	}
}

func worker(ctx context.Context, jobs <-chan string, client *wordpress.Client, folderMgr *wordpress.FolderManager, outputDir string, wg *sync.WaitGroup, errors *int64, bar *progressbar.ProgressBar) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case filePath, ok := <-jobs:
			if !ok {
				return
			}
			resp, err := client.UploadFile(filePath)
			bar.Add(1)

			if err != nil {
				atomic.AddInt64(errors, 1)
				continue
			}

			// 2. Obsługa folderów (tylko jeśli manager istnieje i upload się udał)
			if folderMgr != nil {

				relPath, err := filepath.Rel(outputDir, filePath)
				if err == nil {
					dirName := filepath.Dir(relPath)
					folderID, err := folderMgr.GetFolderID(dirName)

					if err == nil && folderID > 0 {
						client.SetAttachmentFolder(folderID, []int{resp.ID})
					}
				}
			}
		}
	}
}
