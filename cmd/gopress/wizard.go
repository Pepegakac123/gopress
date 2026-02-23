package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Pepegakac123/gopress/internal/version"
)

func runWizard() {
	// Czyszczenie starych plików po aktualizacji
	version.CleanupOldBinary()

	fmt.Println("Tryb interaktywny: Nie podano flag, więc zadam kilka pytań...")

	// Sprawdzenie aktualizacji (Synchronicznie, aby uniknąć problemów z stdin)
	fmt.Print("Sprawdzanie aktualizacji... ")
	newerReleases, err := version.CheckForUpdates("Pepegakac123/gopress")
	if err != nil {
		fmt.Printf("\nBłąd podczas sprawdzania aktualizacji: %v\n", err)
	}

	if len(newerReleases) > 0 {
		latestRelease := newerReleases[len(newerReleases)-1]
		fmt.Printf("\nDostępna jest nowa wersja: %s (Twoja to %s)\n", latestRelease.GetTagName(), version.CurrentVersion)

		var cumulativeChangelog strings.Builder
		for _, release := range newerReleases {
			if release.GetBody() != "" {
				cumulativeChangelog.WriteString(fmt.Sprintf("\n## Zmiany w wersji %s:\n", release.GetTagName()))
				cumulativeChangelog.WriteString(strings.TrimSpace(release.GetBody()) + "\n")
			}
		}

		if cumulativeChangelog.Len() > 0 {
			fmt.Println("\n---")
			fmt.Println(cumulativeChangelog.String())
			fmt.Println("---")
		}

		var update bool
		prompt := &survey.Confirm{
			Message: "Czy chcesz pobrać i zaktualizować automatycznie?",
			Default: true,
		}
		if err := survey.AskOne(prompt, &update); err == nil && update {
			fmt.Println("Pobieranie i instalowanie aktualizacji... (To może chwilę potrwać)")
			if err := version.PerformUpdate(latestRelease); err != nil {
				fmt.Printf("Błąd aktualizacji: %v\n", err)
			} else {
				fmt.Println("Sukces! Zaktualizowano pomyślnie.")
				restartApplication()
			}
		}
	} else if err == nil {
		fmt.Println("(Aktualna)")
	}

	handleSurveyErr := func(err error) {
		if err == nil {
			return
		}
		if err == terminal.InterruptErr {
			fmt.Println("\nPrzerwano przez użytkownika (Ctrl+C). Do widzenia!")
			os.Exit(0)
		}
		fmt.Printf("\nBłąd ankiety: %v\n", err)
		os.Exit(1)
	}
	appConfig.DeleteOriginals = false

	inputPrompt := &survey.Input{
		Message: "Gdzie są zdjęcia (folder wejściowy)?",
		Default: "./raw",
		Suggest: suggestPaths,
	}
	err = survey.AskOne(inputPrompt, &appConfig.InputDir, survey.WithValidator(survey.Required))
	handleSurveyErr(err)

	appConfig.InputDir = sanitizePath(appConfig.InputDir)

	// Obliczamy domyślny output
	var defaultOut string
	if strings.ToLower(filepath.Ext(appConfig.InputDir)) == ".zip" {
		defaultOut = filepath.Join(filepath.Dir(appConfig.InputDir), "webp")
	} else {
		defaultOut = filepath.Join(appConfig.InputDir, "webp")
	}

	// Pytanie 2: Output
	outputPrompt := &survey.Input{
		Message: fmt.Sprintf("Gdzie zapisać wyniki? Zostaw puste, aby użyć domyślnej lokalizacji: %s", defaultOut),
		Default: defaultOut,
		Suggest: suggestPaths,
	}
	err = survey.AskOne(outputPrompt, &appConfig.OutputDir)
	handleSurveyErr(err)

	// Pytanie 3: Jakość
	err = survey.AskOne(&survey.Input{
		Message: "Jakość obrazu WebP (0-100):",
		Default: "80",
	}, &appConfig.Quality, survey.WithValidator(validateRange(0, 100)))
	handleSurveyErr(err)

	// Pytanie 4: Szerokość
	err = survey.AskOne(&survey.Input{
		Message: "Maksymalna szerokość (px):",
		Default: "2560",
	}, &appConfig.MaxWidth, survey.WithValidator(validateRange(10, 10000)))
	handleSurveyErr(err)

	// Pytanie 5: Upload
	err = survey.AskOne(&survey.Confirm{
		Message: "Czy chcesz wysłać pliki do WordPressa?",
		Default: false,
	}, &appConfig.Upload)
	handleSurveyErr(err)

	if appConfig.Upload {
		err = survey.AskOne(&survey.Input{
			Message: "Podaj domenę WP (z https://):",
		}, &appConfig.WpDomain, survey.WithValidator(survey.Required))
		handleSurveyErr(err)

		err = survey.AskOne(&survey.Input{
			Message: "Użytkownik WP:",
			Default: "overflow",
		}, &appConfig.WpUser, survey.WithValidator(survey.Required))
		handleSurveyErr(err)

		err = survey.AskOne(&survey.Password{
			Message: "Hasło Aplikacji (Application Password): ",
		}, &appConfig.WpPassword, survey.WithValidator(survey.Required))
		handleSurveyErr(err)

		err = survey.AskOne(&survey.Password{
			Message: "Token API FileBird (FileBird -> Narzędzia -> Wygeneruj API) - Jeśli nie korzystasz z FileBird, zostaw puste.",
		}, &appConfig.FileBirdToken)
	}
	// Pytania 6: Usuwanie oryginalnych plików
	err = survey.AskOne(&survey.Confirm{
		Message: "Czy usunąć oryginalne pliki po konwersji?",
		Help:    "Oryginały zostaną bezpowrotnie usunięte z dysku. Zostaną tylko pliki WebP.",
		Default: false,
	}, &appConfig.DeleteOriginals)
	handleSurveyErr(err)
}

func validateRange(min, max int) survey.Validator {
	return func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("nieprawidłowy typ danych")
		}

		num, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("to nie jest liczba")
		}

		if num < min || num > max {
			return fmt.Errorf("wartość musi być pomiędzy %d a %d", min, max)
		}
		return nil
	}
}

func suggestPaths(toComplete string) []string {
	files, _ := filepath.Glob(toComplete + "*")
	return files
}
