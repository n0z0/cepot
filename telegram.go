package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// SendImageToTelegram menerima objek gambar image.Image, meng-encode ke format file JPEG,
// lalu mengirimkannya ke Telegram Bot API (/sendPhoto).
func SendImageToTelegram(img image.Image, caption string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN atau TELEGRAM_CHAT_ID tidak diset di environment variables")
	}

	// Encode objek image ke format file JPEG di memori
	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 90}); err != nil {
		return fmt.Errorf("gagal meng-encode gambar ke JPEG untuk Telegram: %w", err)
	}

	return SendPhotoBytesToTelegram(imgBuf.Bytes(), caption, botToken, chatID)
}

// SendPhotoBytesToTelegram mengirimkan byte file gambar ke Telegram Bot API secara multipart HTTP POST
func SendPhotoBytesToTelegram(photoBytes []byte, caption, botToken, chatID string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken)

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)

	// Field chat_id
	if err := writer.WriteField("chat_id", chatID); err != nil {
		return fmt.Errorf("gagal menulis chat_id: %w", err)
	}

	// Field caption (opsional)
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return fmt.Errorf("gagal menulis caption: %w", err)
		}
	}

	// Field photo (file attachment)
	part, err := writer.CreateFormFile("photo", "screenshot.jpg")
	if err != nil {
		return fmt.Errorf("gagal membuat form file photo: %w", err)
	}

	if _, err := part.Write(photoBytes); err != nil {
		return fmt.Errorf("gagal menulis bytes gambar: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("gagal menutup multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, &reqBody)
	if err != nil {
		return fmt.Errorf("gagal membuat request HTTP: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengirim request ke Telegram Bot API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API mengembalikan status %s: %s", resp.Status, string(respBytes))
	}

	return nil
}
