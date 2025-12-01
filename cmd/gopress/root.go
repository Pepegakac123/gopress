package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Pepegakac123/gopress/internal/processor"
	"github.com/Pepegakac123/gopress/internal/scanner"
	"github.com/spf13/cobra"
)

type Config struct {
	InputDir  string
	OutputDir string
}

var appConfig Config

var rootCmd = &cobra.Command{
	Use:   "gopress",
	Short: "A tool for optimalizationa and publishing images to the wordpress",
	Long: `GoPress is a CLI tool written in Golang. It allows user to convert large number of variety of images type to the webp format with optimalization options that make them
	efficient for web usage. The tool provides a simple and intuitive interface for users to easily convert their images to the webp format, while also providing advanced options for fine-tuning the conversion process.`,
	Run: func(cmd *cobra.Command, args []string) {
		if appConfig.InputDir == "" || appConfig.OutputDir == "" {
			runWizard()
		} else {
			fmt.Println("Silent mode")
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
		finalSize := processor.RunSimpleAsync(files, appConfig.OutputDir)

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
	},
}

func init() {
	rootCmd.Flags().StringVarP(&appConfig.InputDir, "input", "i", "", "Ścieżka do folderu z obrazami")
	rootCmd.Flags().StringVarP(&appConfig.OutputDir, "output", "o", "", "Ścieżka gdzie zapisać wyniki")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runWizard() {
	fmt.Println("Tryb interaktywny: Nie podano flag, więc zadam kilka pytań...")
	qs := []*survey.Question{
		{
			Name: "InputDir",
			Prompt: &survey.Input{
				Message: "Gdzie są zdjęcia (folder wejściowy)?",
				Default: "./raw",
			},
			Validate: survey.Required,
		},
		{
			Name: "OutputDir",
			Prompt: &survey.Input{
				Message: "Gdzie zapisać wyniki?",
				Default: "./out",
			},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(qs, &appConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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
