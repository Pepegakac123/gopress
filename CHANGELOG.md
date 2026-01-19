# Changelog

Wszystkie znaczące zmiany w tym projekcie będą dokumentowane w tym pliku.

Format bazuje na [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
a projekt przestrzega [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2026-01-19

W tej wersji skupiliśmy się na poprawie doświadczenia podczas aktualizacji aplikacji. Chcemy, abyś zawsze wiedział, co się zmienia, zanim klikniesz "Tak".

### Added
- **Podgląd zmian przed aktualizacją:** Od teraz, gdy GoPress wykryje nową wersję, wyświetli pełną listę nowości i poprawek bezpośrednio w terminalu. Działa to zarówno w trybie kreatora, jak i w trybie automatycznym (`--update`).
- **Kumulatywna historia zmian:** Jeśli ominąłeś kilka aktualizacji (np. przechodzisz z 1.0.0 na 1.2.0), program pokaże Ci zmiany ze *wszystkich* wersji po drodze, a nie tylko z najnowszej. Masz pełny obraz sytuacji.
- **Integracja z GitHub API:** Aplikacja pobiera teraz dane bezpośrednio z GitHub Releases, co gwarantuje, że informacje są zawsze aktualne i bogate w szczegóły.

### Changed
- Zmieniono sposób budowania opisów wydań na GitHubie. Od teraz są one automatycznie pobierane z tego pliku (`CHANGELOG.md`), co zapewnia spójność między dokumentacją a wydaniem.

## [1.1.0] - 2026-01-19

Ta wersja to przede wszystkim "porządki pod maską" oraz ważna poprawka dla użytkowników Linuxa.

### Fixed
- **Naprawa restartu na Linuxie:** Rozwiązano problem, przez który aplikacja nie uruchamiała się ponownie automatycznie po aktualizacji na systemach Linux (szczególnie w środowiskach używających menedżerów wersji jak `mise` czy `asdf`). Mechanizm restartu jest teraz znacznie bardziej niezawodny.

### Changed
- **Gruntowna przebudowa kodu:** Główny plik aplikacji (`root.go`), który stał się zbyt duży, został podzielony na mniejsze, logiczne moduły (`run.go`, `wizard.go`, `utils.go`). Choć nie widać tego "na zewnątrz", ułatwi to nam szybsze dodawanie nowych funkcji w przyszłości i zmniejsza ryzyko błędów.

## [1.0.0] - 2026-01-18

Pierwsze publiczne wydanie GoPress! 🎉

### Added
- **Inteligentna konwersja:** Automatyczna konwersja folderów ze zdjęciami (JPG, PNG, HEIC) do formatu WebP.
- **Smart Resize:** Inteligentne skalowanie obrazów – zmniejszamy gigantyczne zdjęcia, ale nie ruszamy tych, które już są małe.
- **Obsługa ZIP:** Możliwość podania pliku `.zip` jako wejścia – program sam go rozpakuje, przetworzy i posprząta.
- **WordPress Upload:** Automatyczne przesyłanie przetworzonych zdjęć do biblioteki mediów WordPress.
- **Integracja z FileBird:** Unikalna funkcja odwzorowywania struktury lokalnych folderów w wirtualnych folderach wtyczki FileBird na WordPressie.
- **Dwa tryby pracy:**
    - **Interaktywny Kreator:** Prowadzi za rękę krok po kroku.
    - **Tryb Cichy (CLI):** Pełna automatyzacja za pomocą flag (np. w skryptach CI/CD).
- **Auto-Update:** Wbudowany mechanizm samodzielnej aktualizacji aplikacji.