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

## 🛠️ Tech Stack & Architecture

This project follows idiomatic Go patterns and the "Standard Go Project Layout":

* **Language:** Go 1.25+
* **CLI Framework:** [Cobra](https://github.com/spf13/cobra) (Commands) & [Viper](https://github.com/spf13/viper) (Config).
* **Concurrency:** Worker Pool pattern, Channels for job distribution, `sync.WaitGroup` for synchronization, `context.Context` for cancellation propagation.
* **Image Processing:** `disintegration/imaging` (Resizing), `chai2010/webp` (Encoding), `jdeng/goheif` (HEIC decoding).
* **Interactive UI:** `AlecAivazis/survey` & `schollz/progressbar`.

## 📦 Installation

### Option 1: Download Binary (Recommended for Users)
Download the latest `gopress.exe` from the Releases page. No installation required.

> **⚠️ Note for Windows Users:**
> Since this application is not digitally signed with a corporate certificate, Windows SmartScreen might flag it as unrecognized.
> Click **"More Info"** -> **"Run Anyway"**. The code is open-source and safe.

### Option 2: Build from Source (For Developers)
Requirements: Go 1.21+

```bash
# Clone the repository
git clone [https://github.com/your-username/gopress.git](https://github.com/your-username/gopress.git)
cd gopress

# Build via Makefile (Cross-platform)
make windows  # Creates bin/gopress.exe
make linux    # Creates bin/gopress-linux
make mac      # Creates bin/gopress-mac
````

## 🚀 Usage

You can use GoPress in two modes: **Interactive Wizard** or **Silent CLI**.

### 1\. Interactive Mode (Wizard)

Simply run the executable without arguments. The program will ask you for all necessary details (paths, credentials, settings).

```bash
./gopress.exe
```

### 2\. Silent Mode (CLI / CI/CD)

Ideal for scripts and power users.

```bash
# Basic usage: Convert images in "./raw" and save to "./out"
./gopress.exe -i "./raw" -o "./out"

# Full power: Convert, Optimize, and Upload to WordPress
./gopress.exe -i "./photos" --upload \
  --wp-domain "[https://mysite.com](https://mysite.com)" \
  --wp-user "admin" \
  --wp-secret "abcd xxxx xxxx xxxx" \
  --quality 85 \
  --width 1920
```

### 🚩 Available Flags

| Flag (PL) | Description (EN) | Default |
| :--- | :--- | :--- |
| `--input`, `-i` | Path to source folder with images | (Required) |
| `--output`, `-o` | Path to save processed WebP files | `./[input]/webp` |
| `--quality`, `-q` | WebP Quality (0-100) | `80` |
| `--width`, `-w` | Max width in px (Downscale only) | `2560` |
| `--upload` | Enable upload to WordPress | `false` |
| `--delete`, `-d` | **Delete original files** after success | `false` |
| `--fb-token` | FileBird API Token (for folder support) | `""` |
| `--wp-domain` | WordPress URL (e.g., https://www.google.com/search?q=https://site.com) | `""` |
| `--wp-user` | WordPress Username | `""` |
| `--wp-secret` | WordPress Application Password | `""` |

## 🔌 WordPress Integration

### Application Passwords

GoPress uses **Basic Auth** via WordPress Application Passwords.

1.  Go to your WP Admin -\> Users -\> Profile.
2.  Scroll down to "Application Passwords".
3.  Name it "GoPress", create it, and copy the code. **Do not use your login password.**

### FileBird Support (Folder Mirroring)

To enable folder synchronization:

1.  Install **FileBird** plugin (Lite or Pro) on WordPress.
2.  Go to Settings -\> FileBird -\> API and generate a token.
3.  Pass the token to GoPress via Wizard or `--fb-token` flag.

If provided, GoPress will recreate your local directory tree inside the WordPress Media Library automatically.

## 🏗️ Development

To run the project locally with hot-reload or testing:

```bash
# Run directly
go run ./cmd/gopress

# Run tests
go test ./...
```

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

-----

*Built with ❤️ in Go.*
