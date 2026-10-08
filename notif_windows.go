//go:build windows

package main

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/gen2brain/beeep"
)

//go:embed art/srye.wav
var notifSrye []byte

//go:embed art/ddmushi.wav
var notifDDm []byte

//go:embed art/tot2wuk2.wav
var notifTwk []byte

//go:embed art/utang.wav
var notifUtang []byte

//go:embed art/bat.png
var batPng []byte

var (
	iconPath string
	iconOnce sync.Once
)

// initIcon akan memastikan file PNG ada di disk (temp dir) tepat sekali
func initIcon() (string, error) {
	var err error
	iconOnce.Do(func() {
		tmpDir := os.TempDir()
		iconPath = filepath.Join(tmpDir, "indo.png")

		// cek sudah ada atau belum
		if _, statErr := os.Stat(iconPath); os.IsNotExist(statErr) {
			// belum ada → tulis
			err = os.WriteFile(iconPath, batPng, 0o644)
		}
	})
	return iconPath, err
}

func NotifikasiDesktop(judul, pesan string) error {
	path, err := initIcon()
	if err != nil {
		return err
	}

	return beeep.Notify(judul, pesan, path)
}

func PlayNotificationSound(nomor int) {
	// pilih audio yang akan diputar
	var notifWav []byte
	switch nomor {
	case 1:
		notifWav = notifSrye
	case 2:
		notifWav = notifDDm
	case 3:
		notifWav = notifTwk
	case 4:
		notifWav = notifUtang
	default:
		// Kembalikan error jika indeks tidak valid
		notifWav = notifTwk
	}
	// decode WAV dari byte yang di-embed
	streamer, format, err := wav.Decode(bytes.NewReader(notifWav))
	if err != nil {
		log.Fatal("gagal decode wav:", err)
	}
	defer streamer.Close()

	// init speaker sesuai format audio
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal("gagal init speaker:", err)
	}

	done := make(chan struct{})

	// play sampai selesai
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
}

func PlayNotificationSound2(abjad string) {
	// Ekstrak opsi A-E yang bersih
	cleanOpt := ParseAnswerOption(abjad)
	if cleanOpt == "" {
		cleanOpt = strings.ToUpper(strings.TrimSpace(abjad))
	}

	// pilih audio yang akan diputar
	var notifWav []byte
	switch cleanOpt {
	case "A":
		notifWav = notifSrye
	case "B":
		notifWav = notifDDm
	case "C":
		notifWav = notifTwk
	case "D":
		notifWav = notifUtang
	case "E":
		notifWav = notifTwk
	default:
		// Fallback audio jika indeks tidak valid
		notifWav = notifTwk
	}
	// decode WAV dari byte yang di-embed
	streamer, format, err := wav.Decode(bytes.NewReader(notifWav))
	if err != nil {
		log.Fatal("gagal decode wav:", err)
	}
	defer streamer.Close()

	// init speaker sesuai format audio
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal("gagal init speaker:", err)
	}

	done := make(chan struct{})

	// play sampai selesai
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
}
