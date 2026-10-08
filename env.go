package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnvFiles mencari dan memuat file .env dari beberapa lokasi prioritas:
// 1. Direktori kerja saat ini (./.env)
// 2. Direktori lokasi file executable (<exe_dir>/.env)
// 3. Direktori data pengguna (%LOCALAPPDATA%\cepot\.env atau ~/.config/cepot/.env)
func LoadEnvFiles() {
	locations := []string{
		".env",
	}

	// Cek direktori executable aplikasi
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		locations = append(locations, filepath.Join(exeDir, ".env"))
	}

	// Cek direktori LocalAppData (Windows)
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		locations = append(locations, filepath.Join(localAppData, "cepot", ".env"))
	}

	// Cek direktori home config (Linux / macOS)
	if homeDir, err := os.UserHomeDir(); err == nil {
		locations = append(locations, filepath.Join(homeDir, ".config", "cepot", ".env"))
		locations = append(locations, filepath.Join(homeDir, ".cepot", ".env"))
	}

	loaded := make(map[string]bool)
	for _, path := range locations {
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = path
		}
		if loaded[absPath] {
			continue
		}
		loaded[absPath] = true

		_ = parseAndSetEnv(absPath)
	}
}

// parseAndSetEnv membaca file dan mengeset environment variable jika belum ada di sistem
func parseAndSetEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Hapus kutip ganda atau tunggal di awal dan akhir jika ada
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			if len(val) >= 2 {
				val = val[1 : len(val)-1]
			}
		}

		// Set hanya jika environment variable belum ada di sistem
		if os.Getenv(key) == "" && key != "" {
			_ = os.Setenv(key, val)
		}
	}

	return scanner.Err()
}
