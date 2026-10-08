package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndSetEnv(t *testing.T) {
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	content := `
# Komentar
TEST_CEPOT_VAR1=helloworld
TEST_CEPOT_VAR2="with quotes"
TEST_CEPOT_VAR3='single quotes'
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("Gagal menulis temp .env file: %v", err)
	}

	// Pastikan env belum ada
	os.Unsetenv("TEST_CEPOT_VAR1")
	os.Unsetenv("TEST_CEPOT_VAR2")
	os.Unsetenv("TEST_CEPOT_VAR3")

	if err := parseAndSetEnv(envFile); err != nil {
		t.Fatalf("parseAndSetEnv gagal: %v", err)
	}

	if got := os.Getenv("TEST_CEPOT_VAR1"); got != "helloworld" {
		t.Errorf("TEST_CEPOT_VAR1 = %q; want %q", got, "helloworld")
	}
	if got := os.Getenv("TEST_CEPOT_VAR2"); got != "with quotes" {
		t.Errorf("TEST_CEPOT_VAR2 = %q; want %q", got, "with quotes")
	}
	if got := os.Getenv("TEST_CEPOT_VAR3"); got != "single quotes" {
		t.Errorf("TEST_CEPOT_VAR3 = %q; want %q", got, "single quotes")
	}
}
