# PowerShell script to serve files for wget download
param(
    [int]$Port = 8000
)

Write-Host "Starting HTTP server for file transfer..." -ForegroundColor Green
Write-Host "Serving files from: $(Get-Location)" -ForegroundColor Yellow
Write-Host "Server will be accessible at: http://$(hostname):$Port" -ForegroundColor Cyan
Write-Host ""

# Get local IP addresses
$IPs = Get-NetIPAddress -AddressFamily IPv4 | Where-Object { $_.InterfaceAlias -notlike "*Loopback*" } | Select-Object -ExpandProperty IPAddress
Write-Host "Available on these IPs:" -ForegroundColor Yellow
foreach ($IP in $IPs) {
    Write-Host "  http://$IP`:$Port" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "On your VM (192.168.15.10), run:" -ForegroundColor Yellow
Write-Host "  mkdir -p ~/quic" -ForegroundColor Cyan
Write-Host "  cd ~/quic" -ForegroundColor Cyan

# Find the local network IP (usually 192.168.x.x)
$LocalIP = $IPs | Where-Object { $_ -like "192.168.*" } | Select-Object -First 1
if ($LocalIP) {
    Write-Host "  wget -r -np -nH --cut-dirs=0 -A 'go,mod,sum,pem,html' http://$LocalIP`:$Port/" -ForegroundColor Cyan
} else {
    Write-Host "  wget -r -np -nH --cut-dirs=0 -A 'go,mod,sum,pem,html' http://YOUR_WINDOWS_IP:$Port/" -ForegroundColor Cyan
    Write-Host "  (Replace YOUR_WINDOWS_IP with one of the IPs listed above)" -ForegroundColor Red
}

Write-Host ""
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Red
Write-Host ""

# Start the HTTP server using Python if available
if (Get-Command python -ErrorAction SilentlyContinue) {
    Write-Host "Starting Python HTTP server..." -ForegroundColor Green
    python -m http.server $Port
} elseif (Get-Command python3 -ErrorAction SilentlyContinue) {
    Write-Host "Starting Python3 HTTP server..." -ForegroundColor Green
    python3 -m http.server $Port
} else {
    Write-Host "Python not found. Creating simple Go HTTP server..." -ForegroundColor Yellow
    
    # Create a temporary Go file for serving
    $tempGoFile = "temp_server.go"
    @"
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Printf("Serving files from current directory on port $Port\n")
    log.Fatal(http.ListenAndServe(":$Port", http.FileServer(http.Dir("."))))
}
"@ | Out-File -FilePath $tempGoFile -Encoding UTF8

    try {
        go run $tempGoFile
    } finally {
        Remove-Item $tempGoFile -ErrorAction SilentlyContinue
    }
}