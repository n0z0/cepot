//go:build !windows

package main

func NotifikasiDesktop(judul, pesan string) error {
	return nil
}

func PlayNotificationSound(nomor int) {
	// No-op di platform non-Windows
}

func PlayNotificationSound2(abjad string) {
	// No-op di platform non-Windows
}
