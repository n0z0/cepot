#!/usr/bin/env bash
# ==============================================================================
# Installer Script untuk cepot (Linux / macOS)
# ==============================================================================
# Penggunaan:
#   curl -fsSL https://raw.githubusercontent.com/n0z0/cepot/main/install.sh | bash
# Atau dengan argumen:
#   ./install.sh -k "your_zai_api_key" -t "tg_bot_token" -c "tg_chat_id"
# ==============================================================================

set -e

REPO="n0z0/cepot"
VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_DIR="${HOME}/.config/cepot"

ZAI_API_KEY_ARG=""
TG_TOKEN_ARG=""
TG_CHAT_ARG=""
NON_INTERACTIVE=false

# Parsing argumen baris perintah
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -k|--key) ZAI_API_KEY_ARG="$2"; shift 2 ;;
        -t|--token) TG_TOKEN_ARG="$2"; shift 2 ;;
        -c|--chat) TG_CHAT_ARG="$2"; shift 2 ;;
        -y|--non-interactive) NON_INTERACTIVE=true; shift ;;
        *) shift ;;
    esac
done

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
mkdir -p "$CONFIG_DIR"
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

# 5. Konfigurasi Environment & Pembuatan File .env
ENV_FILE="$CONFIG_DIR/.env"
echo ""
echo "--- Konfigurasi Environment Variables ---"

CURRENT_ZAI="${ZAI_API_KEY:-}"
FINAL_ZAI="$CURRENT_ZAI"
if [ -n "$ZAI_API_KEY_ARG" ]; then
    FINAL_ZAI="$ZAI_API_KEY_ARG"
elif [ -z "$CURRENT_ZAI" ] && [ -t 0 ] && [ "$NON_INTERACTIVE" = false ]; then
    read -r -p "Masukkan ZAI_API_KEY (Tekan Enter untuk lewati): " INPUT_KEY
    if [ -n "$INPUT_KEY" ]; then
        FINAL_ZAI="$INPUT_KEY"
    fi
fi

CURRENT_TG_TOKEN="${TELEGRAM_BOT_TOKEN:-}"
FINAL_TG_TOKEN="$CURRENT_TG_TOKEN"
if [ -n "$TG_TOKEN_ARG" ]; then
    FINAL_TG_TOKEN="$TG_TOKEN_ARG"
elif [ -z "$CURRENT_TG_TOKEN" ] && [ -t 0 ] && [ "$NON_INTERACTIVE" = false ]; then
    read -r -p "Masukkan TELEGRAM_BOT_TOKEN (Opsional, Enter untuk lewati): " INPUT_TG
    if [ -n "$INPUT_TG" ]; then
        FINAL_TG_TOKEN="$INPUT_TG"
    fi
fi

CURRENT_TG_CHAT="${TELEGRAM_CHAT_ID:-}"
FINAL_TG_CHAT="$CURRENT_TG_CHAT"
if [ -n "$TG_CHAT_ARG" ]; then
    FINAL_TG_CHAT="$TG_CHAT_ARG"
elif [ -z "$CURRENT_TG_CHAT" ] && [ -t 0 ] && [ "$NON_INTERACTIVE" = false ]; then
    read -r -p "Masukkan TELEGRAM_CHAT_ID (Opsional, Enter untuk lewati): " INPUT_CHAT
    if [ -n "$INPUT_CHAT" ]; then
        FINAL_TG_CHAT="$INPUT_CHAT"
    fi
fi

# Tulis file .env
cat <<EOF > "$ENV_FILE"
# Konfigurasi Environment cepot
ZAI_API_KEY=$FINAL_ZAI
TELEGRAM_BOT_TOKEN=$FINAL_TG_TOKEN
TELEGRAM_CHAT_ID=$FINAL_TG_CHAT
EOF
echo "[+] File konfigurasi .env tersimpan di: $ENV_FILE"

# Salin juga ke samping binary agar selalu terbaca jika dijalankan lokal
cp "$ENV_FILE" "$INSTALL_DIR/.env" 2>/dev/null || true

# 6. Cek PATH
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
echo "Lokasi Binary    : $TARGET_BIN"
echo "Lokasi File .env : $ENV_FILE"
echo ""
if [ -z "$FINAL_ZAI" ]; then
    echo "[!] PERINGATAN: ZAI_API_KEY belum diisi."
    echo "    Silakan edit file: $ENV_FILE"
    echo "    atau tambahkan ke file profile shell Anda: export ZAI_API_KEY=\"api_key_anda\""
else
    echo "[✓] Konfigurasi environment siap digunakan!"
fi
echo ""
echo "Jalankan aplikasi:"
echo "   cepot"
echo ""
