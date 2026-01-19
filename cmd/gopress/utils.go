package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func sanitizePath(path string) string {
	path = strings.TrimSpace(path)
	return strings.Trim(path, "\"'")
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func openFolder(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("nieobsługiwany system operacyjny: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func restartApplication() {
	fmt.Println("Restartowanie aplikacji...")

	executablePath := os.Args[0]

	if !filepath.IsAbs(executablePath) {
		path, err := exec.LookPath(executablePath)
		if err != nil {
			fmt.Printf("Błąd: Nie można znaleźć pliku wykonywalnego '%s' w PATH: %v\n", executablePath, err)
			return
		}
		executablePath = path
	}

	cmd := exec.Command(executablePath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if err := syscall.Exec(executablePath, os.Args, os.Environ()); err != nil {
			fmt.Printf("Błąd podczas restartu (syscall.Exec): %v\n", err)
			// Fallback do cmd.Run() jeśli syscall.Exec zawiedzie
			if err := cmd.Run(); err != nil {
				fmt.Printf("Błąd podczas restartu (cmd.Run fallback): %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		if err := cmd.Run(); err != nil {
			fmt.Printf("Błąd podczas restartu: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
