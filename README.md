# Cepot Agent

**Cepot Agent** adalah agen AI cerdas (*Agentic AI*) berbasis visi komputer dan aktuasi sistem untuk otomasi **Capture and Post**. Agen ini beroperasi dalam siklus otonom *Perceive ➔ Reason ➔ Act*: mengamati layar (*Perception*), bernalar menggunakan model visi AI GLM-4.6V-Flash (*Reasoning*), mendistribusikan data ke Telegram, dan mengeksekusi aksi nyata ke antarmuka sistem operasi seperti menggerakkan kursor mouse ke 5 arah sebagai sinyal hening, memainkan audio, dan mengirimkan notifikasi (*Action*).

---

## 📋 Ikhtisar Proyek (Project Overview)

`Cepot Agent` adalah implementasi *Task-Oriented Vision & Action Agent* berbasis Go. Agen ini dirancang untuk membaca dan menganalisis konten visual di layar monitor (misal soal latihan/tes) secara mandiri, melakukan inferensi penalaran multimodal via API Z.AI, menyiarkan tangkapan layar secara simultan ke Telegram Bot via goroutine, dan memberikan umpan balik aksi (*grounded action*) secara instan melalui gerakan fisik kursor mouse, audio tersamar, dan notifikasi desktop Windows.

### Alur Eksekusi (Architecture Flow)
```text
[Layar / Soal]
      │
      ▼
1. Screenshot Layar ───► Menghasilkan Objek Gambar (img)
                               │
      ┌────────────────────────┴────────────────────────┐
      ▼ (Goroutine Paralel)                             ▼ (Main Thread)
2a. Kirim Gambar ke Telegram               2b. Encode Gambar ke Base64 (Z.AI)
    - Encode ke format JPEG                    - Encode ke format JPEG
    - Kirim via Telegram Bot API               - Encode ke Base64 Data URL
      (POST /sendPhoto)                                │
                                                       ▼
                                           3. Request ke API Z.AI (GLM-4.6V-Flash)
                                                       │
                                                       ▼
                                           4. Ekstraksi Jawaban (A, B, C, D, E)
                                                       │
               ┌───────────────────────────────────────┼───────────────────────────────────────┐
               ▼                                       ▼                                       ▼
       [Gerakan Mouse]                        [Audio Notifikasi]                     [Desktop Notifikasi]
   Meluncur ke salah satu                       Suara unik WAV                       Notifikasi tersamar:
       dari 5 arah                              per opsi abjad                         "Power & Battery"
```

---

## 📂 Struktur File (Repository Structure)

| File / Folder | Deskripsi |
| :--- | :--- |
| [`main.go`](main.go) | Titik masuk utama: alur screenshot, pemanggilan goroutine Telegram, HTTP client Z.AI API, dan orkestrasi feedback. |
| [`env.go`](env.go) | Modul pembaca otomatis file konfigurasi `.env` dari folder kerja, folder executable, atau LocalAppData. |
| [`env_test.go`](env_test.go) | Unit test parser environment file `.env`. |
| [`telegram.go`](telegram.go) | Modul pengiriman tangkapan layar ke Telegram Bot API (`/sendPhoto`) via multipart HTTP POST. |
| [`telegram_test.go`](telegram_test.go) | Unit test validasi modul Telegram (termasuk penanganan aman jika env kosong). |
| [`mouse_common.go`](mouse_common.go) | Definisi arah, pemetaan koordinat, dan parser jawaban (`ParseAnswerOption`) lintas platform. |
| [`mouse_windows.go`](mouse_windows.go) | Implementasi Win32 API (`user32.dll`) pergerakan kursor mouse ke 5 arah khusus Windows. |
| [`mouse_other.go`](mouse_other.go) | Stub / fallback untuk platform non-Windows. |
| [`notif_windows.go`](notif_windows.go) | Pengatur suara notifikasi audio WAV ter-embed (`beep`) dan notifikasi desktop tersamar (`beeep`) di Windows. |
| [`notif_other.go`](notif_other.go) | Stub notifikasi audio untuk platform non-Windows. |
| [`install.ps1`](install.ps1) | Skrip instalasi otomatis Windows (mendukung build, download, mendaftarkan PATH, & konfigurasi .env). |
| [`install.sh`](install.sh) | Skrip instalasi otomatis Linux & macOS (mendukung argumen CLI, PATH, & konfigurasi .env). |
| [`.env.example`](.env.example) | Contoh template berkas konfigurasi environment. |
| [`.github/workflows/build-and-release.yml`](.github/workflows/build-and-release.yml) | Workflow CI/CD GitHub Actions kompilasi otomatis ke Windows, Linux, dan macOS. |
| [`art/`](art/) | Aset suara WAV (`srye.wav`, `ddmushi.wav`, `tot2wuk2.wav`, `utang.wav`) dan ikon baterai (`bat.png`). |

---

## 🧭 Pemetaan Arah Gerakan Mouse (5 Arah Jawaban)

Saat model AI menentukan pilihan jawaban, kursor mouse akan meluncur secara mulus (*smooth glide* selama ~350 milidetik) ke arah yang sesuai:

| Pilihan Jawaban | Angka | Arah Gerakan Mouse | Gerakan Piksel | Deskripsi Gerakan |
| :---: | :---: | :---: | :---: | :--- |
| **A** | **1** | **↑ Atas (Up)** | `dx: 0, dy: -220` | Kursor meluncur lurus ke atas |
| **B** | **2** | **→ Kanan (Right)** | `dx: +220, dy: 0` | Kursor meluncur lurus ke kanan |
| **C** | **3** | **↓ Bawah (Down)** | `dx: 0, dy: +220` | Kursor meluncur lurus ke bawah |
| **D** | **4** | **← Kiri (Left)** | `dx: -220, dy: 0` | Kursor meluncur lurus ke kiri |
| **E** | **5** | **↗ Kanan-Atas (Diagonal)** | `dx: +160, dy: -160` | Kursor meluncur serong ke kanan atas |

> **Fitur Proteksi Batas Layar (*Boundary Guard*):**  
> Jika kursor pengguna berada terlalu dekat dengan batas layar monitor, posisi awal kursor akan otomatis disesuaikan agar luncuran tetap terlihat utuh dan tidak terpotong di tepi layar.

---

## ⚙️ Persyaratan & Konfigurasi Environment

Aplikasi mendukung 3 cara fleksibel dalam mengatur konfigurasi:
1. **Otomatis via Installer Script** (Paling Mudah)
2. **File `.env`** (Diletakkan di samping `cepot.exe` atau di `%LOCALAPPDATA%\cepot\.env`)
3. **Environment Variable Sistem Operasi**

### Daftar Variabel:
1. **`ZAI_API_KEY` (Wajib)**: API Key dari platform [Z.AI (BigModel)](https://open.bigmodel.cn/).
2. **`TELEGRAM_BOT_TOKEN` (Opsional)**: Token bot Telegram dari [@BotFather](https://t.me/BotFather) untuk mengirim screenshot secara otomatis.
3. **`TELEGRAM_CHAT_ID` (Opsional)**: ID chat / grup target Telegram (misal `123456789` atau `-100123456789`).

---

### Cara 1: Mengatur Saat Menjalankan Installer
- **Windows (PowerShell)**:
  ```powershell
  .\install.ps1 -ZaiApiKey "your_zai_key" -TelegramBotToken "123:ABC" -TelegramChatId "987654"
  ```
  *(Atau jalankan `.\install.ps1` tanpa argumen dan isi prompt interaktif yang muncul)*

- **Linux / macOS**:
  ```bash
  ./install.sh -k "your_zai_key" -t "123:ABC" -c "987654"
  ```

---

### Cara 2: Menggunakan File `.env`
Salin template [`.env.example`](.env.example) menjadi `.env` di samping binary atau di folder instalasi:
```env
ZAI_API_KEY=your_zai_api_key_here
TELEGRAM_BOT_TOKEN=123456789:ABCdefGhIkl_ZYXwvutsRqPoNMLkji
TELEGRAM_CHAT_ID=123456789
```

---

### Cara 3: Menyetel Langsung di OS
- **Windows CMD**:
  ```cmd
  setx ZAI_API_KEY "your_zai_api_key_here"
  setx TELEGRAM_BOT_TOKEN "your_bot_token"
  setx TELEGRAM_CHAT_ID "your_chat_id"
  ```
- **PowerShell**:
  ```powershell
  [System.Environment]::SetEnvironmentVariable('ZAI_API_KEY', 'your_zai_api_key_here', 'User')
  [System.Environment]::SetEnvironmentVariable('TELEGRAM_BOT_TOKEN', 'your_bot_token', 'User')
  [System.Environment]::SetEnvironmentVariable('TELEGRAM_CHAT_ID', 'your_chat_id', 'User')
  ```
- **Linux / macOS**:
  ```bash
  export ZAI_API_KEY="your_zai_api_key_here"
  export TELEGRAM_BOT_TOKEN="your_bot_token"
  export TELEGRAM_CHAT_ID="your_chat_id"
  ```

<img width="264" height="589" alt="image" src="https://github.com/user-attachments/assets/489709c5-157f-433f-b06d-bf2272ed79c5" />

---

## ⚡ Instalasi Cepat (One-Liner Installation)

### Windows (PowerShell)
Jalankan perintah berikut di PowerShell untuk mengunduh/mengompilasi dan memasang `cepot` ke `%LOCALAPPDATA%\cepot` serta otomatis mendaftarkannya ke User PATH:
```powershell
irm https://raw.githubusercontent.com/n0z0/cepot/main/install.ps1 | iex
```
*Atau secara lokal:*
```powershell
.\install.ps1
```
*(Tambahkan `-ConsoleMode` jika ingin versi terminal konsol).*

### Linux & macOS (Bash/Zsh)
Jalankan perintah berikut di terminal:
```bash
curl -fsSL https://raw.githubusercontent.com/n0z0/cepot/main/install.sh | bash
```
*Atau secara lokal:*
```bash
chmod +x install.sh && ./install.sh
```

---

## 🔨 Kompilasi Manual & Workflow Multi-Platform

### 1. Menjalankan Unit Test
```sh
go test -v ./...
```

### 2. Opsi Build Manual (Windows)

- **Mode Konsol (Untuk Debugging / Log Terminal)**:
  ```sh
  go build .
  ```
  *Karakteristik*: Menghasilkan `cepot.exe` berbasis console subsystem. Terminal CMD akan muncul dan menampilkan log respons JSON AI, pengiriman Telegram, serta jumlah penggunaan token.

- **Mode Senyap / Silent GUI (Rekomendasi Pemakaian Asli)**:
  ```sh
  go build -ldflags "-H=windowsgui" .
  ```
  *Karakteristik*: Menghasilkan `cepot.exe` bertipe GUI subsystem. Saat dieksekusi, **tidak akan muncul jendela hitam command prompt sama sekali**, sehingga sangat hening.

### 3. Workflow Otomatis GitHub Actions (`.github/workflows/build-and-release.yml`)
Repositori ini telah dilengkapi dengan workflow CI/CD GitHub Actions yang otomatis mengompilasi binary untuk beberapa platform:
- **Windows**: `cepot-windows-amd64.exe`, `cepot-windows-amd64-gui.exe`, `cepot-windows-arm64.exe`, `cepot-windows-arm64-gui.exe`
- **Linux**: `cepot-linux-amd64`
- **macOS**: `cepot-darwin-amd64` (Intel), `cepot-darwin-arm64` (Apple Silicon M-Series)

Setiap pembuatan tag versi (contoh `git tag v1.0.0 && git push origin v1.0.0`), seluruh binary hasil kompilasi akan otomatis diunggah ke **GitHub Releases**.

---

## 🎯 Panduan Menjalankan Aplikasi

1. Buka materi atau contoh soal di browser (misalnya [https://s.id/cehv13practice](https://s.id/cehv13practice)).
2. Pastikan area soal terlihat jelas di layar dan tidak tertutup jendela aplikasi lain.
3. Jalankan `cepot.exe` (baik via terminal, klik langsung, atau tombol shortcut).
4. Amati alur:
   - Tangkapan layar otomatis dikirimkan ke Telegram di background.
   - **Gerakan Mouse**: Kursor bergerak ke salah satu dari 5 arah.
   - **Notifikasi Desktop**: Timbul notifikasi tersamar *"Energy saver is on 1[jawaban]%"*.
   - **Suara Notifikasi**: Nada audio terputar sesuai indeks pilihan.

---

## ⌨️ Integrasi Tombol Cepat (AutoHotkey)

Untuk kemudahan pemanggilan instan menggunakan tombol keyboard (contoh tombol `F8`), gunakan [AutoHotkey v2](https://www.autohotkey.com/):

Simpan skrip berikut sebagai `cepot.ahk`:
```ahk
#Requires AutoHotkey v2.0
F8::Run("C:\path\ke\cepot.exe")
```
*(Ganti `C:\path\ke\cepot.exe` dengan lokasi absolut binary Anda).*

---

## 🗺️ Rencana Pengembangan (Future Roadmap)

Fitur berikut adalah rencana pengembangan yang dapat ditambahkan pada iterasi selanjutnya:
- Mengintegrasikan pengiriman hasil tangkapan layar dan jawaban ke WhatsApp (grup atau rekan tim) menggunakan library [whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow).
