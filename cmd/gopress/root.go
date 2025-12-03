package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Pepegakac123/gopress/internal/processor"
	"github.com/Pepegakac123/gopress/internal/scanner"
	"github.com/Pepegakac123/gopress/internal/uploader"
	"github.com/Pepegakac123/gopress/internal/wordpress"
	"github.com/spf13/cobra"
)

type Config struct {
	InputDir        string
	OutputDir       string
	Upload          bool
	WpDomain        string
	WpUser          string
	WpPassword      string
	FileBirdToken   string
	Quality         int
	MaxWidth        int
	DeleteOriginals bool
}

var appConfig Config

var rootCmd = &cobra.Command{
	Use:   "gopress [folder-ze-zdjeciami]", // Bardziej jasne niż [input-dir]
	Short: "Automat do zmniejszania zdjęć i wysyłania na WordPressa",
	Long: `GoPress to Twój asystent do zadań specjalnych.
	Bierze cały folder zdjęć (JPG, PNG, a nawet HEIC z iPhone'a), automatycznie przerabia je na szybki format WebP, zmniejsza do odpowiedniego rozmiaru i wysyła na stronę internetową.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if appConfig.Quality < 0 || appConfig.Quality > 100 {
			return fmt.Errorf("nieprawidłowa jakość (%d). Podaj wartość między 0 a 100", appConfig.Quality)
		}
		if appConfig.MaxWidth <= 10 {
			return fmt.Errorf("szerokość (%d) jest zbyt mała. Podaj wartość większą niż 10", appConfig.MaxWidth)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if appConfig.InputDir == "" && len(args) > 0 {
			appConfig.InputDir = args[0]
		}
		if appConfig.InputDir == "" {
			runWizard()
		} else {
			if appConfig.OutputDir == "" {
				appConfig.OutputDir = filepath.Join(appConfig.InputDir, "webp")
			}
			fmt.Println("Tryb cichy: Używam ustawień startowych.")
		}

		var wpClient *wordpress.Client
		if appConfig.Upload {
			if appConfig.WpDomain == "" || appConfig.WpUser == "" || appConfig.WpPassword == "" {
				log.Fatal("❌ Błąd: Tryb --upload wymaga podania --wp-domain, --wp-user i --wp-secret. Uruchom program bez parametrów, aby włączyć kreatora.")
			}
			fmt.Println("\n Łączenie z WordPress...")
			wpClient = wordpress.NewClient(appConfig.WpDomain, appConfig.WpUser, appConfig.WpPassword, appConfig.FileBirdToken)
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
		fmt.Printf("⚙️  Parametry: Jakość %d%%, Max Szerokość %dpx\n", appConfig.Quality, appConfig.MaxWidth)
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
	},
}

func init() {
	// Input/Output
	rootCmd.Flags().StringVarP(&appConfig.InputDir, "input", "i", "", "Ścieżka do folderu ze zdjęciami (możesz też przeciągnąć folder na okno)")
	rootCmd.Flags().StringVarP(&appConfig.OutputDir, "output", "o", "", "Gdzie zapisać gotowe pliki (domyślnie tworzy folder 'webp' w środku)")

	// Upload
	rootCmd.Flags().BoolVar(&appConfig.Upload, "upload", false, "Wyślij gotowe pliki na serwer WordPress")

	// WP Config
	rootCmd.Flags().StringVar(&appConfig.WpDomain, "wp-domain", "", "Adres strony (np. https://mojastrona.pl)")
	rootCmd.Flags().StringVar(&appConfig.WpUser, "wp-user", "", "Twój login do WordPressa")
	// hasło WP
	rootCmd.Flags().StringVar(&appConfig.WpPassword, "wp-secret", "", "Hasło Aplikacji (NIE twoje hasło do logowania!). Wygeneruj w: Użytkownicy -> Profil")
	// Jakość
	rootCmd.Flags().IntVarP(&appConfig.Quality, "quality", "q", 80, "Jakość obrazu (0-100). 80 to złoty środek.")
	// Wymiary
	rootCmd.Flags().IntVarP(&appConfig.MaxWidth, "width", "w", 2560, "Maksymalna szerokość w px (program pomniejszy duże zdjęcia, ale nie powiększy małych)")
	// Delete - Zostawmy to mocne ostrzeżenie
	rootCmd.Flags().BoolVarP(&appConfig.DeleteOriginals, "delete", "d", false, "USUŃ oryginały po konwersji (Ostrożnie! Tej operacji nie da się cofnąć)")
	// FileBird
	rootCmd.Flags().StringVar(&appConfig.FileBirdToken, "fb-token", "", "Token FileBird (jeśli chcesz zachować strukturę folderów) i strona używa wtyczki FileBird")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runWizard() {
	fmt.Println("Tryb interaktywny: Nie podano flag, więc zadam kilka pytań...")

	handleSurveyErr := func(err error) {
		if err == nil {
			return
		}
		if err == terminal.InterruptErr {
			fmt.Println("\n🛑 Przerwano przez użytkownika (Ctrl+C). Do widzenia!")
			os.Exit(0)
		}
		fmt.Printf("\n❌ Błąd ankiety: %v\n", err)
		os.Exit(1)
	}

	inputPrompt := &survey.Input{
		Message: "Gdzie są zdjęcia (folder wejściowy)?",
		Default: "./raw",
	}
	err := survey.AskOne(inputPrompt, &appConfig.InputDir, survey.WithValidator(survey.Required))
	handleSurveyErr(err)

	// Obliczamy domyślny output
	defaultOut := filepath.Join(appConfig.InputDir, "webp")

	// Pytanie 2: Output
	outputPrompt := &survey.Input{
		Message: fmt.Sprintf("Gdzie zapisać wyniki? Zostaw puste, aby użyć domyślnej lokalizacji: %s", defaultOut),
		Default: defaultOut,
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
			Default: "admin",
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
		Message: "⚠️  Czy usunąć oryginalne pliki po konwersji?",
		Help:    "Oryginały zostaną bezpowrotnie usunięte z dysku. Zostaną tylko pliki WebP.",
		Default: false,
	}, &appConfig.DeleteOriginals)
	handleSurveyErr(err)
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
func prepareFileBirdToken(client *wordpress.Client) {
	if appConfig.FileBirdToken == "" {
		return
	}

	fmt.Print("📂 Weryfikacja tokenu FileBird... ")
	if err := client.CheckFileBirdConnection(); err != nil {
		fmt.Printf("\n❌ BŁĄD weryfikacji tokenu: %v\n", err)

		var continueWithoutFolders bool
		prompt := &survey.Confirm{
			Message: "Token FileBird jest nieprawidłowy. Czy chcesz kontynuować upload BEZ obsługi folderów (płasko)?",
			Default: false,
		}

		if err := survey.AskOne(prompt, &continueWithoutFolders); err != nil {
			fmt.Println("\n🛑 Operacja anulowana.")
			os.Exit(0)
		}

		if !continueWithoutFolders {
			fmt.Println("🛑 Anulowano. Popraw token i spróbuj ponownie.")
			os.Exit(0)
		}

		fmt.Println("⚠️  Zrozumiałem. Kontynuuję upload w trybie płaskim.")
		appConfig.FileBirdToken = ""
	} else {
		fmt.Println("✅ OK")
	}
}
