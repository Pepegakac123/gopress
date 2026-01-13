package version

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// CurrentVersion to aktualna wersja aplikacji.
// Wartość powinna być nadpisywana przez ldflags podczas budowania (np. -X ...CurrentVersion=v1.2.1).
var CurrentVersion = "v0.0.0"

// CheckUpdate sprawdza, czy dostępna jest nowsza wersja w podanym repozytorium (slug: "owner/repo").
// Zwraca znaleziony release, flagę found (czy znaleziono cokolwiek) oraz ewentualny błąd.
func CheckUpdate(slug string) (*selfupdate.Release, bool, error) {
	latest, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return nil, false, fmt.Errorf("błąd sprawdzania aktualizacji: %w", err)
	}
	if !found {
		return nil, false, nil
	}

	// Parsowanie obecnej wersji
	vCurrent, err := semver.ParseTolerant(CurrentVersion)
	if err != nil {
		// Jeśli obecna wersja jest "v0.0.0" lub niepoprawna, zakładamy, że chcemy zaktualizować (dev mode)
		// Ale dla bezpieczeństwa w produkcji:
		return latest, true, nil
	}

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
