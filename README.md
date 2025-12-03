# GoPress CLI 🚀

> ⚠️ **Note:** This tool was developed for a Polish-speaking team. All CLI commands, help messages, flags descriptions, and interactive wizard prompts are in **Polish**.

**GoPress** is a high-performance CLI tool written in **Go (Golang)** designed to automate the tedious process of preparing and publishing images for the web.

It recursively scans local directories, optimizes images (Smart Resize + WebP conversion), and uploads them to WordPress, **mirroring your local folder structure** directly into the media library using the FileBird integration.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/build-passing-brightgreen)

## ✨ Key Features

* **⚡ High Performance:** Utilizes **Worker Pools** and Goroutines to process and upload images concurrently, maximizing CPU and network usage.
* **🖼️ Smart Optimization:**
    * Converts JPG, PNG, BMP, TIFF, and **HEIC (iPhone)** to optimized **WebP**.
    * **Smart Downscaling:** Reduces image resolution only if it exceeds the limit (e.g., 2560px), preventing upscaling artifacts.
* **📂 Folder Mirroring (FileBird):** Automatically replicates your local directory structure (e.g., `2024/Summer/Events`) inside WordPress Media Library using FileBird API.
* **🛡️ Thread Safety:** Implements **Mutexes** and **Atomic** counters to safely manage state across concurrent workers.
* **🛑 Graceful Shutdown:** Handles `SIGINT` (Ctrl+C) signals correctly, ensuring current tasks are completed before exiting to prevent data corruption.
* **🧙‍♂️ Interactive Wizard:** Features a user-friendly terminal UI (via `survey`) for non-technical users.

## 📖 How to Use (User Guide)

This section is for end-users who want to run the tool without compiling code.

### 🪟 Windows

1.  **Download:** Get `gopress.exe` from the latest Release.
2.  **First Run (Security Warning):**
    * Since this tool is built internally and not signed with a paid certificate, Windows **SmartScreen** will likely block it.
    * Click **"More Info" (Więcej informacji)** -> **"Run Anyway" (Uruchom mimo to)**.
    * *This only happens once.*

**Option A: Interactive Mode (Easiest)**
1.  Move `gopress.exe` to the folder containing your images (optional, but convenient).
2.  **Double-click** `gopress.exe`.
3.  A black terminal window will open.
4.  Answer the questions (Input folder, WordPress URL, etc.) and press Enter.

**Option B: Power User Mode (PowerShell/CMD)**
1.  Open PowerShell/CMD in the folder where `gopress.exe` is located.
2.  Run commands with flags:
    ```powershell
    .\gopress.exe -i "C:\Photos\2024" --upload
    ```

### 🍎 macOS

1.  **Download:** Get the binary for your Mac (`gopress-mac-m1` for Apple Silicon or `gopress-mac-intel`).
2.  **Permissions:** You need to make the file executable. Open Terminal, navigate to the download folder, and run:
    ```bash
    chmod +x gopress-mac-m1
    ```
3.  **Security Warning:** macOS will block apps from unidentified developers.
    * **Right-click** the file in Finder -> Select **Open**.
    * Click **Open** in the dialog box.
    * *(Alternatively: Go to System Settings -> Privacy & Security -> Allow "gopress-mac..." to run).*

**Running the tool:**
You must run it via Terminal:
```bash
# Interactive Wizard
./gopress-mac-m1

# Silent Mode
./gopress-mac-m1 -i "/Users/name/Pictures/Project" --upload
````

### 🐧 Linux

1.  **Download:** Get `gopress-linux`.
2.  **Permissions:**
    ```bash
    chmod +x gopress-linux
    ```
3.  **Run:**
    ```bash
    ./gopress-linux
    ```

-----

## 💡 Examples

### 1\. The "I want to be guided" approach (Wizard)

Run the tool without arguments. It will ask you step-by-step:

  * Where are the photos?
  * Do you want to upload them to WordPress?
  * (If yes) Provide domain, user, and Application Password.
  * (Optional) Provide FileBird Token to mirror folders.

### 2\. The "Quick Convert" approach

Convert all images in `./raw` folder to WebP and save them in `./optimized`:

```bash
gopress -i "./raw" -o "./optimized"
```

### 3\. The "Full Automation" approach

Convert, Resize (max 1920px), and Upload to WordPress preserving folder structure:

```bash
gopress -i "./photos" --upload \
  --wp-domain "[https://mysite.com](https://mysite.com)" \
  --wp-user "admin" \
  --wp-secret "xxxx xxxx xxxx xxxx" \
  --fb-token "your-filebird-api-token" \
  --width 1920
```

-----

## 🛠️ Tech Stack & Architecture

This project follows idiomatic Go patterns and the "Standard Go Project Layout":

  * **Language:** Go 1.25+
  * **CLI Framework:** [Cobra](https://github.com/spf13/cobra) (Commands) & [Viper](https://github.com/spf13/viper) (Config).
  * **Concurrency:** Worker Pool pattern, Channels for job distribution, `sync.WaitGroup` for synchronization, `context.Context` for cancellation propagation.
  * **Image Processing:** `disintegration/imaging` (Resizing), `chai2010/webp` (Encoding), `jdeng/goheif` (HEIC decoding).
  * **Interactive UI:** `AlecAivazis/survey` & `schollz/progressbar`.

## 📦 Building from Source (Developers)

Requirements:

  * **Go 1.25+**
  * **Zig** (Required for cross-compiling CGO dependencies like HEIC support)

<!-- end list -->

```bash
# Clone the repository
git clone [https://github.com/your-username/gopress.git](https://github.com/your-username/gopress.git)
cd gopress

# Build via Makefile (Cross-platform using Zig cc)
make windows  # Creates bin/gopress.exe
make linux    # Creates bin/gopress-linux
make mac      # Creates bin/gopress-mac
```

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

-----

*Built with ❤️ in Go.*

