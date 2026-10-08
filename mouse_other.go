//go:build !windows

package main

import (
	"fmt"
)

// MoveMouseByAnswer stub untuk sistem non-Windows
func MoveMouseByAnswer(jawaban string) (string, error) {
	option := ParseAnswerOption(jawaban)
	if option == "" {
		return "", fmt.Errorf("jawaban '%s' tidak dapat dipetakan ke pilihan A-E", jawaban)
	}

	dir, exists := DirectionMap[option]
	if !exists {
		return option, fmt.Errorf("arah untuk opsi '%s' tidak ditemukan", option)
	}

	fmt.Printf("[Mouse] (Non-Windows stub) Opsi terdeteksi: %s - %s\n", option, dir.Name)
	return option, nil
}
