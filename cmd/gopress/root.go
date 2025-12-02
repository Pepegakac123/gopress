package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Pepegakac123/gopress/internal/processor"
	"github.com/Pepegakac123/gopress/internal/scanner"
	"github.com/Pepegakac123/gopress/internal/wordpress"
	"github.com/spf13/cobra"
)

type Config struct {
	InputDir   string
	OutputDir  string
	Upload     bool
	WpDomain   string
	WpUser     string
	WpPassword string
}

var appConfig Config

var rootCmd = &cobra.Command{
	Use:   "gopress [input-dir]",
	Short: "A tool for optimalizationa and publishing images to the wordpress",
	Long: `GoPress is a CLI tool written in Golang. It allows user to convert large number of variety of images type to the webp format with optimalization options that make them
	efficient for web usage. The tool provides a simple and intuitive interface for users to easily convert their images to the webp format, while also providing advanced options for fine-tuning the conversion process.`,
	Run: func(cmd *cobra.Command, args []string) {
		if appConfig.InputDir == "" && len(args) > 0 {
			appConfig.InputDir = args[0]
		}
		if appConfig.InputDir == "" {
			runWizard()
		} else {
			if appConfig.OutputDir == "" {
				appConfig.OutputDir = filepath.Join(appConfig.InputDir, "webp")
			}
			fmt.Println("Silent mode")
		}

		var wpClient *wordpress.Client
		if appConfig.Upload {
			if appConfig.WpDomain == "" || appConfig.WpUser == "" || appConfig.WpPassword == "" {
				log.Fatal("❌ Błąd: Tryb --upload wymaga podania --wp-domain, --wp-user i --wp-secret")
			}
			fmt.Println("\n Łączenie z WordPress...")
			wpClient = wordpress.NewClient(appConfig.WpDomain, appConfig.WpUser, appConfig.WpPassword)
			if err := wpClient.CheckConnection(); err != nil {
				log.Fatalf("Błąd połączenia z WP: %v", err)
			}
			fmt.Println("✅ Połączono z WordPress (Autoryzacja OK)")
		}

		fmt.Printf("🔍 Skanowanie folderu: %s\n", appConfig.InputDir)

		files, initialSize, err := scanner.FindImages(appConfig.InputDir)
		if err != nil {
			log.Fatalf("Bląd podczas skanowania %v", err)
		}
		if len(files) == 0 {
			log.Fatal("⚠️ Nie znaleziono plików")
			return
		}
		fmt.Printf("✅ Znaleziono %d obrazów do przetworzenia.\n", len(files))
		if _, err := os.Stat(appConfig.OutputDir); os.IsNotExist(err) {
			os.MkdirAll(appConfig.OutputDir, 0755)
		}
		finalSize := processor.RunWorkerPool(files, appConfig.OutputDir)

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
		if appConfig.Upload {
			fmt.Println("🚀 (Tu wkrótce nastąpi wysyłanie plików do WP...)")
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&appConfig.InputDir, "input", "i", "", "Ścieżka do folderu z obrazami")
	rootCmd.Flags().StringVarP(&appConfig.OutputDir, "output", "o", "", "Ścieżka gdzie zapisać wyniki")
	rootCmd.Flags().BoolVar(&appConfig.Upload, "upload", false, "Włącz wysyłanie na WP")
	rootCmd.Flags().StringVar(&appConfig.WpDomain, "wp-domain", "", "Domena WP (np. https://mojastrona.pl)")
	rootCmd.Flags().StringVar(&appConfig.WpUser, "wp-user", "", "Użytkownik WP")
	rootCmd.Flags().StringVar(&appConfig.WpPassword, "wp-secret", "", "Hasło Aplikacji WP w formacie XXXX XXXX XXXX XXXX XXXX XXXX")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runWizard() {
	fmt.Println("Tryb interaktywny: Nie podano flag, więc zadam kilka pytań...")
	inputPrompt := &survey.Input{
		Message: "Gdzie są zdjęcia (folder wejściowy)?",
		Default: "./raw",
	}
	err := survey.AskOne(inputPrompt, &appConfig.InputDir, survey.WithValidator(survey.Required))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defaultOut := filepath.Join(appConfig.InputDir, "webp")
	outputPrompt := &survey.Input{
		Message: fmt.Sprintf("Gdzie zapisać wyniki? Zostaw puste, aby użyć domyślnej lokalizacji: %s", defaultOut),
		Default: defaultOut,
	}
	err = survey.AskOne(outputPrompt, &appConfig.OutputDir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	survey.AskOne(&survey.Confirm{
		Message: "Czy chcesz wysłać pliki do WordPressa?",
		Default: false,
	}, &appConfig.Upload)
	if appConfig.Upload {
		survey.AskOne(&survey.Input{
			Message: "Podaj domenę WP (z https://):",
		}, &appConfig.WpDomain, survey.WithValidator(survey.Required))

		survey.AskOne(&survey.Input{
			Message: "Użytkownik WP:",
			Default: "admin",
		}, &appConfig.WpUser, survey.WithValidator(survey.Required))

		survey.AskOne(&survey.Password{
			Message: "Hasło Aplikacji (Application Password) w formacie XXXX XXXX XXXX XXXX XXXX XXXX:",
		}, &appConfig.WpPassword, survey.WithValidator(survey.Required))
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
