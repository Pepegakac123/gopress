package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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
	if appConfig.WpDomain == "" || appConfig.FileBirdToken == "" {
		fmt.Println("Błąd: Musisz podać przynajmniej --wp-domain oraz --fb-token.")
		fmt.Println("Przykład: gopress links --wp-domain https://example.com --fb-token \"XXXX\"")
		os.Exit(1)
	}

	client := wordpress.NewClient(appConfig.WpDomain, appConfig.WpUser, appConfig.WpPassword, appConfig.FileBirdToken)

	fmt.Println("Łączenie z WordPress i pobieranie danych z FileBird...")
	folders, err := client.GetFileBirdFolders()
	if err != nil {
		fmt.Printf("Błąd: %v\n", err)
		os.Exit(1)
	}

	// Heurystyka: szukamy folderu głównego rekurencyjnie
	targetName := "5 zdjęć do każdej usługi"
	shortTargetName := "5 zdjęć"

	root := findRootRecursive(folders, targetName, shortTargetName)
	if root == nil {
		fmt.Printf("Nie znaleziono folderu o nazwie zawierającej \"%s\" lub \"%s\".\n", targetName, shortTargetName)
		os.Exit(1)
	}

	fmt.Printf("Znaleziono folder: %s (ID: %d)\n", root.Text, root.ID)

	if len(root.Children) == 0 {
		fmt.Println("Brak podfolderów w folderze głównym.")
		return
	}

	// Menu TUI
	showTUI(client, root)
}

func findRootRecursive(folders []wordpress.FbFolder, name, shortName string) *wordpress.FbFolder {
	n := strings.ToLower(name)
	ns := strings.ToLower(shortName)

	for _, f := range folders {
		lowerText := strings.ToLower(f.Text)
		if strings.Contains(lowerText, n) || (ns != "" && strings.Contains(lowerText, ns)) {
			return &f
		}
		if len(f.Children) > 0 {
			res := findRootRecursive(f.Children, name, shortName)
			if res != nil {
				return res
			}
		}
	}
	return nil
}

func showTUI(client *wordpress.Client, root *wordpress.FbFolder) {
	// Przygotuj listę nazw podfolderów do wyboru
	subfoldersMap := make(map[string]int)
	subfolderNames := []string{}

	for _, sub := range root.Children {
		subfoldersMap[sub.Text] = sub.ID
		subfolderNames = append(subfolderNames, sub.Text)
	}
	sort.Strings(subfolderNames)

	const exitOpt = "❌ Wyjdź"
	const allOpt = "📁 Wszystkie (lista zbiorcza)"

	for {
		options := append([]string{allOpt}, subfolderNames...)
		options = append(options, exitOpt)

		var selected string
		prompt := &survey.Select{
			Message: "Wybierz podfolder, aby zobaczyć linki:",
			Options: options,
			PageSize: 15,
		}

		err := survey.AskOne(prompt, &selected)
		if err != nil {
			if err == terminal.InterruptErr {
				break
			}
			fmt.Printf("Błąd: %v\n", err)
			break
		}

		if selected == exitOpt {
			break
		}

		if selected == allOpt {
			printAllFolders(client, root.Children)
		} else {
			folderID := subfoldersMap[selected]
			printFolderLinks(client, selected, folderID)
		}

		fmt.Println("\n---")
	}
}

func printFolderLinks(client *wordpress.Client, name string, id int) {
	fmt.Printf("\n🔍 Pobieranie linków dla: %s...\n", name)
	attachmentIDs, err := client.GetFileBirdAttachmentIDs(id)
	if err != nil {
		fmt.Printf("  Błąd: %v\n", err)
		return
	}

	if len(attachmentIDs) == 0 {
		fmt.Println("  (Brak zdjęć w tym folderze)")
		return
	}

	for _, aid := range attachmentIDs {
		url, err := client.GetMediaSourceURL(aid)
		if err != nil {
			fmt.Printf("  Błąd URL dla ID %s: %v\n", aid, err)
			continue
		}
		fmt.Println(url)
	}
}

func printAllFolders(client *wordpress.Client, subs []wordpress.FbFolder) {
	fmt.Println("\n📦 Generowanie pełnej listy linków...")
	for _, sub := range subs {
		fmt.Printf("\n%s\n", sub.Text)
		attachmentIDs, err := client.GetFileBirdAttachmentIDs(sub.ID)
		if err != nil {
			fmt.Printf("  Błąd: %v\n", err)
			continue
		}
		if len(attachmentIDs) == 0 {
			fmt.Println("  (Brak zdjęć)")
			continue
		}
		for _, aid := range attachmentIDs {
			url, err := client.GetMediaSourceURL(aid)
			if err != nil {
				continue
			}
			fmt.Println(url)
		}
	}
}
