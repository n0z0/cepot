# cepot

Aplikasi otomatisasi untuk **Capture and Post** (Pengambilan tangkapan layar, analisis soal berbasis AI GLM-4.6V-Flash / Z.AI, serta indikator jawaban multimedia & gerakan mouse).

---

## 📋 Ikhtisar Proyek (Project Overview)

`cepot` adalah utilitas desktop berbasis Go (Windows) yang dirancang untuk menangkap tampilan layar berisi pertanyaan (misal soal latihan/tes), mengirimkannya ke model visi AI (GLM-4.6V-Flash via API Z.AI), dan mengembalikan jawaban secara instan melalui feedback multisensorik yang hening (*stealth*).

### Alur Eksekusi (Architecture Flow)
```text
[Layar / Soal]
      │
      ▼
1. Screenshot Layar (screenshot.CaptureRect)
      │
      ▼
2. Encode Gambar JPEG & Base64
      │
      ▼
3. Kirim ke API Z.AI (GLM-4.6V-Flash)
      │
      ▼
4. Ekstraksi Jawaban (ParseAnswerOption: A, B, C, D, E)
      │
      ├───────────────────────┼────────────────────────┐
      ▼                       ▼                        ▼
[Gerakan Mouse]        [Audio Notifikasi]      [Desktop Notifikasi]
Meluncur ke salah       Suara unik WAV          Notifikasi tersamar:
satu dari 5 arah        per opsi abjad          "Power & Battery"
```

---

## 📂 Struktur File (Repository Structure)

| File / Folder | Deskripsi |
| :--- | :--- |
| [`main.go`](file:///c:/Users/Windows%2011%20Pro/Documents/joe/cepot/main.go) | Titik masuk utama: alur screenshot, HTTP client Z.AI API, dan orkestrasi feedback. |
| [`mouse.go`](file:///c:/Users/Windows%2011%20Pro/Documents/joe/cepot/mouse.go) | Modul Win32 API (`user32.dll`) untuk pergerakan kursor mouse ke 5 arah dan ekstraksi opsi jawaban (`ParseAnswerOption`). |
| [`mouse_test.go`](file:///c:/Users/Windows%2011%20Pro/Documents/joe/cepot/mouse_test.go) | Unit test untuk parser jawaban dan validasi arah mouse. |
| [`notif.go`](file:///c:/Users/Windows%2011%20Pro/Documents/joe/cepot/notif.go) | Pengatur suara notifikasi audio WAV ter-embed (`beep`) dan notifikasi desktop tersamar (`beeep`). |
| [`art/`](file:///c:/Users/Windows%2011%20Pro/Documents/joe/cepot/art) | Aset suara WAV (`srye.wav`, `ddmushi.wav`, `tot2wuk2.wav`, `utang.wav`) dan ikon baterai (`bat.png`). |

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

## ⚙️ Persyaratan & Konfigurasi

### 1. Environment Variable
Aplikasi memerlukan API Key dari platform [Z.AI (BigModel)](https://open.bigmodel.cn/). Daftarkan API Key pada environment variable sistem dengan nama `ZAI_API_KEY`:

- **Windows CMD**:
  ```cmd
  setx ZAI_API_KEY "your_api_key_here"
  ```
- **PowerShell**:
  ```powershell
  [System.Environment]::SetEnvironmentVariable('ZAI_API_KEY', 'your_api_key_here', 'User')
  ```

<img width="264" height="589" alt="image" src="https://github.com/user-attachments/assets/489709c5-157f-433f-b06d-bf2272ed79c5" />

---

## 🔨 Kompilasi (Build & Testing)

### Menjalankan Unit Test
```sh
go test -v ./...
```

### Opsi Build Binary

1. **Mode Konsol (Standar untuk Debugging)**:
   ```sh
   go build .
   ```
   *Karakteristik*: Menghasilkan `cepot.exe` berbasis console subsystem. Terminal CMD akan muncul dan menampilkan log respons JSON AI serta jumlah penggunaan token.

2. **Mode Senyap / Silent GUI (Rekomendasi Penggunaan Nyata)**:
   ```sh
   go build -ldflags "-H=windowsgui" .
   ```
   *Karakteristik*: Menghasilkan `cepot.exe` bertipe GUI subsystem. Saat dieksekusi, **tidak akan muncul jendela hitam command prompt sama sekali**, sehingga sangat hening dan tidak mengganggu tampilan layar.

---

## 🎯 Panduan Menjalankan Aplikasi

1. Buka materi atau contoh soal di browser (misalnya [https://s.id/cehv13practice](https://s.id/cehv13practice)).
2. Pastikan area soal terlihat jelas di layar dan tidak tertutup jendela aplikasi lain.
3. Jalankan `cepot.exe` (baik via terminal, klik langsung, atau tombol shortcut).
4. Amati sinyal jawaban yang diberikan:
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
