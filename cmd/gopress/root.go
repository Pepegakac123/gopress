package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
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
		fmt.Printf("Config -> Input: %s | Output: %s\n", appConfig.InputDir, appConfig.OutputDir)
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
