<#
.SYNOPSIS
    Skrip Instalasi Cepat cepot untuk Windows (PowerShell)
.DESCRIPTION
    Mengunduh / mengompilasi binary cepot dan memasangnya ke sistem lokal ($env:LOCALAPPDATA\cepot).
    Mendukung konfigurasi otomatis environment variables (ZAI_API_KEY, TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID).
.EXAMPLE
    .\install.ps1 -ZaiApiKey "your_zai_key" -TelegramBotToken "123:ABC" -TelegramChatId "999"
.EXAMPLE
    irm https://raw.githubusercontent.com/n0z0/cepot/main/install.ps1 | iex
#>

[CmdletBinding()]
param (
    [string]$Version = "latest",
    [switch]$ConsoleMode,
    [switch]$NonInteractive,
    [string]$InstallDir = "$env:LOCALAPPDATA\cepot",
    [string]$ZaiApiKey,
    [string]$TelegramBotToken,
    [string]$TelegramChatId
)

$ErrorActionPreference = "Stop"

Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "          Installer cepot untuk Windows       " -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

# 1. Tentukan Direktori Instalasi
if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Host "[+] Direktori instalasi dibuat: $InstallDir" -ForegroundColor Green
}

$targetExe = Join-Path $InstallDir "cepot.exe"

# 2. Cek apakah script dijalankan di folder source code cepot yang memiliki main.go dan Go terpasang
$builtFromSource = $false
if ((Test-Path "main.go") -and (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "[*] Source code dan Go compiler terdeteksi di direktori saat ini." -ForegroundColor Yellow
    Write-Host "[*] Mengompilasi langsung dari source code..." -ForegroundColor Yellow
    
    $ldflags = "-s -w"
    if (-not $ConsoleMode) {
        $ldflags = $ldflags + " -H=windowsgui"
        Write-Host "[*] Mode: Silent GUI (tanpa console window)" -ForegroundColor Gray
    } else {
        Write-Host "[*] Mode: Console" -ForegroundColor Gray
    }

    try {
        go build -trimpath -ldflags "$ldflags" -o "$targetExe" .
        $builtFromSource = $true
        Write-Host "[+] Berhasil mengompilasi binary ke: $targetExe" -ForegroundColor Green
    } catch {
        Write-Warning "Kompilasi lokal gagal: $_. Mencoba mengunduh dari GitHub..."
    }
}

# 3. Jika belum ter-build dari source, unduh rilis dari GitHub
if (-not $builtFromSource) {
    $repo = "n0z0/cepot"
    $arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    
    $binaryName = if ($ConsoleMode) {
        "cepot-windows-" + $arch + ".exe"
    } else {
        "cepot-windows-" + $arch + "-gui.exe"
    }

    Write-Host "[*] Mengambil informasi rilis dari GitHub ($repo)..." -ForegroundColor Yellow
    
    $apiUrl = if ($Version -eq "latest") {
        "https://api.github.com/repos/" + $repo + "/releases/latest"
    } else {
        "https://api.github.com/repos/" + $repo + "/releases/tags/" + $Version
    }

    $downloadUrl = $null
    try {
        $releaseJson = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "cepot-installer" }
        $asset = $releaseJson.assets | Where-Object { $_.name -eq $binaryName }
        if ($asset) {
            $downloadUrl = $asset.browser_download_url
        }
    } catch {
        Write-Host "[-] Catatan rilis belum ditemukan di GitHub Releases API, beralih ke fallback repository URL." -ForegroundColor DarkYellow
    }

    if (-not $downloadUrl) {
        $downloadUrl = "https://github.com/" + $repo + "/releases/download/" + $Version + "/" + $binaryName
    }

    Write-Host "[*] Mengunduh $binaryName ..." -ForegroundColor Yellow
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $targetExe -UseBasicParsing
        Write-Host "[+] Binary berhasil diunduh ke: $targetExe" -ForegroundColor Green
    } catch {
        Write-Host "[-] Gagal mengunduh binary: $_" -ForegroundColor Red
        Write-Host "    Silakan pastikan rilis sudah tersedia di https://github.com/$repo/releases" -ForegroundColor Gray
        exit 1
    }
}

# 4. Tambahkan Direktori Instalasi ke User PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$InstallDir*") {
    $newPath = $userPath + ";" + $InstallDir
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = $env:Path + ";" + $InstallDir
    Write-Host "[+] Direktori $InstallDir telah ditambahkan ke User PATH!" -ForegroundColor Green
} else {
    Write-Host "[OK] $InstallDir sudah terdaftar di PATH." -ForegroundColor Gray
}

# 5. Buat Template AutoHotkey (cepot.ahk) di Folder Instalasi
$ahkPath = Join-Path $InstallDir "cepot.ahk"
if (-not (Test-Path $ahkPath)) {
    $ahkContent = "; Shortcut tombol F8 untuk menjalankan cepot`r`nF8::Run('`"$targetExe`"')"
    [System.IO.File]::WriteAllText($ahkPath, $ahkContent)
    Write-Host "[+] Template AutoHotkey dibuat di: $ahkPath" -ForegroundColor Green
}

# 6. Konfigurasi Environment Variables & File .env
Write-Host ""
Write-Host "--- Konfigurasi Environment Variables ---" -ForegroundColor Cyan

$envFilePath = Join-Path $InstallDir ".env"
$currentZai = [Environment]::GetEnvironmentVariable("ZAI_API_KEY", "User")
if (-not $currentZai) { $currentZai = $env:ZAI_API_KEY }

$canPrompt = [Environment]::UserInteractive -and (-not [Console]::IsInputRedirected) -and (-not $NonInteractive)

# ZAI_API_KEY
if ($ZaiApiKey) {
    [Environment]::SetEnvironmentVariable("ZAI_API_KEY", $ZaiApiKey, "User")
    $env:ZAI_API_KEY = $ZaiApiKey
    Write-Host "[+] ZAI_API_KEY berhasil disimpan ke User Environment!" -ForegroundColor Green
} elseif (-not $currentZai) {
    if ($canPrompt) {
        $promptKey = Read-Host "Masukkan ZAI_API_KEY Anda (Tekan Enter untuk lewati)"
        if ($promptKey) {
            [Environment]::SetEnvironmentVariable("ZAI_API_KEY", $promptKey, "User")
            $env:ZAI_API_KEY = $promptKey
            $ZaiApiKey = $promptKey
            Write-Host "[+] ZAI_API_KEY berhasil disimpan!" -ForegroundColor Green
        }
    }
} else {
    Write-Host "[OK] ZAI_API_KEY sudah terdaftar di sistem." -ForegroundColor Gray
    $ZaiApiKey = $currentZai
}

# TELEGRAM_BOT_TOKEN & TELEGRAM_CHAT_ID
$currentTgToken = [Environment]::GetEnvironmentVariable("TELEGRAM_BOT_TOKEN", "User")
if (-not $currentTgToken) { $currentTgToken = $env:TELEGRAM_BOT_TOKEN }
if ($TelegramBotToken) {
    [Environment]::SetEnvironmentVariable("TELEGRAM_BOT_TOKEN", $TelegramBotToken, "User")
    $env:TELEGRAM_BOT_TOKEN = $TelegramBotToken
    Write-Host "[+] TELEGRAM_BOT_TOKEN berhasil disimpan ke User Environment!" -ForegroundColor Green
} elseif (-not $currentTgToken -and $canPrompt) {
    $promptToken = Read-Host "Masukkan TELEGRAM_BOT_TOKEN (Opsional, tekan Enter untuk lewati)"
    if ($promptToken) {
        [Environment]::SetEnvironmentVariable("TELEGRAM_BOT_TOKEN", $promptToken, "User")
        $env:TELEGRAM_BOT_TOKEN = $promptToken
        $TelegramBotToken = $promptToken
    }
}

$currentTgChat = [Environment]::GetEnvironmentVariable("TELEGRAM_CHAT_ID", "User")
if (-not $currentTgChat) { $currentTgChat = $env:TELEGRAM_CHAT_ID }
if ($TelegramChatId) {
    [Environment]::SetEnvironmentVariable("TELEGRAM_CHAT_ID", $TelegramChatId, "User")
    $env:TELEGRAM_CHAT_ID = $TelegramChatId
    Write-Host "[+] TELEGRAM_CHAT_ID berhasil disimpan ke User Environment!" -ForegroundColor Green
} elseif (-not $currentTgChat -and $canPrompt) {
    $promptChat = Read-Host "Masukkan TELEGRAM_CHAT_ID (Opsional, tekan Enter untuk lewati)"
    if ($promptChat) {
        [Environment]::SetEnvironmentVariable("TELEGRAM_CHAT_ID", $promptChat, "User")
        $env:TELEGRAM_CHAT_ID = $promptChat
        $TelegramChatId = $promptChat
    }
}

# Tulis atau perbarui file .env di direktori instalasi
$finalZai = if ($ZaiApiKey) { $ZaiApiKey } else { $currentZai }
$finalTgToken = if ($TelegramBotToken) { $TelegramBotToken } else { $currentTgToken }
$finalTgChat = if ($TelegramChatId) { $TelegramChatId } else { $currentTgChat }

$envContent = "# Konfigurasi Environment cepot`r`n"
$envContent += "ZAI_API_KEY=$finalZai`r`n"
$envContent += "TELEGRAM_BOT_TOKEN=$finalTgToken`r`n"
$envContent += "TELEGRAM_CHAT_ID=$finalTgChat`r`n"
[System.IO.File]::WriteAllText($envFilePath, $envContent)
Write-Host "[+] File konfigurasi tersimpan di: $envFilePath" -ForegroundColor Green

# 7. Selesai
Write-Host ""
Write-Host "==============================================" -ForegroundColor Green
Write-Host "          Instalasi cepot Selesai!            " -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Green
Write-Host "Lokasi Executable: $targetExe" -ForegroundColor White
Write-Host "Lokasi File .env : $envFilePath" -ForegroundColor White
Write-Host ""
if (-not $finalZai) {
    Write-Host "[!] PERINGATAN: ZAI_API_KEY belum diisi." -ForegroundColor Yellow
    Write-Host "    Silakan edit file: $envFilePath" -ForegroundColor Yellow
    Write-Host "    atau jalankan: [System.Environment]::SetEnvironmentVariable('ZAI_API_KEY', 'api_key_anda', 'User')" -ForegroundColor Cyan
} else {
    Write-Host "[OK] Konfigurasi environment siap digunakan!" -ForegroundColor Green
}
Write-Host ""
Write-Host "Cara Menjalankan:" -ForegroundColor White
Write-Host "   cepot.exe" -ForegroundColor Cyan
Write-Host "   atau aktifkan AutoHotkey: $ahkPath" -ForegroundColor Gray
Write-Host ""
