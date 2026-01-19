package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Pepegakac123/gopress/internal/processor"
	"github.com/Pepegakac123/gopress/internal/scanner"
	"github.com/Pepegakac123/gopress/internal/uploader"
	"github.com/Pepegakac123/gopress/internal/wordpress"
	"github.com/spf13/cobra"
)

func runGopress(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if appConfig.InputDir == "" && len(args) > 0 {
		appConfig.InputDir = args[0]
	}

	if appConfig.InputDir == "" {
		runWizard()
	} else {
		fmt.Println("Tryb cichy: Używam ustawień startowych.")
	}

	appConfig.InputDir = sanitizePath(appConfig.InputDir)
	if appConfig.OutputDir != "" {
		appConfig.OutputDir = sanitizePath(appConfig.OutputDir)
	}

	originalInput := appConfig.InputDir
	stat, err := os.Stat(originalInput)
	if err != nil {
		log.Fatalf("Błąd dostępu do %s: %v", originalInput, err)
	}
	isZip := !stat.IsDir() && strings.ToLower(filepath.Ext(originalInput)) == ".zip"

	if appConfig.OutputDir == "" {
		if isZip {
			appConfig.OutputDir = filepath.Join(filepath.Dir(originalInput), "webp")
		} else {
			appConfig.OutputDir = filepath.Join(appConfig.InputDir, "webp")
		}
	}

	var shouldCleanup = true
	if isZip {
		tempDir, err := os.MkdirTemp("", "gopress_unzip_*")
		if err != nil {
			log.Fatalf("Błąd tworzenia katalogu tymczasowego: %v", err)
		}
		defer func() {
			if shouldCleanup {
				os.RemoveAll(tempDir)
			}
		}()

		fmt.Printf("Rozpakowywanie %s do %s...\n", originalInput, tempDir)
		if err := unzip(originalInput, tempDir); err != nil {
			log.Fatalf("Błąd rozpakowywania zip: %v", err)
		}
		appConfig.InputDir = tempDir
	}

	var wpClient *wordpress.Client
	if appConfig.Upload {
		if appConfig.WpDomain == "" || appConfig.WpUser == "" || appConfig.WpPassword == "" {
			log.Fatal("Błąd: Tryb --upload wymaga podania --wp-domain, --wp-user i --wp-secret. Uruchom program bez parametrów, aby włączyć kreatora.")
		}
		fmt.Println("\nŁączenie z WordPress...")
		wpClient = wordpress.NewClient(appConfig.WpDomain, appConfig.WpUser, appConfig.WpPassword, appConfig.FileBirdToken)
		if err := wpClient.CheckConnection(); err != nil {
			log.Fatalf("Błąd połączenia z WP: %v", err)
		}
		fmt.Println("Połączono z WordPress (Autoryzacja OK)")
	}

	fmt.Printf("Skanowanie folderu: %s\n", appConfig.InputDir)

	files, initialSize, err := scanner.FindImages(appConfig.InputDir)
	if err != nil {
		log.Fatalf("Bląd podczas skanowania %v", err)
	}
	if len(files) == 0 {
		log.Fatal("Nie znaleziono plików")
		return
	}

	fmt.Printf("Znaleziono %d obrazów do przetworzenia.\n", len(files))
	fmt.Printf("Parametry: Jakość %d%%, Max Szerokość %dpx\n", appConfig.Quality, appConfig.MaxWidth)
	if _, err := os.Stat(appConfig.OutputDir); os.IsNotExist(err) {
		os.MkdirAll(appConfig.OutputDir, 0755)
	}
	finalSize, convertedFiles := processor.RunWorkerPool(ctx, files, appConfig.InputDir, appConfig.OutputDir, appConfig.Quality, appConfig.MaxWidth, appConfig.DeleteOriginals)

	var savings float64
	if initialSize > 0 {
		savings = (1.0 - float64(finalSize)/float64(initialSize)) * 100
	}
	fmt.Println("\n--- 📊 Podsumowanie ---")
	fmt.Printf("✅ Przetworzono obrazów: %d\n", len(files))
	fmt.Printf("📦 Rozmiar przed:       %s\n", formatBytes(initialSize))
	fmt.Printf("💾 Rozmiar po:          %s\n", formatBytes(finalSize))
	fmt.Printf("📉 Oszczędność:         %.2f%%\n", savings)
	fmt.Printf("📂 Folder wynikowy:     %s\n", appConfig.OutputDir)
	if appConfig.Upload && len(convertedFiles) > 0 {
		prepareFileBirdToken(wpClient)
		useFileBird := appConfig.FileBirdToken != ""
		uploader.Run(ctx, wpClient, convertedFiles, appConfig.OutputDir, useFileBird, 0)
	}

	var openResult bool
	prompt := &survey.Confirm{
		Message: "Czy otworzyć folder z wynikami?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &openResult); err == nil && openResult {
		if isZip {
			shouldCleanup = false
			fmt.Println("Tymczasowe pliki (rozpakowany ZIP) zostały zachowane, ponieważ otwierasz folder.")
		}
		openFolder(appConfig.OutputDir)
	}
}

func prepareFileBirdToken(client *wordpress.Client) {
	if appConfig.FileBirdToken == "" {
		return
	}

	fmt.Print("Weryfikacja tokenu FileBird... ")
	if err := client.CheckFileBirdConnection(); err != nil {
		fmt.Printf("\nBŁĄD weryfikacji tokenu: %v\n", err)

		var continueWithoutFolders bool
		prompt := &survey.Confirm{
			Message: "Token FileBird jest nieprawidłowy. Czy chcesz kontynuować upload BEZ obsługi folderów (płasko)?",
			Default: false,
		}

		if err := survey.AskOne(prompt, &continueWithoutFolders); err != nil {
			fmt.Println("\nOperacja anulowana.")
			os.Exit(0)
		}

		if !continueWithoutFolders {
			fmt.Println("Anulowano. Popraw token i spróbuj ponownie.")
			os.Exit(0)
		}

		fmt.Println("Zrozumiałem. Kontynuuję upload w trybie płaskim.")
		appConfig.FileBirdToken = ""
	} else {
		fmt.Println("OK")
	}
}

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
