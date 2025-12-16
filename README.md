# cepot

Aplikasi untuk Capture and Post

1. Untuk bisa menggunakan aplikasi ini, perlu memasukkan API Key dari Z.AI di environment variables dengan nama variables `ZAI_API_KEY`

<img width="264" height="589" alt="image" src="https://github.com/user-attachments/assets/489709c5-157f-433f-b06d-bf2272ed79c5" />

2. build aplikasi ini, 
    ```sh
    go build .
    ```
    kemudian buka contoh soal https://s.id/cehv13practice pilih salah satu modul yang ada soalnya. Posisikan soal di browser terlihat jelas di tengah tidak ditutupi window aplikasi apapun, jalankan aplikasi cepot.exe menggunakan windows CMD kecil di sebelah kanan bawah.seharusnya ada notif baterai dan suara yang terdengar. diantara angka 1 dan % adalah jawaban dari soal tersebut.


3. uji coba build dengan opsi
    ```sh
    go build -ldflags -H=windowsgui
    ```
    jelaskan perbedaannya

4. gunakan https://www.autohotkey.com/ untuk membuat tombol F8 dilaptop menjalankan cepot.exe, contoh file cepot.ahk
    ```ahk
    F8::Run "C:\Users\PC-10\Documents\cepot\cepot.exe"
    ```
    Buka soal contoh kemudian pada posisi soal terbukan di browser, pastikan tidak ada window aplikasi apapun yang menutupi soal. Pastikan ketika menekan tombol F8 di laptop, aplikasi cepot.exe akan mengeluarkan notif batere dan suara yang terdengar. buat juga shortcut di taskbar dengan icon yang lucu untuk menjalankan cepot.exe

5. Kembangkan lagi aplikasi ini, dengan cara: setelah capture screen. hasil capture screen user dikirim juga melalui whatsapp ke nomor teman satu team atau ke nomor grup wa satu team. gunakan library https://pkg.go.dev/go.mau.fi/whatsmeow untuk mewujudkan fitur ini.
