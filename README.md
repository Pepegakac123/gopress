# GoPress CLI 🚀

> ⚠️ **Uwaga:** To narzędzie zostało stworzone dla polskojęzycznego zespołu. Wszystkie komendy CLI, komunikaty pomocy, opisy flag i interaktywny kreator są w **języku polskim**.

**GoPress** to inteligentne narzędzie automatyzacji napisane w **Go (Golang)**, zaprojektowane, aby zaoszczędzić godziny ręcznej pracy przy publikowaniu obrazów w sieci.

Pobiera folder pełen surowych zdjęć (JPG, PNG, HEIC), optymalizuje je do użytku w internecie (WebP), inteligentnie zmienia ich rozmiar i przesyła do WordPressa, **odwzorowując lokalną strukturę folderów** bezpośrednio w bibliotece mediów.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/build-passing-brightgreen)

## ✨ Kluczowe Funkcje

* **⚡ Szybkie i Wydajne:** Przetwarza wiele obrazów jednocześnie (Worker Pools), znacznie szybciej niż ręczna konwersja.
* **📦 Obsługa ZIP:** Wrzuć spakowany folder (ZIP), a program zajmie się resztą – rozpakuje, przetworzy i posprząta.
* **🖼️ Inteligentna Optymalizacja:**
    * Konwertuje wszystkie formaty (JPG, PNG, TIFF, **iPhone HEIC**) do **WebP**.
    * **Smart Resize:** Automatycznie zmniejsza ogromne obrazy do rozmiaru przyjaznego dla sieci (np. 1920px), ale zachowuje małe obrazy bez zmian.
* **📂 Odwzorowanie Folderów:** Jeśli używasz wtyczki **FileBird**, GoPress odtwarza lokalne foldery (np. `2024/Lato/Wydarzenia`) wewnątrz WordPressa automatycznie.
* **🧙‍♂️ Łatwe dla Każdego:** Nie musisz być programistą. Po prostu uruchom, a **Kreator** poprowadzi Cię krok po kroku.

---

## 🎛️ Dostępne Opcje (Flagi)

Możesz kontrolować program używając tych "flag", jeśli chcesz pominąć Kreatora.

| Flaga | Opis | Domyślne zachowanie (jeśli nie ustawione) |
| :--- | :--- | :--- |
| `--input`, `-i` | **Folder źródłowy** zawierający obrazy (lub plik **.zip**). | Program zapyta przez Kreatora. |
| `--output`, `-o` | **Gdzie zapisać** zoptymalizowane pliki WebP. | Dla folderu: wewnątrz. Dla ZIP: obok pliku ZIP. |
| `--quality`, `-q` | **Jakość obrazu** (0-100). Niższa = mniejszy rozmiar pliku. | Używa **80** (świetna równowaga jakości/rozmiaru). |
| `--width`, `-w` | **Maksymalna szerokość** w pikselach. Obrazy szersze zostaną przeskalowane. | Używa **2560px**. (Małe obrazy NIE są powiększane). |
| `--upload` | **Włącz przesyłanie** do WordPressa. | Tylko konwertuje pliki lokalnie. |
| `--delete`, `-d` | **Usuń oryginały**. Usuwa pliki źródłowe po sukcesie. | Zachowuje oryginalne pliki bezpiecznie. |
| `--fb-token` | **Token FileBird**. Włącza odwzorowanie folderów w WP. | Przesyła obrazy płasko (bez folderów). |
| `--wp-domain` | URL Twojej strony (np. `https://strona.pl`). | Program zapyta przez Kreatora. |
| `--wp-user` | Twoja nazwa użytkownika WordPress. | Program zapyta przez Kreatora. |
| `--wp-secret` | Twoje **Hasło Aplikacji** (Nie hasło logowania!). | Program zapyta przez Kreatora. |

---

## 📖 Jak Używać (Przewodnik Użytkownika)

Wybierz swój system operacyjny poniżej, aby zobaczyć jak uruchomić narzędzie.

<details>
<summary><strong>🪟 Windows (Kliknij aby rozwinąć)</strong></summary>

### 1. Pobierz
Pobierz plik `gopress.exe` z najnowszego Release.

### 2. Pierwsze Uruchomienie (Ostrzeżenie Bezpieczeństwa)
Ponieważ to narzędzie jest zbudowane wewnętrznie i nie jest "podpisane" płatnym certyfikatem firmowym, Windows **SmartScreen** może próbować je zablokować.
* Kliknij **"Więcej informacji"**.
* Kliknij **"Uruchom mimo to"**.
* *To dzieje się tylko raz.*

### 3. Jak to uruchomić?

**Metoda A: Kreator (Najłatwiejsza)**
1.  Po prostu **kliknij dwukrotnie** `gopress.exe` gdziekolwiek się znajduje.
2.  Pojawi się czarne okno (terminal).
3.  Odpowiedz na pytania (przeciąganie i upuszczanie folderów do okna również działa!).

**Metoda B: Zaawansowany Użytkownik (Linia Komend)**
1.  Otwórz PowerShell lub CMD.
2.  Przejdź do folderu z narzędziem.
3.  Uruchom z flagami, aby pominąć pytania:
```powershell
    .\gopress.exe -i "C:\MojeZdjecia" --upload
```
</details>

<details>
<summary><strong>🍎 macOS (Kliknij aby rozwinąć)</strong></summary>

### 1. Pobierz
Pobierz plik binarny dla swojego Maca (`gopress-mac-m1` dla Apple Silicon lub `gopress-mac-intel`).

### 2. Uprawnienia
MacOS jest restrykcyjny. Musisz pozwolić na uruchomienie pliku.
1.  Otwórz **Terminal**.
2.  Wpisz `chmod +x ` i przeciągnij plik do okna terminala.
3.  Naciśnij Enter.

### 3. Pierwsze Uruchomienie (Ostrzeżenie Bezpieczeństwa)
1.  **Kliknij prawym przyciskiem** plik w Finderze.
2.  Wybierz **Otwórz**.
3.  Kliknij **Otwórz** w oknie dialogowym (to dodaje aplikację do białej listy).

### 4. Jak to uruchomić?
Przeciągnij plik do Terminala i naciśnij Enter, lub uruchom:
```bash
./gopress-mac-m1
```
</details>

<details>
<summary><strong>🐧 Linux (Kliknij aby rozwinąć)</strong></summary>

1.  Pobierz `gopress-linux`.
2.  Nadaj uprawnienia wykonywania: `chmod +x gopress-linux`.
3.  Uruchom: `./gopress-linux`.

</details>

---

## 💡 Przykłady

### 1. Podejście "Chcę być prowadzony" (Kreator)

Po prostu kliknij dwukrotnie aplikację. Zapyta Cię:

* *"Gdzie są zdjęcia?"*
* *"Czy chcesz je przesłać?"*
* *"Jakie jest Twoje hasło WP?"*

### 2. Podejście "Szybka Konwersja"

Konwertuj wszystkie obrazy w folderze `raw`. Ponieważ `--output` nie jest podany, automatycznie tworzy folder `raw/webp`.
```bash
gopress -i "./raw"
```

**LUB użyj pliku ZIP:**
```bash
gopress -i "./zdjecia.zip"
# Stworzy folder ./webp obok pliku zip
```

### 3. Podejście "Pełna Automatyzacja"

Konwertuj, zmień rozmiar do Full HD (1920px) i prześlij do WordPressa zachowując strukturę folderów:
```bash
gopress -i "./zdjecia" --upload \
  --wp-domain "https://mojastrona.pl" \
  --wp-user "admin" \
  --wp-secret "xxxx xxxx xxxx xxxx" \
  --fb-token "twoj-token-api-filebird" \
  --width 1920
```

---

## 🔌 Integracja z WordPressem

Aby przesyłanie działało, potrzebujesz **Hasła Aplikacji**. Jest to bezpieczniejsze niż prawdziwe hasło.

1.  Przejdź do **WP Admin** -> **Użytkownicy** -> **Profil**.
2.  Przewiń w dół do "Hasła aplikacji".
3.  Nazwij je "GoPress", utwórz i skopiuj kod.
4.  Wklej ten kod do GoPress gdy zostaniesz poproszony.

**Bonus: Wsparcie dla FileBird**
Jeśli chcesz mieć foldery w WordPressie:

1.  Zainstaluj wtyczkę **FileBird**.
2.  Przejdź do Ustawienia -> FileBird -> API i wygeneruj token.
3.  Podaj ten token do GoPress.

---

## 🛠️ Stack Technologiczny (Dla Programistów)

* **Język:** Go 1.25+
* **Rdzeń:** `Cobra` (CLI), `Viper` (Konfiguracja)
* **Współbieżność:** Worker Pools, Mutexy, Liczniki Atomowe
* **Grafika:** `imaging` (resampling Lanczos3), `goheif` (wiązania CGO dla HEIC)
* **System Budowania:** Zig (Cross-kompilacja)

## 📦 Budowanie ze Źródeł

Wymagania: **Go 1.25+** i **Zig**.
```bash
git clone https://github.com/twoja-nazwa-uzytkownika/gopress.git
cd gopress
make windows  # Buduje bin/gopress.exe
```

## 📄 Licencja

Dystrybuowane na licencji MIT.

---

*Zbudowane z ❤️ w Go.*