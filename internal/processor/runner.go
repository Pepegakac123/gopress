package processor

import (
	"fmt"
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
func RunWorkerPool(files []string, outputDir string) int64 {
	start := time.Now()
	var totalSize int64

	totalFiles := len(files)
	numWorkers := runtime.NumCPU()
	fmt.Printf("🚀 Rozpoczynam przetwarzanie %d plików (tryb: Worker Pool)...\n", totalFiles)

	bar := progressbar.Default(int64(totalFiles))

	jobs := make(chan string, totalFiles)
	results := make(chan error, totalFiles)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, outputDir, &wg, &totalSize)
	}
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	var errorCount int
	for i := 0; i < totalFiles; i++ {
		err := <-results
		bar.Add(1)
		if err != nil {
			errorCount++
		}
	}
	wg.Wait()
	fmt.Printf("\n\n🏁 Zakończono w %s\n", time.Since(start).Round(time.Millisecond))
	if errorCount > 0 {
		fmt.Printf("⚠️  Liczba błędów: %d\n", errorCount)
	}
	return atomic.LoadInt64(&totalSize)
}

// worker wykonuje zadania z kanału jobs
func worker(id int, jobs <-chan string, results chan<- error, outputDir string, wg *sync.WaitGroup, totalSize *int64) {
	defer wg.Done()

	for filepath := range jobs {
		size, err := ConvertFile(filepath, outputDir)
		if err == nil {
			atomic.AddInt64(totalSize, int64(size))
		}
		results <- err
	}
}
