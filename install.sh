#!/usr/bin/env bash
# ==============================================================================
# Installer Script untuk cepot (Linux / macOS)
# ==============================================================================
# Penggunaan:
#   curl -fsSL https://raw.githubusercontent.com/n0z0/cepot/main/install.sh | bash
# atau secara lokal:
#   chmod +x install.sh && ./install.sh
# ==============================================================================

set -e

REPO="n0z0/cepot"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

echo "=============================================="
echo "        🚀 Installer cepot (Unix)             "
echo "=============================================="

# 1. Deteksi Sistem Operasi
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux*)  PLATFORM="linux" ;;
    darwin*) PLATFORM="darwin" ;;
    *)
        echo "[-] Sistem operasi $OS tidak didukung secara otomatis oleh skrip ini."
        exit 1
        ;;
esac

# 2. Deteksi Arsitektur Mesin
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)   TARGET_ARCH="amd64" ;;
    arm64|aarch64)  TARGET_ARCH="arm64" ;;
    *)
        echo "[-] Arsitektur $ARCH tidak didukung."
        exit 1
        ;;
esac

mkdir -p "$INSTALL_DIR"
TARGET_BIN="$INSTALL_DIR/cepot"

# 3. Cek apakah source code dan Go compiler tersedia secara lokal
BUILT_FROM_SOURCE=false
if [ -f "main.go" ] && command -v go >/dev/null 2>&1; then
    echo "[*] Source code dan Go compiler terdeteksi."
    echo "[*] Mengompilasi binary langsung dari source..."
    if go build -trimpath -ldflags "-s -w" -o "$TARGET_BIN" .; then
        chmod +x "$TARGET_BIN"
        BUILT_FROM_SOURCE=true
        echo "[+] Berhasil mengompilasi binary ke: $TARGET_BIN"
    else
        echo "[-] Gagal mengompilasi lokal, mencoba mengunduh dari rilis GitHub..."
    fi
fi

# 4. Jika belum ter-build dari source, unduh binary rilis dari GitHub
if [ "$BUILT_FROM_SOURCE" = false ]; then
    BINARY_NAME="cepot-${PLATFORM}-${TARGET_ARCH}"
    echo "[*] Mencari rilis $BINARY_NAME dari GitHub ($REPO)..."

    if [ "$VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME"
    else
        DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"
    fi

    echo "[*] Mengunduh $DOWNLOAD_URL ..."
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TARGET_BIN" || {
            echo "[-] Gagal mengunduh $BINARY_NAME dari GitHub Release."
            echo "    Pastikan rilis telah tersedia di https://github.com/$REPO/releases"
            exit 1
        }
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$TARGET_BIN" "$DOWNLOAD_URL" || {
            echo "[-] Gagal mengunduh $BINARY_NAME dari GitHub Release."
            exit 1
        }
    else
        echo "[-] Memerlukan curl atau wget untuk mengunduh binary."
        exit 1
    fi

    chmod +x "$TARGET_BIN"
    echo "[+] Binary berhasil diunduh dan dipasang di: $TARGET_BIN"
fi

# 5. Cek PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "[!] Peringatan: $INSTALL_DIR belum ada di PATH shell Anda."
    echo "    Tambahkan baris berikut ke ~/.bashrc atau ~/.zshrc Anda:"
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo ""
echo "=============================================="
echo "        🎉 Instalasi cepot Selesai!           "
echo "=============================================="
echo "Lokasi Binary : $TARGET_BIN"
echo ""
echo "Langkah Konfigurasi:"
echo "1. Daftarkan API Key Z.AI:"
echo "   export ZAI_API_KEY=\"api_key_anda\""
echo ""
echo "2. (Opsional) Daftarkan Telegram Bot:"
echo "   export TELEGRAM_BOT_TOKEN=\"token_bot\""
echo "   export TELEGRAM_CHAT_ID=\"id_chat\""
echo ""
echo "3. Jalankan aplikasi:"
echo "   cepot"
echo ""
