package processor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func RunSimpleAsync(files []string, outputDir string) int64 {
	start := time.Now()
	var totalSize int64

	var wg sync.WaitGroup
	fmt.Printf("🚀 Rozpoczynam przetwarzanie %d plików (tryb: Simple Async)...\n", len(files))

	for _, file := range files {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			size, err := ConvertFile(f, outputDir)
			if err != nil {
				fmt.Printf("❌ Błąd podczas przetwarzania pliku %s: %v\n", f, err)
			} else {
				fmt.Printf("✅ Plik %s przetworzony\n", f)
				if size > 0 {
					atomic.AddInt64(&totalSize, size)
				}

			}
		}(file)
	}
	wg.Wait()
	fmt.Printf("\n⏱️  Czas całkowity: %s\n", time.Since(start))
	return atomic.LoadInt64(&totalSize)
}
