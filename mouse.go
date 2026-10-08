package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

type POINT struct {
	X int32
	Y int32
}

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

// Direction merepresentasikan arah pergerakan mouse
type Direction struct {
	Option string
	Name   string
	Dx     int32
	Dy     int32
}

// 5 arah pergerakan kursor mouse untuk masing-masing opsi pilihan (A, B, C, D, E atau 1, 2, 3, 4, 5)
// - A: Atas (Up ↑)
// - B: Kanan (Right →)
// - C: Bawah (Down ↓)
// - D: Kiri (Left ←)
// - E: Kanan-Atas / Diagonal (Up-Right ↗)
var DirectionMap = map[string]Direction{
	"A": {Option: "A", Name: "Atas (Up ↑)", Dx: 0, Dy: -220},
	"B": {Option: "B", Name: "Kanan (Right →)", Dx: 220, Dy: 0},
	"C": {Option: "C", Name: "Bawah (Down ↓)", Dx: 0, Dy: 220},
	"D": {Option: "D", Name: "Kiri (Left ←)", Dx: -220, Dy: 0},
	"E": {Option: "E", Name: "Kanan-Atas (Diagonal ↗)", Dx: 160, Dy: -160},
}

var (
	// Pola prefix eksplisit: misal "Jawaban: B", "Pilihan: C", "Answer: A"
	reExplicitPrefix = regexp.MustCompile(`(?i)(?:jawaban|answer|option|pilihan|opsi)\s*[:=\-–]?\s*[*_` + "`" + `]*([A-E1-5])\b`)
	// Pola huruf atau angka yang berdiri sendiri atau diapit tanda kurung/markdown/titik
	reIsolatedOption = regexp.MustCompile(`(?i)(?:^|[\s\(\[\{<*_"'\-:.,])([A-E1-5])(?:$|[\s\)\]\}>*_"'\-:.,])`)
)

// ParseAnswerOption mengekstrak opsi pilihan (A-E) dari teks respons AI
func ParseAnswerOption(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	digitToLetter := func(match string) string {
		match = strings.ToUpper(match)
		switch match {
		case "1":
			return "A"
		case "2":
			return "B"
		case "3":
			return "C"
		case "4":
			return "D"
		case "5":
			return "E"
		default:
			return match
		}
	}

	// 1. Cek jika diawali kata kunci eksplisit seperti "Jawaban: B"
	if matches := reExplicitPrefix.FindStringSubmatch(trimmed); len(matches) > 1 {
		return digitToLetter(matches[1])
	}

	// 2. Cek karakter standalone / isolated option
	if matches := reIsolatedOption.FindStringSubmatch(trimmed); len(matches) > 1 {
		return digitToLetter(matches[1])
	}

	// 3. Fallback: jika string sangat pendek, cari karakter A-E atau 1-5 pertama
	if len(trimmed) <= 5 {
		for _, r := range strings.ToUpper(trimmed) {
			if r >= 'A' && r <= 'E' {
				return string(r)
			}
			if r >= '1' && r <= '5' {
				return digitToLetter(string(r))
			}
		}
	}

	return ""
}

// GetCurrentCursorPos mengambil koordinat kursor mouse saat ini
func GetCurrentCursorPos() (POINT, error) {
	var pt POINT
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return pt, err
	}
	return pt, nil
}

// SetCursorPosition mengatur posisi kursor mouse ke koordinat X, Y
func SetCursorPosition(x, y int32) error {
	ret, _, err := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return err
	}
	return nil
}

// GetScreenDimensions mengambil resolusi layar utama (lebar dan tinggi)
func GetScreenDimensions() (int32, int32) {
	w, _, _ := procGetSystemMetrics.Call(uintptr(SM_CXSCREEN))
	h, _, _ := procGetSystemMetrics.Call(uintptr(SM_CYSCREEN))
	if w == 0 {
		w = 1920
	}
	if h == 0 {
		h = 1080
	}
	return int32(w), int32(h)
}

// SmoothMove menggerakkan kursor mouse secara mulus dari posisi awal ke target
func SmoothMove(startX, startY, targetX, targetY int32, duration time.Duration) {
	const steps = 30
	stepDuration := duration / steps

	for i := 1; i <= steps; i++ {
		// Easing function (Cosine Ease-In-Out) untuk gerakan halus dan natural
		progress := float64(i) / float64(steps)
		ease := (1.0 - math.Cos(progress*math.Pi)) / 2.0

		currX := int32(math.Round(float64(startX) + float64(targetX-startX)*ease))
		currY := int32(math.Round(float64(startY) + float64(targetY-startY)*ease))

		_ = SetCursorPosition(currX, currY)
		time.Sleep(stepDuration)
	}
}

// MoveMouseByAnswer menggerakkan mouse ke arah tertentu berdasarkan jawaban (A-E atau 1-5)
func MoveMouseByAnswer(jawaban string) (string, error) {
	option := ParseAnswerOption(jawaban)
	if option == "" {
		return "", fmt.Errorf("jawaban '%s' tidak dapat dipetakan ke pilihan A-E", jawaban)
	}

	dir, exists := DirectionMap[option]
	if !exists {
		return option, fmt.Errorf("arah untuk opsi '%s' tidak ditemukan", option)
	}

	pt, err := GetCurrentCursorPos()
	if err != nil {
		return option, fmt.Errorf("gagal mendapatkan posisi kursor: %w", err)
	}

	screenWidth, screenHeight := GetScreenDimensions()

	startX := pt.X
	startY := pt.Y

	// Antisipasi jika kursor berada terlalu dekat dengan tepi layar sehingga arah tidak terlihat jelas.
	// Jika terlalu dekat dengan tepi tujuan gerakan, geser sedikit posisi awal ke arah berlawanan terlebih dahulu.
	const margin int32 = 40
	if dir.Dx > 0 && startX+dir.Dx >= screenWidth-margin {
		startX = screenWidth - margin - dir.Dx - 50
	} else if dir.Dx < 0 && startX+dir.Dx <= margin {
		startX = margin - dir.Dx + 50
	}

	if dir.Dy > 0 && startY+dir.Dy >= screenHeight-margin {
		startY = screenHeight - margin - dir.Dy - 50
	} else if dir.Dy < 0 && startY+dir.Dy <= margin {
		startY = margin - dir.Dy + 50
	}

	// Pastikan startX dan startY tetap dalam layar
	if startX < margin {
		startX = margin
	} else if startX > screenWidth-margin {
		startX = screenWidth - margin
	}
	if startY < margin {
		startY = margin
	} else if startY > screenHeight-margin {
		startY = screenHeight - margin
	}

	// Atur posisi awal jika ada penyesuaian
	if startX != pt.X || startY != pt.Y {
		_ = SetCursorPosition(startX, startY)
		time.Sleep(20 * time.Millisecond)
	}

	targetX := startX + dir.Dx
	targetY := startY + dir.Dy

	// Clamp target agar tidak keluar batas layar
	if targetX < margin {
		targetX = margin
	} else if targetX > screenWidth-margin {
		targetX = screenWidth - margin
	}
	if targetY < margin {
		targetY = margin
	} else if targetY > screenHeight-margin {
		targetY = screenHeight - margin
	}

	fmt.Printf("[Mouse] Menggerakkan kursor untuk Opsi %s: %s (dari %d,%d -> %d,%d)\n",
		option, dir.Name, startX, startY, targetX, targetY)

	// Gerakkan mouse secara smooth selama 350ms
	SmoothMove(startX, startY, targetX, targetY, 350*time.Millisecond)

	return option, nil
}
