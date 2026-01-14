package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// CurrentVersion to aktualna wersja aplikacji.
// Wartość powinna być nadpisywana przez ldflags podczas budowania (np. -X ...CurrentVersion=v1.2.1).
var CurrentVersion = "v0.0.0"

// EnsureBinaryName sprawdza, czy plik wykonywalny ma oczekiwaną nazwę.
// Jeśli nie, zmienia ją (np. z gopress_darwin_arm64 na gopress).
func EnsureBinaryName(desiredName string) error {
	// Pomiń w trybie deweloperskim
	if CurrentVersion == "v0.0.0" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exe)
	filename := filepath.Base(exe)

	// Obsługa .exe dla Windows
	targetName := desiredName
	if runtime.GOOS == "windows" {
		if !strings.EqualFold(filepath.Ext(targetName), ".exe") {
			targetName += ".exe"
		}
		if strings.EqualFold(filename, targetName) {
			return nil
		}
	} else {
		if filename == targetName {
			return nil
		}
	}

	targetPath := filepath.Join(dir, targetName)

	// Jeśli plik docelowy istnieje, usuń go (zastąpienie)
	if _, err := os.Stat(targetPath); err == nil {
		// Windows nie pozwala na nadpisanie rename, trzeba usunąć
		if err := os.Remove(targetPath); err != nil {
			return fmt.Errorf("nie można usunąć istniejącego pliku %s: %w", targetName, err)
		}
	}

	if err := os.Rename(exe, targetPath); err != nil {
		return fmt.Errorf("nie udało się zmienić nazwy pliku na %s: %w", targetName, err)
	}

	fmt.Printf("ℹ️ Program zmienił nazwę na: %s\n", targetName)
	return nil
}

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
		// Jeśli obecna wersja jest niepoprawna, zakładamy, że chcemy zaktualizować (dev mode)
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
