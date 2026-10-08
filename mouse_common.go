package main

import (
	"regexp"
	"strings"
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
