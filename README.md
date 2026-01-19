# GoPress CLI 

**Uwaga:** To narzędzie zostało stworzone dla polskojęzycznego zespołu. Wszystkie komendy CLI, komunikaty pomocy i interaktywny kreator są w języku polskim.

GoPress to narzędzie automatyzacji napisane w Go, służące do optymalizacji i publikacji obrazów. Pobiera folder zdjęć (JPG, PNG, HEIC), konwertuje je do WebP, zmienia rozmiar i przesyła do WordPressa, odwzorowując lokalną strukturę folderów (wymagana wtyczka FileBird).

## Wymagania systemowe

Ze względu na obsługę formatu HEIC (Apple), aplikacja wykorzystuje CGO i biblioteki systemowe. Aby uruchomić GoPress na Linuxie lub macOS, musisz zainstalować bibliotekę `libheif`.

### Linux (Ubuntu/Debian)
```bash
sudo apt-get update && sudo apt-get install -y libheif-dev
```

### macOS
```bash
brew install libheif
```

### Windows
Wersja dla Windows zawiera wszystkie niezbędne biblioteki statycznie wbudowane w plik `.exe`. Nie wymaga dodatkowej instalacji.

## Uruchomienie

### Windows
Pobierz plik `gopress.exe` i uruchom go dwukrotnie. Windows SmartScreen może wyświetlić ostrzeżenie – kliknij "Więcej informacji" i "Uruchom mimo to".

### Linux / macOS
Pobierz odpowiedni plik binarny, nadaj mu uprawnienia wykonywania i uruchom w terminalu:

```bash
chmod +x gopress-linux
./gopress-linux
```

## Sposób użycia

Aplikacja może działać w dwóch trybach:

1. **Tryb Interaktywny (Kreator):** Uruchom program bez żadnych argumentów. Aplikacja poprowadzi Cię krok po kroku, pytając o folder źródłowy, dane do WordPressa itp.

2. **Tryb Cichy (CLI):** Uruchom program z flagami, aby w pełni zautomatyzować proces.

### Dostępne Flagi

| Flaga | Opis |
| :--- | :--- |
| `--input`, `-i` | Folder źródłowy lub plik .zip. |
| `--output`, `-o` | Folder docelowy (domyślnie: podfolder `webp`). |
| `--upload` | Włącza przesyłanie do WordPressa. |
| `--quality`, `-q` | Jakość WebP (0-100, domyślnie 80). |
| `--width`, `-w` | Maksymalna szerokość w px (domyślnie 2560). |
| `--delete`, `-d` | Usuwa pliki źródłowe po udanej konwersji. |
| `--wp-domain` | Adres strony WordPress. |
| `--wp-user` | Nazwa użytkownika WordPress. |
| `--wp-secret` | Hasło Aplikacji (Application Password). |
| `--fb-token` | Token API wtyczki FileBird (do struktury folderów). |
| `--update` | Sprawdza i instaluje aktualizacje. |

### Przykłady

**Szybka konwersja folderu:**
```bash
./gopress -i "./zdjecia"
```

**Konwersja pliku ZIP:**
```bash
./gopress -i "./wakacje.zip"
```

**Pełna automatyzacja (konwersja + upload do WP):**
```bash
./gopress -i "./nowe-produkty" --upload \
  --wp-domain "https://sklep.pl" \
  --wp-user "admin" \
  --wp-secret "xxxx xxxx xxxx xxxx" \
  --fb-token "moj-token-filebird" \
  --width 1920
```

## Integracja z WordPress

Aby przesyłanie plików działało, musisz wygenerować **Hasło Aplikacji** (Application Password). Nie używaj swojego hasła do logowania!

1. W panelu WP wejdź w **Użytkownicy -> Profil**.
2. Znajdź sekcję "Hasła aplikacji".
3. Nadaj nazwę (np. "GoPress"), wygeneruj hasło i skopiuj je.

Jeśli chcesz zachować strukturę folderów, zainstaluj wtyczkę **FileBird**, wygeneruj token API w jej ustawieniach i podaj go jako parametr `--fb-token`.

## Budowanie ze źródeł

Wymagania: Go 1.25+, Make, Zig (do cross-kompilacji).

```bash
make all
```
Pliki wynikowe pojawią się w folderze `bin/`.
