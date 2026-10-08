<#
.SYNOPSIS
    Skrip Instalasi Cepat cepot untuk Windows (PowerShell)
.DESCRIPTION
    Mengunduh / mengompilasi binary cepot dan memasangnya ke sistem lokal ($env:LOCALAPPDATA\cepot).
    Dapat dijalankan langsung:
        irm https://raw.githubusercontent.com/n0z0/cepot/main/install.ps1 | iex
    atau dijalankan secara lokal:
        .\install.ps1
#>

[CmdletBinding()]
param (
    [string]$Version = "latest",
    [switch]$ConsoleMode,
    [string]$InstallDir = "$env:LOCALAPPDATA\cepot"
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

# 6. Selesai & Panduan Konfigurasi
Write-Host ""
Write-Host "==============================================" -ForegroundColor Green
Write-Host "          Instalasi cepot Selesai!            " -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Green
Write-Host "Lokasi Executable: $targetExe" -ForegroundColor White
Write-Host ""
Write-Host "Langkah selanjutnya:" -ForegroundColor Yellow
Write-Host "1. Atur API Key Z.AI Anda di PowerShell:" -ForegroundColor White
Write-Host "   [System.Environment]::SetEnvironmentVariable('ZAI_API_KEY', 'api_key_anda', 'User')" -ForegroundColor Cyan
Write-Host ""
Write-Host "2. (Opsional) Atur Telegram Bot API:" -ForegroundColor White
Write-Host "   [System.Environment]::SetEnvironmentVariable('TELEGRAM_BOT_TOKEN', 'token_bot', 'User')" -ForegroundColor Cyan
Write-Host "   [System.Environment]::SetEnvironmentVariable('TELEGRAM_CHAT_ID', 'id_chat', 'User')" -ForegroundColor Cyan
Write-Host ""
Write-Host "3. Jalankan aplikasi:" -ForegroundColor White
Write-Host "   cepot.exe" -ForegroundColor Cyan
Write-Host "   atau gunakan file AutoHotkey di: $ahkPath" -ForegroundColor Gray
Write-Host ""
