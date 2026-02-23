package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/Pepegakac123/gopress/internal/wordpress"
	"github.com/spf13/cobra"
)

var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "Wypisuje linki do zdjęć z folderów FileBird",
	Long: `Ta komenda wyszukuje folder "5 zdjęć do każdej usługi" (lub o podobnej nazwie) w FileBird,
przeszukuje jego podfoldery i dla każdego z nich wypisuje listę bezpośrednich linków do zdjęć.`,
	Run: runLinksLister,
}

func runLinksLister(cmd *cobra.Command, args []string) {
	if appConfig.WpDomain == "" || appConfig.WpUser == "" || appConfig.WpPassword == "" || appConfig.FileBirdToken == "" {
		fmt.Println("Błąd: Musisz podać --wp-domain, --wp-user, --wp-secret oraz --fb-token.")
		fmt.Println("Przykład: gopress links --wp-domain https://example.com --wp-user admin --wp-secret \"XXXX XXXX XXXX\" --fb-token \"XXXX\"")
		os.Exit(1)
	}

	client := wordpress.NewClient(appConfig.WpDomain, appConfig.WpUser, appConfig.WpPassword, appConfig.FileBirdToken)

	fmt.Println("Łączenie z WordPress i pobieranie danych z FileBird...")
	results, err := client.ListImagesFromFolder("5 zdjęć do każdej usługi", "5 zdjęć")
	if err != nil {
		fmt.Printf("Błąd: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("Nie znaleziono żadnych podfolderów ani zdjęć w folderze głównym.")
		return
	}

	// Sortujemy nazwy podfolderów dla ładniejszego wyjścia
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, folderName := range keys {
		links := results[folderName]
		fmt.Printf("\n%s\n", folderName)
		if len(links) == 0 {
			fmt.Println("  (Brak zdjęć)")
			continue
		}
		for _, url := range links {
			fmt.Println(url)
		}
	}
}
