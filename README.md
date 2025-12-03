# GoPress CLI 🚀

> ⚠️ **Note:** This tool was developed for a Polish-speaking team. All CLI commands, help messages, flags descriptions, and interactive wizard prompts are in **Polish**.

**GoPress** is a smart automation tool written in **Go (Golang)** designed to save hours of manual work when publishing images for the web.

It takes a folder full of raw images (JPG, PNG, HEIC), optimizes them for the web (WebP), resizes them intelligently, and uploads them to WordPress, **mirroring your local folder structure** directly into the media library.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/build-passing-brightgreen)

## ✨ Key Features

* **⚡ Fast & Efficient:** Processes multiple images at the same time (Worker Pools), making it much faster than manual conversion.
* **🖼️ Smart Optimization:**
    * Converts all formats (JPG, PNG, TIFF, **iPhone HEIC**) to **WebP**.
    * **Smart Resize:** Automatically shrinks huge images to a web-friendly size (e.g., 1920px) but keeps small images as they are.
* **📂 Folder Mirroring:** If you use the **FileBird** plugin, GoPress recreates your local folders (e.g., `2024/Summer/Events`) inside WordPress automatically.
* **🧙‍♂️ Easy for Everyone:** You don't need to be a programmer. Just run it, and the **Wizard** will guide you step-by-step.

---

## 🎛️ Available Options (Flags)

You can control the program using these "flags" if you want to skip the Wizard.

| Flag | Description | Default Behavior (if not set) |
| :--- | :--- | :--- |
| `--input`, `-i` | **The source folder** containing your images. | Program asks you via Wizard. |
| `--output`, `-o` | **Where to save** the optimized WebP files. | Creates a **`webp`** folder inside your input folder. |
| `--quality`, `-q` | **Image Quality** (0-100). Lower = smaller file size. | Uses **80** (Great balance of quality/size). |
| `--width`, `-w` | **Max Width** in pixels. Images wider than this will be resized. | Uses **2560px**. (Small images are NOT stretched). |
| `--upload` | **Enable Upload** to send files to WordPress. | Only converts files locally. |
| `--delete`, `-d` | **Delete Originals**. Removes source files after success. | Keeps original files safe. |
| `--fb-token` | **FileBird Token**. Enables folder mirroring on WP. | Uploads images flat (no folders). |
| `--wp-domain` | Your Website URL (e.g., `https://site.com`). | Program asks you via Wizard. |
| `--wp-user` | Your WordPress Username. | Program asks you via Wizard. |
| `--wp-secret` | Your **Application Password** (Not login password!). | Program asks you via Wizard. |

---

## 📖 How to Use (User Guide)

Select your operating system below to see how to run the tool.

<details>
<summary><strong>🪟 Windows (Click to expand)</strong></summary>

### 1. Download
Get the `gopress.exe` file from the latest Release.

### 2. First Run (Security Warning)
Because this tool is built internally and not "signed" with a paid corporate certificate, Windows **SmartScreen** might try to block it.
* Click **"More Info" (Więcej informacji)**.
* Click **"Run Anyway" (Uruchom mimo to)**.
* *This only happens once.*

### 3. How to run it?

**Method A: The Wizard (Easiest)**
1.  Just **double-click** `gopress.exe` wherever it is.
2.  A black window (terminal) will appear.
3.  Answer the questions (drag & drop folders into the window works too!).

**Method B: Power User (Command Line)**
1.  Open PowerShell or CMD.
2.  Navigate to the folder with the tool.
3.  Run it with flags to skip questions:
    ```powershell
    .\gopress.exe -i "C:\MyPhotos" --upload
    ```
</details>

<details>
<summary><strong>🍎 macOS (Click to expand)</strong></summary>

### 1. Download
Get the binary for your Mac (`gopress-mac-m1` for Apple Silicon or `gopress-mac-intel`).

### 2. Permissions
MacOS is strict. You need to allow the file to run.
1.  Open **Terminal**.
2.  Type `chmod +x ` and drag the file into the terminal window.
3.  Press Enter.

### 3. First Run (Security Warning)
1.  **Right-click** the file in Finder.
2.  Select **Open**.
3.  Click **Open** in the dialog box (this whitelists the app).

### 4. How to run it?
Drag the file into your Terminal and press Enter, or run:
```bash
./gopress-mac-m1
````

\</details\>

\<details\>
\<summary\>\<strong\>🐧 Linux (Click to expand)\</strong\>\</summary\>

1.  Download `gopress-linux`.
2.  Make it executable: `chmod +x gopress-linux`.
3.  Run it: `./gopress-linux`.

\</details\>

-----

## 💡 Examples

### 1\. The "I want to be guided" approach (Wizard)

Simply double-click the app. It will ask you:

  * *"Where are the photos?"*
  * *"Do you want to upload them?"*
  * *"What is your WP password?"*

### 2\. The "Quick Convert" approach

Convert all images in `raw` folder. Since `--output` is missing, it creates a `raw/webp` folder automatically.

```bash
gopress -i "./raw"
```

### 3\. The "Full Automation" approach

Convert, Resize to Full HD (1920px), and Upload to WordPress preserving folder structure:

```bash
gopress -i "./photos" --upload \
  --wp-domain "[https://mysite.com](https://mysite.com)" \
  --wp-user "admin" \
  --wp-secret "xxxx xxxx xxxx xxxx" \
  --fb-token "your-filebird-api-token" \
  --width 1920
```

-----

## 🔌 WordPress Integration

To make the upload work, you need an **Application Password**. This is safer than your real password.

1.  Go to your **WP Admin** -\> **Users** -\> **Profile**.
2.  Scroll down to "Application Passwords".
3.  Name it "GoPress", create it, and copy the code.
4.  Paste this code into GoPress when asked.

**Bonus: FileBird Support**
If you want folders in WordPress:

1.  Install **FileBird** plugin.
2.  Go to Settings -\> FileBird -\> API and generate a token.
3.  Provide this token to GoPress.

-----

## 🛠️ Tech Stack (For Developers)

  * **Language:** Go 1.25+
  * **Core:** `Cobra` (CLI), `Viper` (Config)
  * **Concurrency:** Worker Pools, Mutexes, Atomic Counters
  * **Graphics:** `imaging` (Lanczos3 resampling), `goheif` (CGO bindings for HEIC)
  * **Build System:** Zig (Cross-compilation)

## 📦 Building from Source

Requirements: **Go 1.25+** and **Zig**.

```bash
git clone [https://github.com/your-username/gopress.git](https://github.com/your-username/gopress.git)
cd gopress
make windows  # Builds bin/gopress.exe
```

## 📄 License

Distributed under the MIT License.

-----

*Built with ❤️ in Go.*
