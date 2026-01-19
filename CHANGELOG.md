# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Your new feature here.

## [1.1.0] - 2026-01-19

### Added
- Display changelog from GitHub release notes before updating.

### Fixed
- A bug where the application failed to restart after an update on Linux, especially when using version managers like `mise`. The restart logic now uses `os.Args[0]` and `syscall.Exec` for more robustness.

### Changed
- Refactored the massive `cmd/gopress/root.go` file into smaller, logical modules:
  - `run.go`: Contains the main application execution logic.
  - `wizard.go`: Contains the interactive wizard flow.
  - `utils.go`: Contains common utility functions.
  This improves code readability and maintainability.

## [1.0.0] - 2026-01-18

### Added
- Initial version of GoPress.
- Image processing (to WebP) and resizing.
- ZIP file support.
- WordPress uploader with FileBird integration.
- Interactive wizard and silent mode with flags.
- Self-update mechanism.
