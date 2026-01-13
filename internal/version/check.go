package version

import (
	"fmt"
	"os"
	"runtime"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// CurrentVersion to aktualna wersja aplikacji.
// Wartość powinna być nadpisywana przez ldflags podczas budowania (np. -X ...CurrentVersion=v1.2.1).
var CurrentVersion = "v0.0.0"

// CheckUpdate sprawdza, czy dostępna jest nowsza wersja w podanym repozytorium (slug: "owner/repo").
// Zwraca znaleziony release, flagę found (czy znaleziono cokolwiek) oraz ewentualny błąd.
func CheckUpdate(slug string) (*selfupdate.Release, bool, error) {
	fmt.Println("DEBUG: Rozpoczynam wykrywanie wersji...", slug)
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return nil, false, fmt.Errorf("błąd sprawdzania aktualizacji: %w", err)
	}
	if !found {
		fmt.Println("DEBUG: Nie znaleziono żadnej wersji na GitHub.")
		return nil, false, nil
	}

	fmt.Printf("DEBUG: Znaleziono najnowszą wersję na GitHub: %s\n", latest.Version.String())
	fmt.Printf("DEBUG: Obecna wersja lokalna: %s\n", CurrentVersion)

	// Parsowanie obecnej wersji
	vCurrent, err := semver.ParseTolerant(CurrentVersion)
	if err != nil {
		fmt.Printf("DEBUG: Błąd parsowania obecnej wersji (%s): %v. Wymuszam aktualizację.\n", CurrentVersion, err)
		return latest, true, nil
	}

	fmt.Printf("DEBUG: Porównanie: Czy %s > %s? Wynik: %v\n", latest.Version.String(), vCurrent.String(), latest.Version.GT(vCurrent))

	// Jeśli najnowsza wersja jest nowsza niż obecna
	if latest.Version.GT(vCurrent) {
		return latest, true, nil
	}

	return nil, false, nil
}

// PerformUpdate pobiera i instaluje nową wersję.
func PerformUpdate(slug string) error {
	// Ponowne pobranie najnowszej wersji, aby przekazać ją do UpdateTo
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("nie znaleziono wydania do aktualizacji")
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("nie udało się ustalić ścieżki pliku wykonywalnego: %w", err)
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return err
	}

	return nil
}

// CleanupOldBinary usuwa pozostałości po aktualizacji (pliki .old na Windows)
func CleanupOldBinary() {
	if runtime.GOOS != "windows" {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	oldExe := exe + ".old"
	// Próbujemy usunąć, ignorujemy błędy (np. jeśli plik nie istnieje lub jest zablokowany)
	_ = os.Remove(oldExe)
}
