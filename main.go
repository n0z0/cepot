// main.go go build -ldflags -H=windowsgui
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/kbinani/screenshot"
)

// --- Struct untuk REQUEST ke API z.ai ---
type ZAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Thinking    *Thinking `json:"thinking,omitempty"`
}

type Message struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"` // Untuk request, ini adalah array
}

type ContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Thinking struct {
	Type string `json:"type"` // "enabled" atau "disabled"
}

// --- Struct untuk RESPONSE dari API z.ai ---
type ZAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"` // Menggunakan struct ResponseMessage
	FinishReason string          `json:"finish_reason"`
}

// PERUBAHAN: Struct baru khusus untuk response pesan
type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"` // Untuk response, ini adalah string
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func main() {
	// --- Konfigurasi ---
	apiKey := os.Getenv("ZAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Environment variable ZAI_API_KEY tidak ditemukan. Silakan set terlebih dahulu.")
	}

	modelName := "glm-4.6v-flash"
	promptText := "Tolong bantu saya menjawab soal dengan tema ethical hacking ini, jawaban cukup singkat saja hanya abjad atau angka pilihannya saja untuk menghemat token"
	apiURL := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	fmt.Println("Memulai proses screenshot...")

	// 1. Ambil Screenshot
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Fatalf("Gagal melakukan screenshot: %v", err)
	}

	// 2. Encode gambar ke format JPEG lalu ke Base64
	var imgBuffer bytes.Buffer
	err = jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("Gagal meng-encode gambar ke JPEG: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())
	imageDataURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)

	fmt.Println("Screenshot berhasil di-encode ke Base64.")

	// 3. Buat payload JSON (menggunakan struct untuk request)
	requestPayload := ZAIRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentItem{
					{Type: "text", Text: promptText},
					{Type: "image_url", ImageURL: &ImageURL{URL: imageDataURL}},
				},
			},
		},
		MaxTokens:   1024,
		Temperature: 0.7,
		Stream:      false,
		Thinking:    &Thinking{Type: "enabled"},
	}

	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("Gagal membuat JSON payload: %v", err)
	}

	fmt.Printf("Payload JSON berhasil dibuat. Mengirim request ke model %s...\n", modelName)

	// 4. Buat dan kirim HTTP Request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Gagal membuat HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Gagal mengirim request ke API: %v", err)
	}
	defer resp.Body.Close()

	// 5. Baca dan tampilkan response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Gagal membaca response dari API: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API mengembalikan error. Status: %s, Body: %s", resp.Status, string(respBody))
	}

	// Parsing JSON response menggunakan struct yang benar
	var apiResponse ZAIResponse
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		log.Fatalf("Gagal mem-parsing response JSON: %v", err)
	}
	var konten string

	// Tampilkan hasil dari AI
	if len(apiResponse.Choices) > 0 {
		fmt.Println("\n--- Hasil Analisis dari Z.AI (GLM-4.6V-Flash) ---")
		//fmt.Println(apiResponse.Choices[0].Message.Content) // Akan berhasil karena Content sekarang string
		konten = apiResponse.Choices[0].Message.Content
		fmt.Println(konten)
		fmt.Println("------------------------------------------------")
		fmt.Printf("Penggunaan Token: %d\n", apiResponse.Usage.TotalTokens)
	} else {
		fmt.Println("Tidak ada respons dari model AI.")
		konten = "0"
	}
	//err = NotifikasiDesktop("Power & Battery", "Energy saver is on 1"+strings.TrimSpace(konten)+"%")
	jawaban := strings.TrimSpace(konten)

	// Konversi string ke integer
	if angka, err := strconv.Atoi(jawaban); err == nil {
		PlayNotificationSound(angka)
	} else {
		PlayNotificationSound2(jawaban)
	}
	err = beeep.Notify("Power & Battery", "Energy saver is on 1"+jawaban+"%", "")
	fmt.Println(err)
}
