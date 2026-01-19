# GEMINI.md: GoPress Project

## Project Overview

GoPress is a command-line utility written in Go (Golang) designed to automate the process of optimizing and uploading images to WordPress. Its primary user base is Polish-speaking, so all CLI commands, help messages, and interactive prompts are in Polish.

The tool can take a directory of images (or a .zip file), convert them to the WebP format, intelligently resize them, and then upload them to a WordPress media library. A key feature is its ability to replicate the local folder structure within the WordPress media library, which requires the "FileBird" plugin to be active on the WordPress site.

The application can be run in two modes:
1.  **Silent Mode:** By providing command-line flags (e.g., `--input`, `--upload`), the tool can be fully automated.
2.  **Interactive Wizard:** If run without flags, it launches a step-by-step wizard that guides the user through the process.

The project also features a built-in self-update mechanism to keep the binary current.

### Core Technologies
*   **Language:** Go
*   **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
*   **Interactive Wizard:** [Survey](https://github.com/AlecAivazis/survey)
*   **Image Processing:**
    *   `disintegration/imaging` for resizing.
    *   `chai2010/webp` for WebP conversion.
    *   `jdeng/goheif` for HEIC (Apple) image support.
*   **Concurrency:** A worker pool model is used for parallel image processing, maximizing CPU usage.
*   **Cross-Compilation:** The `Makefile` is configured to use `Zig` as a C/C++ cross-compiler to build executables for Windows, macOS, and Linux.

## Building and Running

The project uses a `Makefile` to simplify the build process.

### Building the Application

*   **Build for all platforms (Windows, Linux, macOS):**
    ```bash
    make all
    ```
*   **Build for a specific platform:**
    ```bash
    make windows
    make linux
    make mac
    ```
*   **Clean build artifacts:**
    ```bash
    make clean
    ```
The compiled binaries are placed in the `bin/` directory.

### Running the Application

To run the application after building, you can execute the binary directly.

*   **Run with the interactive wizard:**
    ```bash
    ./bin/gopress
    ```
*   **Run in silent mode (example):**
    ```bash
    ./bin/gopress -i "./path/to/images" --upload --wp-domain "https://example.com" ...
    ```

### Running Tests

The project includes a suite of tests. To run all tests, use the standard Go test command:
```bash
go test ./...
```

## Development Conventions

*   **Code Style:** The code follows standard Go formatting (`gofmt`).
*   **Structure:** The project is organized into `cmd` for the main application entrypoint and `internal` for the core logic, which is a standard Go project layout. The internal packages are divided by concern:
    *   `processor`: Handles image conversion and processing.
    *   `scanner`: Finds images in the source directory.
    *   `uploader`: Manages the upload process to WordPress.
    *   `wordpress`: Contains the WordPress REST API client.
    *   `version`: Implements the self-update check.
*   **Testing:** Tests are located alongside the code they are testing (e.g., `converter.go` and `converter_test.go`). This indicates a convention of writing unit tests for core functionality.
*   **Dependencies:** Go Modules are used for dependency management. All dependencies are listed in `go.mod`.
*   **Localization:** As noted, all user-facing strings are in Polish. This is a hardcoded convention.
**NIE UŻYWAJ IKONEK ANI ZBĘDNYCH KOMENATRZY IMPORTANT**
