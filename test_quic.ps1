# PowerShell script to setup and test QUIC/HTTP3 server
param(
    [string]$ServerIP = "127.0.0.1",
    [switch]$TestOnly = $false,
    [switch]$StartChrome = $false
)

Write-Host "QUIC/HTTP3 Server Setup and Test Script" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host ""

if (-not $TestOnly) {
    Write-Host "Checking for certificate files..." -ForegroundColor Yellow
    
    if (-not (Test-Path "localhost+3.pem") -or -not (Test-Path "localhost+3-key.pem")) {
        Write-Host "Certificate files not found!" -ForegroundColor Red
        Write-Host "Please run the following command first:" -ForegroundColor Yellow
        Write-Host "mkcert localhost 127.0.0.1 ::1 $ServerIP" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "This will create:" -ForegroundColor Yellow
        Write-Host "  - localhost+3.pem" -ForegroundColor Cyan
        Write-Host "  - localhost+3-key.pem" -ForegroundColor Cyan
        exit 1
    }
    
    Write-Host "Certificate files found!" -ForegroundColor Green
    Write-Host ""
}

if ($TestOnly -or $StartChrome) {
    Write-Host "Testing QUIC/HTTP3 connectivity..." -ForegroundColor Yellow
    Write-Host ""
    
    # Test basic HTTPS connectivity
    Write-Host "Testing HTTPS connectivity:" -ForegroundColor Cyan
    try {
        $response = Invoke-WebRequest -Uri "https://$ServerIP`:9444/api/test" -SkipCertificateCheck -TimeoutSec 10
        Write-Host "✓ HTTPS connection successful" -ForegroundColor Green
        
        # Check for Alt-Svc header
        if ($response.Headers.ContainsKey("Alt-Svc")) {
            Write-Host "✓ Alt-Svc header found: $($response.Headers['Alt-Svc'])" -ForegroundColor Green
        } else {
            Write-Host "⚠ Alt-Svc header not found" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "✗ HTTPS connection failed: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    Write-Host ""
    
    # Test using curl if available
    if (Get-Command curl -ErrorAction SilentlyContinue) {
        Write-Host "Testing with curl:" -ForegroundColor Cyan
        
        Write-Host "HTTP/3 test:" -ForegroundColor Yellow
        & curl -v --http3-only -k "https://$ServerIP`:9444/api/test" 2>&1 | Select-String -Pattern "(alt-svc|HTTP/3|protocol)"
        
        Write-Host ""
        Write-Host "HTTP/2 test:" -ForegroundColor Yellow
        & curl -v --http2 -k "https://$ServerIP`:9444/api/test" 2>&1 | Select-String -Pattern "(alt-svc|HTTP/2|protocol)"
    } else {
        Write-Host "curl not found - install curl for more detailed testing" -ForegroundColor Yellow
    }
    
    Write-Host ""
}

if ($StartChrome) {
    Write-Host "Starting Chrome with QUIC enabled..." -ForegroundColor Yellow
    
    $chromePath = "C:\Program Files\Google\Chrome\Application\chrome.exe"
    if (-not (Test-Path $chromePath)) {
        $chromePath = "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe"
    }
    
    if (Test-Path $chromePath) {
        $chromeArgs = "--enable-quic --origin-to-force-quic-on=$ServerIP`:9444"
        Write-Host "Launching: $chromePath $chromeArgs" -ForegroundColor Cyan
        Start-Process -FilePath $chromePath -ArgumentList $chromeArgs
        
        Write-Host ""
        Write-Host "Chrome started with QUIC enabled!" -ForegroundColor Green
        Write-Host "Navigate to: https://$ServerIP`:9444" -ForegroundColor Cyan
    } else {
        Write-Host "Chrome not found at expected location" -ForegroundColor Red
        Write-Host "Please install Chrome or start it manually with:" -ForegroundColor Yellow
        Write-Host "chrome.exe --enable-quic --origin-to-force-quic-on=$ServerIP`:9444" -ForegroundColor Cyan
    }
}

if (-not $TestOnly -and -not $StartChrome) {
    Write-Host "Usage examples:" -ForegroundColor Yellow
    Write-Host "  .\test_quic.ps1 -ServerIP 192.168.0.100 -TestOnly" -ForegroundColor Cyan
    Write-Host "  .\test_quic.ps1 -ServerIP 192.168.0.100 -StartChrome" -ForegroundColor Cyan
    Write-Host "  .\test_quic.ps1 -ServerIP 192.168.0.100 -TestOnly -StartChrome" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Manual testing URLs:" -ForegroundColor Yellow
    Write-Host "  https://$ServerIP`:9444        (HTTPS with HTTP/3)" -ForegroundColor Cyan
    Write-Host "  https://$ServerIP`:9444/api/test (API endpoint)" -ForegroundColor Cyan
    Write-Host "  http://$ServerIP`:8080         (HTTP only)" -ForegroundColor Cyan
}