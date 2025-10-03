@echo off
echo Starting Chrome with QUIC enabled for your server...
echo.
echo Please replace 192.168.0.XXX with your actual server IP address
echo.
set /p SERVER_IP="Enter your server IP address (e.g., 192.168.0.100): "

if "%SERVER_IP%"=="" (
    echo No IP address provided, using localhost
    set SERVER_IP=127.0.0.1
)

echo.
echo Starting Chrome with QUIC enabled for %SERVER_IP%:9444
echo.

"C:\Program Files\Google\Chrome\Application\chrome.exe" --enable-quic --origin-to-force-quic-on=%SERVER_IP%:9444

echo.
echo Chrome should now be running with QUIC enabled.
echo Navigate to: https://%SERVER_IP%:9444
echo.
pause