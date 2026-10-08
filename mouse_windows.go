//go:build windows

package main

import (
	"fmt"
	"math"
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

// GetCurrentCursorPos mengambil koordinat kursor mouse saat ini di Windows
func GetCurrentCursorPos() (POINT, error) {
	var pt POINT
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return pt, err
	}
	return pt, nil
}

// SetCursorPosition mengatur posisi kursor mouse ke koordinat X, Y di Windows
func SetCursorPosition(x, y int32) error {
	ret, _, err := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return err
	}
	return nil
}

// GetScreenDimensions mengambil resolusi layar utama di Windows
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

// SmoothMove menggerakkan kursor mouse secara mulus dari posisi awal ke target di Windows
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

// MoveMouseByAnswer menggerakkan mouse ke arah tertentu berdasarkan jawaban (A-E atau 1-5) di Windows
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

	// Antisipasi jika kursor berada terlalu dekat dengan tepi layar.
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
