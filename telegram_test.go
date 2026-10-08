package main

import (
	"os"
	"testing"
)

func TestSendPhotoToTelegram_MissingCredentials(t *testing.T) {
	// Pastikan ketika env kosong, mengembalikan error yang deskriptif dan aman tanpa panik/crash
	origToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	origChatID := os.Getenv("TELEGRAM_CHAT_ID")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origToken)
		os.Setenv("TELEGRAM_CHAT_ID", origChatID)
	}()

	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_CHAT_ID")

	err := SendPhotoToTelegram([]byte("dummy_bytes"), "test caption")
	if err == nil {
		t.Fatal("Expected error when credentials are not set, got nil")
	}
}
