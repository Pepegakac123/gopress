# Changelog

Wszystkie znaczące zmiany w tym projekcie będą dokumentowane w tym pliku.

Format bazuje na [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
a projekt przestrzega [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Podgląd zmian przed aktualizacją:** Wyświetlanie listy zmian (changelog) pobranej z GitHub Releases przed zainstalowaniem nowej wersji.
- **Kumulatywna historia zmian:** Wyświetlanie zmian ze wszystkich pominiętych wersji przy aktualizacji.
- **Integracja z GitHub API:** Użycie klienta GitHub do pobierania szczegółowych informacji o wydaniach.

## [1.3.2] - 2026-01-15

### Fixed
- **Naprawa restartu na Linuxie:** Rozwiązano krytyczny błąd, przez który aplikacja nie uruchamiała się ponownie po aktualizacji na systemach Linux. Poprawiono mechanizm restartu, aby był bardziej niezawodny, szczególnie w środowiskach z menedżerami wersji (np. `mise`).

## [1.3.1] - 2026-01-14

### Changed
- **Optymalizacja kodu:** Przeprowadzono gruntowną refaktoryzację głównego modułu aplikacji (`root.go`). Kod został podzielony na mniejsze, wyspecjalizowane pliki (`wizard.go`, `run.go`, `utils.go`), co zwiększa stabilność i ułatwia przyszły rozwój.

## [1.3.0] - 2026-01-14

### Added
- **Obsługa FileBird:** Dodano wsparcie dla wtyczki FileBird w WordPressie. Użytkownicy mogą teraz podać token API, aby automatycznie tworzyć i odwzorowywać strukturę folderów w bibliotece mediów WordPress.

## [1.2.1] - 2026-01-14

### Fixed
- Drobne poprawki w logice aktualizacji automatycznych.

## [1.1.1] - 2026-01-14

### Fixed
- Naprawiono błąd przy rozpakowywaniu niektórych plików ZIP zawierających zagnieżdżone katalogi.

## [1.1.0] - 2026-01-06

### Added
- **Obsługa plików ZIP:** Dodano możliwość podania archiwum `.zip` jako źródła zdjęć. Aplikacja automatycznie rozpakuje archiwum, przetworzy zdjęcia i posprząta pliki tymczasowe.
- **Nowe flagi CLI:** Rozszerzono tryb cichy o nowe opcje konfiguracyjne.

## [1.0.1] - 2025-12-19

### Fixed
- Poprawiono obsługę ścieżek plików w systemie Windows (problem z backslashami `\`).

## [1.0.0] - 2025-12-19

Pierwsze publiczne wydanie GoPress! 🎉

### Added
- **Konwersja do WebP:** Automatyczna konwersja obrazów (JPG, PNG, HEIC).
- **Smart Resize:** Inteligentne skalowanie dużych obrazów.
- **WordPress Upload:** Przesyłanie przetworzonych zdjęć do WordPressa.
- **Interaktywny Kreator:** Prosty w obsłudze tryb pytań i odpowiedzi.
- **Tryb Cichy:** Automatyzacja dla zaawansowanych użytkowników.
