package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Pepegakac123/gopress/internal/version"
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
	Version         bool
	Update          bool
}

var appConfig Config

var rootCmd = &cobra.Command{
	Use:   "gopress [folder-ze-zdjeciami]",
	Short: "Automat do zmniejszania zdjęć i wysyłania na WordPressa",
	Long: `GoPress to Twój asystent do zadań specjalnych.
	Bierze cały folder zdjęć (JPG, PNG, a nawet HEIC z iPhone'a), automatycznie przerabia je na szybki format WebP, zmniejsza do odpowiedniego rozmiaru i wysyła na stronę internetową.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := version.EnsureBinaryName("gopress"); err != nil {
			fmt.Printf("Nie udało się zmienić nazwy programu: %v\n", err)
		}

		if appConfig.Version {
			fmt.Printf("GoPress version %s\n", version.CurrentVersion)
			os.Exit(0)
		}

		if appConfig.Update {
			fmt.Print("Sprawdzanie aktualizacji... ")
			newerReleases, err := version.CheckForUpdates("Pepegakac123/gopress")
			if err != nil {
				fmt.Printf("\nBłąd podczas sprawdzania aktualizacji: %v\n", err)
				os.Exit(1)
			}

			if len(newerReleases) == 0 {
				fmt.Println("Masz już najnowszą wersję.")
				os.Exit(0)
			}

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

			fmt.Println("Pobieranie...")
			if err := version.PerformUpdate(latestRelease); err != nil {
				fmt.Printf("Błąd aktualizacji: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Sukces! Zaktualizowano pomyślnie.")
			restartApplication()
		}

		if appConfig.Quality < 0 || appConfig.Quality > 100 {
			return fmt.Errorf("nieprawidłowa jakość (%d). Podaj wartość między 0 a 100", appConfig.Quality)
		}
		if appConfig.MaxWidth <= 10 {
			return fmt.Errorf("szerokość (%d) jest zbyt mała. Podaj wartość większą niż 10", appConfig.MaxWidth)
		}

		return nil
	},
	Run: runGopress,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	// Input/Output
	rootCmd.Flags().StringVarP(&appConfig.InputDir, "input", "i", "", "Ścieżka do folderu ze zdjęciami (możesz też przeciągnąć folder na okno)")
	rootCmd.Flags().StringVarP(&appConfig.OutputDir, "output", "o", "", "Gdzie zapisać gotowe pliki (domyślnie tworzy folder 'webp' w środku)")

	_ = rootCmd.MarkFlagFilename("input", "zip", "jpg", "jpeg", "png", "heic")
	_ = rootCmd.MarkFlagDirname("output")

	// Upload
	rootCmd.Flags().BoolVar(&appConfig.Upload, "upload", false, "Wyślij gotowe pliki na serwer WordPress")

	// WP Config
	rootCmd.PersistentFlags().StringVar(&appConfig.WpDomain, "wp-domain", "", "Adres strony (np. https://mojastrona.pl)")
	rootCmd.PersistentFlags().StringVar(&appConfig.WpUser, "wp-user", "", "Twój login do WordPressa")
	rootCmd.PersistentFlags().StringVar(&appConfig.WpPassword, "wp-secret", "", "Hasło Aplikacji (NIE twoje hasło do logowania!). Wygeneruj w: Użytkownicy -> Profil")
	rootCmd.PersistentFlags().StringVar(&appConfig.FileBirdToken, "fb-token", "", "Token FileBird (jeśli chcesz zachować strukturę folderów) i strona używa wtyczki FileBird")

	rootCmd.Flags().IntVarP(&appConfig.Quality, "quality", "q", 80, "Jakość obrazu (0-100). 80 to złoty środek.")
	rootCmd.Flags().IntVarP(&appConfig.MaxWidth, "width", "w", 2560, "Maksymalna szerokość w px (program pomniejszy duże zdjęcia, ale nie powiększy małych)")
	rootCmd.Flags().BoolVarP(&appConfig.DeleteOriginals, "delete", "d", false, "USUŃ oryginały po konwersji (Ostrożnie! Tej operacji nie da się cofnąć)")

	// Globalne info/aktualizacja
	rootCmd.Flags().BoolVarP(&appConfig.Version, "version", "v", false, "Wyświetl wersję programu")
	rootCmd.Flags().BoolVar(&appConfig.Update, "update", false, "Sprawdź i zainstaluj aktualizację")

	rootCmd.AddCommand(linksCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
