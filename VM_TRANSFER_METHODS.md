# Method 2: Direct File Transfer to VM

## Option A: Using SCP (if you have SSH access)

# From your Windows machine, copy files to VM:
scp -r "c:\Users\User\OneDrive\Desktop\New folder\quic\*" user@192.168.15.10:~/quic/

# Then SSH to VM and continue setup:
ssh user@192.168.15.10
cd ~/quic
# Follow setup steps below

## Option B: Using wget from a temporary web server

# 1. Start a simple HTTP server on Windows (in your project directory):
# Open PowerShell in your quic directory and run:
python -m http.server 8000
# OR if you have Python 2:
python -m SimpleHTTPServer 8000
# OR using Go:
go run -c "package main; import \"net/http\"; func main() { http.ListenAndServe(\":8000\", http.FileServer(http.Dir(\".\"))) }"

# 2. On your VM (192.168.15.10), download the files:
wget -r -np -nH --cut-dirs=1 http://YOUR_WINDOWS_IP:8000/
# Replace YOUR_WINDOWS_IP with your Windows machine's IP

## Option C: Using GitHub (if you push your code)

# 1. On Windows, commit and push your changes:
git add .
git commit -m "Updated QUIC server for VM deployment"
git push origin main

# 2. On VM, clone the repository:
git clone https://github.com/sereymony01/quic.git
cd quic

## VM Setup Commands (run these on 192.168.15.10)

# Install dependencies
sudo apt update
sudo apt install -y golang-go curl wget

# Install mkcert
curl -JLO "https://dl.filippo.io/mkcert/latest?for=linux/amd64"
chmod +x mkcert-v*-linux-amd64
sudo cp mkcert-v*-linux-amd64 /usr/local/bin/mkcert

# Generate certificates for the VM
mkcert -install
mkcert localhost 127.0.0.1 ::1 192.168.15.10

# Create public directory and simple index.html
mkdir -p public
echo '<!DOCTYPE html>
<html><head><title>HTTP/3 VM Server</title></head>
<body>
<h1>HTTP/3 Server Running on VM!</h1>
<p>Protocol: <span id="proto">checking...</span></p>
<script>
fetch("/api/test").then(r=>r.json()).then(d=>document.getElementById("proto").textContent=d.protocol);
</script>
</body></html>' > public/index.html

# Install Go dependencies
go mod tidy

# Build the server
go build -o quic-server main.go

# Run the server (needs sudo for port 443)
sudo ./quic-server

## Testing from client machines

# Test basic HTTPS
curl -vk https://192.168.15.10:443/

# Test HTTP/3 specifically
curl -v --http3-only -k https://192.168.15.10:443/api/test

# Test API endpoint
curl -vk https://192.168.15.10:443/api/test

## Browser setup on client machines

# Copy rootCA.pem from VM to client
scp user@192.168.15.10:~/.local/share/mkcert/rootCA.pem ./

# Install on client (see previous guides for OS-specific steps)

# Start browser with QUIC
google-chrome --enable-quic --origin-to-force-quic-on=192.168.15.10:443
# OR on Windows:
"C:\Program Files\Google\Chrome\Application\chrome.exe" --enable-quic --origin-to-force-quic-on=192.168.15.10:443