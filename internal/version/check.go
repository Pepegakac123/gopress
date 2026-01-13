package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// CurrentVersion to aktualna wersja aplikacji
// Wartość tej zmiennej powinna być nadpisywana podczas budowania za pomocą flagi ldflags
var CurrentVersion = "v0.0.0"

type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// CheckForUpdates sprawdza, czy dostępna jest nowsza wersja na GitHub
func CheckForUpdates(repoOwner, repoName string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd pobierania wersji: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	// Porównanie wersji
	// Zakładamy format vX.Y.Z
	if compareVersions(release.TagName, CurrentVersion) > 0 {
		return &release, nil
	}

	return nil, nil
}

func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("nieobsługiwany system operacyjny")
	}
	return cmd.Start()
}

// compareVersions zwraca 1 jeśli v1 > v2, -1 jeśli v1 < v2, 0 jeśli v1 == v2
func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")
	
	// Jeśli wersje są identyczne
	if v1 == v2 {
		return 0
	}

	// Proste porównanie leksykograficzne (nie jest idealne dla semver, ale wystarczy na start)
	// v1.10.0 > v1.9.0 (z stringami by nie zadziałało poprawnie, ale przy założeniu poprawnych semver zazwyczaj działa)
	// Lepiej by było rozparsować na inty, ale dla uproszczenia:
	
	// Spróbujmy jednak parsować podstawowy semver
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		n1 := 0
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		
		n2 := 0
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}

	return 0
}
