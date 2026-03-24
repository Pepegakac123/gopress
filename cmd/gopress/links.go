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

	// Heurystyka: szukamy folderu głównego rekurencyjnie (tylko przy pierwszym uruchomieniu pętli)
	targetName := "5 zdjęć do każdej usługi"
	shortTargetName := "5 zdjęć"
	root := findRootRecursive(folders, targetName, shortTargetName)

	for {
		if root == nil {
			var err error
			root, err = selectFolderInteractively(folders)
			if err != nil {
				if err == terminal.InterruptErr {
					fmt.Println("\nPrzerwano przez użytkownika.")
					return
				}
				fmt.Printf("Błąd wyboru folderu: %v\n", err)
				return
			}
			// Jeśli użytkownik wybrał wyjście (zwrócony root jest nil)
			if root == nil {
				return
			}
		}

		fmt.Printf("\n--- Aktualny folder: %s (ID: %d) ---\n", root.Text, root.ID)

		if len(root.Children) == 0 {
			printFolderLinks(client, root, false)
			fmt.Println("\n(Koniec folderu bez podfolderów)")
		} else {
			// Menu TUI
			showTUI(client, root)
		}

		// Resetujemy root, aby przy kolejnej iteracji użytkownik mógł wybrać inny folder
		root = nil
		fmt.Println("\n--- Powrót do menu wyboru folderu ---")
	}
}

func selectFolderInteractively(folders []wordpress.FbFolder) (*wordpress.FbFolder, error) {
	const exitLabel = "❌ WYJDŹ Z PROGRAMU"
	type folderOption struct {
		path string
		f    *wordpress.FbFolder
	}

	var options []folderOption
	var flatten func(fs []wordpress.FbFolder, prefix string)
	flatten = func(fs []wordpress.FbFolder, prefix string) {
		for i := range fs {
			f := &fs[i]
			fullPath := f.Text
			if prefix != "" {
				fullPath = prefix + " > " + f.Text
			}
			options = append(options, folderOption{path: fullPath, f: f})
			if len(f.Children) > 0 {
				flatten(f.Children, fullPath)
			}
		}
	}

	flatten(folders, "")

	if len(options) == 0 {
		return nil, fmt.Errorf("brak dostępnych folderów w WordPress")
	}

	sort.Slice(options, func(i, j int) bool {
		return strings.ToLower(options[i].path) < strings.ToLower(options[j].path)
	})

	optionStrings := make([]string, 0, len(options)+1)
	optionStrings = append(optionStrings, exitLabel)
	for _, o := range options {
		optionStrings = append(optionStrings, o.path)
	}

	var selectedPath string
	prompt := &survey.Select{
		Message:  "Wybierz folder z listy:",
		Options:  optionStrings,
		PageSize: 15,
	}

	err := survey.AskOne(prompt, &selectedPath)
	if err != nil {
		return nil, err
	}

	if selectedPath == exitLabel {
		return nil, nil // Sygnał wyjścia
	}

	for _, o := range options {
		if o.path == selectedPath {
			return o.f, nil
		}
	}

	return nil, fmt.Errorf("nieoczekiwany błąd wyboru")
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
	const exitOpt = "❌ Wyjdź / Powrót"
	const allOpt = "📁 Wszystkie z tego poziomu (lista zbiorcza)"
	const recursiveOpt = "🔄 Rekurencyjnie wszystkie (ten folder + podfoldery)"

	for {
		// Mapujemy opcje na foldery
		optionsMap := make(map[string]*wordpress.FbFolder)
		var options []string

		options = append(options, recursiveOpt)
		options = append(options, allOpt)

		// Dodaj podfoldery do listy opcji
		var subfolderNames []string
		for i := range root.Children {
			sub := &root.Children[i]
			label := "📁 " + sub.Text
			if len(sub.Children) > 0 {
				label = "📂 " + sub.Text + " (zawiera podfoldery)"
			}
			optionsMap[label] = sub
			subfolderNames = append(subfolderNames, label)
		}
		sort.Strings(subfolderNames)
		options = append(options, subfolderNames...)
		options = append(options, exitOpt)

		var selected string
		prompt := &survey.Select{
			Message:  fmt.Sprintf("Folder: %s. Wybierz opcję:", root.Text),
			Options:  options,
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

		if selected == recursiveOpt {
			printFolderLinks(client, root, true)
		} else if selected == allOpt {
			printAllFolders(client, root.Children, false)
		} else {
			subFolder := optionsMap[selected]
			if len(subFolder.Children) > 0 {
				// Jeśli ma podfoldery, zapytaj czy wejść czy wypisać
				var action string
				actionPrompt := &survey.Select{
					Message: fmt.Sprintf("Folder \"%s\" ma podfoldery. Co chcesz zrobić?", subFolder.Text),
					Options: []string{"Otwórz folder", "Wypisz linki bezpośrednie", "Wypisz wszystko rekurencyjnie", "Anuluj"},
				}
				survey.AskOne(actionPrompt, &action)

				switch action {
				case "Otwórz folder":
					showTUI(client, subFolder)
				case "Wypisz linki bezpośrednie":
					printFolderLinks(client, subFolder, false)
				case "Wypisz wszystko rekurencyjnie":
					printFolderLinks(client, subFolder, true)
				}
			} else {
				printFolderLinks(client, subFolder, false)
			}
		}

		fmt.Println("\n---")
	}
}

func printFolderLinks(client *wordpress.Client, folder *wordpress.FbFolder, recursive bool) {
	if recursive {
		fmt.Printf("\n🔄 Pobieranie linków REKURENCYJNIE dla: %s...\n", folder.Text)
		printFolderLinksRecursive(client, folder)
		return
	}

	fmt.Printf("\n🔍 Pobieranie linków dla: %s...\n", folder.Text)
	attachmentIDs, err := client.GetFileBirdAttachmentIDs(folder.ID)
	if err != nil {
		fmt.Printf("  Błąd: %v\n", err)
		return
	}

	if len(attachmentIDs) == 0 {
		fmt.Println("  (Brak bezpośrednich zdjęć w tym folderze)")
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

func printFolderLinksRecursive(client *wordpress.Client, folder *wordpress.FbFolder) {
	// Wypisz zdjęcia z bieżącego folderu
	attachmentIDs, err := client.GetFileBirdAttachmentIDs(folder.ID)
	if err == nil && len(attachmentIDs) > 0 {
		fmt.Printf("\n--- Folder: %s ---\n", folder.Text)
		for _, aid := range attachmentIDs {
			url, err := client.GetMediaSourceURL(aid)
			if err == nil {
				fmt.Println(url)
			}
		}
	}

	// Wypisz zdjęcia z podfolderów
	for i := range folder.Children {
		printFolderLinksRecursive(client, &folder.Children[i])
	}
}

func printAllFolders(client *wordpress.Client, subs []wordpress.FbFolder, recursive bool) {
	fmt.Println("\n📦 Generowanie listy linków...")
	for i := range subs {
		printFolderLinks(client, &subs[i], recursive)
	}
}
