package version

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/Pepegakac123/gopress/internal/github"
	"github.com/blang/semver"
	gh "github.com/google/go-github/v58/github"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var CurrentVersion = "v0.0.0"

func EnsureBinaryName(desiredName string) error {
	if CurrentVersion == "v0.0.0" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Dir(exe)
	filename := filepath.Base(exe)

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

	if _, err := os.Stat(targetPath); err == nil {
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

// CheckForUpdates sprawdza, czy dostępne są nowsze wersje w podanym repozytorium.
// Zwraca posortowaną listę nowszych wydań.
func CheckForUpdates(slug string) ([]*gh.RepositoryRelease, error) {
	parts := strings.Split(slug, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("nieprawidłowy format repozytorium: %s", slug)
	}
	owner, repo := parts[0], parts[1]

	client := github.NewClient()
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("błąd pobierania wydań z GitHub: %w", err)
	}

	vCurrent, err := semver.ParseTolerant(CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("nieprawidłowy format obecnej wersji '%s': %w", CurrentVersion, err)
	}

	var newerReleases []*gh.RepositoryRelease
	for _, release := range releases {
		vRelease, err := semver.ParseTolerant(release.GetTagName())
		if err != nil {
			continue // Pomiń tagi, które nie są poprawnymi wersjami
		}

		if vRelease.GT(vCurrent) {
			newerReleases = append(newerReleases, release)
		}
	}

	// Sortuj od najstarszej do najnowszej
	sort.Slice(newerReleases, func(i, j int) bool {
		vI, _ := semver.ParseTolerant(newerReleases[i].GetTagName())
		vJ, _ := semver.ParseTolerant(newerReleases[j].GetTagName())
		return vI.LT(vJ)
	})

	return newerReleases, nil
}

// PerformUpdate pobiera i instaluje podaną wersję.
func PerformUpdate(release *gh.RepositoryRelease) error {
	if release == nil {
		return fmt.Errorf("brak wydania do zainstalowania")
	}

	assetURL := ""
	for _, asset := range release.Assets {
		// Prosta logika dopasowania - można ją ulepszyć
		if strings.Contains(asset.GetName(), runtime.GOOS) && strings.Contains(asset.GetName(), runtime.GOARCH) {
			assetURL = asset.GetBrowserDownloadURL()
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("nie znaleziono odpowiedniego pliku binarnego dla %s/%s w wydaniu %s", runtime.GOOS, runtime.GOARCH, release.GetTagName())
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("nie udało się ustalić ścieżki pliku wykonywalnego: %w", err)
	}

	if err := selfupdate.UpdateTo(assetURL, exe); err != nil {
		return err
	}

	return nil
}

func CleanupOldBinary() {
	if runtime.GOOS != "windows" {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	oldExe := exe + ".old"
	_ = os.Remove(oldExe)
}
