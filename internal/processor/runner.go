package processor

import (
	"fmt"
	"sync"
	"time"
)

func RunSimpleAsync(files []string, outputDir string) {
	start := time.Now()

	var wg sync.WaitGroup
	fmt.Printf("🚀 Rozpoczynam przetwarzanie %d plików (tryb: Simple Async)...\n", len(files))

	for _, file := range files {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			err := ConvertFile(f, outputDir)
			if err != nil {
				fmt.Printf("❌ Błąd podczas przetwarzania pliku %s: %v\n", f, err)
			} else {
				fmt.Printf("✅ Plik %s przetworzony\n", f)
			}
		}(file)
	}
	wg.Wait()
	fmt.Printf("\n⏱️  Czas całkowity: %s\n", time.Since(start))
}
