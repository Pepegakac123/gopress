package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

// func RunSimpleAsync(files []string, outputDir string) int64 {
// 	start := time.Now()
// 	var totalSize int64

// 	var wg sync.WaitGroup
// 	fmt.Printf("🚀 Rozpoczynam przetwarzanie %d plików (tryb: Simple Async)...\n", len(files))

// 	for _, file := range files {
// 		wg.Add(1)

// 		go func(f string) {
// 			defer wg.Done()
// 			size, err := ConvertFile(f, outputDir)
// 			if err != nil {
// 				fmt.Printf("❌ Błąd podczas przetwarzania pliku %s: %v\n", f, err)
// 			} else {
// 				fmt.Printf("✅ Plik %s przetworzony\n", f)
// 				if size > 0 {
// 					atomic.AddInt64(&totalSize, size)
// 				}

// 			}
// 		}(file)
// 	}
// 	wg.Wait()
// 	fmt.Printf("\n⏱️  Czas całkowity: %s\n", time.Since(start))
// 	return atomic.LoadInt64(&totalSize)
// }

// RunWorkerPool to bezpieczna wersja przetwarzania równoległego.
// Zwraca łączny rozmiar przetworzonych plików (w bajtach).
func RunWorkerPool(ctx context.Context, files []string, inputRoot string, outputRoot string, quality, maxWidth int, deleteOriginals bool) (int64, []string) {
	start := time.Now()
	var totalSize int64
	var convertedFiles []string
	var mu sync.Mutex
	totalFiles := len(files)
	numWorkers := runtime.NumCPU()
	fmt.Printf("🚀 Rozpoczynam przetwarzanie %d plików (Pełna moc procesora)...\n", totalFiles)

	bar := progressbar.Default(int64(totalFiles))

	jobs := make(chan string, totalFiles)
	results := make(chan error, totalFiles)

	var wg sync.WaitGroup
	for i := range numWorkers {
		wg.Add(1)
		go worker(ctx, i, jobs, results, inputRoot, outputRoot, quality, maxWidth, deleteOriginals, &wg, &totalSize, &convertedFiles, &mu)

	}
	go func() {
		for _, file := range files {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- file:
			}
		}
		close(jobs)
	}()

	var errorCount int
	done := make(chan bool)
	go func() {
		for err := range results {
			bar.Add(1)
			if err != nil {
				errorCount++
			}
		}
		done <- true
	}()
	wg.Wait()
	close(results)
	<-done
	fmt.Printf("\n\n🏁 Zakończono w %s\n", time.Since(start).Round(time.Millisecond))
	if ctx.Err() == context.Canceled {
		fmt.Println("\n Operacja anulowana przez użytkownika!")
	} else {
		fmt.Printf("\n\nZakończono w %s\n", time.Since(start).Round(time.Millisecond))
	}

	if errorCount > 0 {
		fmt.Printf("⚠️ Liczba błędów konwersji: %d\n", errorCount)
	}
	return atomic.LoadInt64(&totalSize), convertedFiles
}

// worker wykonuje zadania z kanału jobs
func worker(ctx context.Context, id int, jobs <-chan string, results chan<- error, inputRoot, outputRoot string, quality, maxWidth int, deleteOriginals bool, wg *sync.WaitGroup, totalSize *int64, convertedFiles *[]string, mu *sync.Mutex) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case filePath, ok := <-jobs:
			if !ok {
				return
			}
			relPath, err := filepath.Rel(inputRoot, filePath)
			if err != nil {
				relPath = filepath.Base(filePath)
			}
			targetDir := filepath.Join(outputRoot, filepath.Dir(relPath))

			if err := os.MkdirAll(targetDir, 0755); err != nil {
				results <- fmt.Errorf("błąd tworzenia katalogu %s: %w", targetDir, err)
				continue
			}
			size, outPath, err := ConvertFile(filePath, targetDir, quality, maxWidth)

			if err == nil {
				mu.Lock()
				*convertedFiles = append(*convertedFiles, outPath)
				mu.Unlock()
				atomic.AddInt64(totalSize, int64(size))
				if deleteOriginals {
					if rmErr := os.Remove(filePath); rmErr != nil {
						results <- fmt.Errorf("skonwertowano, ale błąd usuwania źródła %s: %w", filePath, rmErr)
						continue
					}
					dir := filepath.Dir(filePath)
					_ = os.Remove(dir)
				}
			}
			results <- err
		}
	}
}
