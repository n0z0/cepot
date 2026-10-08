package main

import (
	"image"
	"image/color"
	"os"
	"testing"
)

func TestSendImageToTelegram_MissingCredentials(t *testing.T) {
	origToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	origChatID := os.Getenv("TELEGRAM_CHAT_ID")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origToken)
		os.Setenv("TELEGRAM_CHAT_ID", origChatID)
	}()

	os.Unsetenv("TELEGRAM_BOT_TOKEN")
	os.Unsetenv("TELEGRAM_CHAT_ID")

	// Buat gambar dummy 10x10
	dummyImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			dummyImg.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}

	err := SendImageToTelegram(dummyImg, "test caption")
	if err == nil {
		t.Fatal("Expected error when credentials are not set, got nil")
	}
}
