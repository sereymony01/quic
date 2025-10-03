#!/bin/bash
# VM Setup Script for QUIC Server
# Run this on your 192.168.15.10 VM

echo "QUIC Server Setup on VM"
echo "======================="

# Step 1: Install dependencies
echo "Installing dependencies..."
sudo apt update
sudo apt install -y git golang-go curl wget

# Step 2: Clone the repository
echo "Cloning repository..."
git clone https://github.com/sereymony01/quic.git
cd quic

# Step 3: Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Step 4: Install mkcert
echo "Installing mkcert..."
curl -JLO "https://dl.filippo.io/mkcert/latest?for=linux/amd64"
chmod +x mkcert-v*-linux-amd64
sudo cp mkcert-v*-linux-amd64 /usr/local/bin/mkcert

# Step 5: Generate certificates
echo "Generating certificates..."
mkcert -install
mkcert localhost 127.0.0.1 ::1 192.168.15.10

# Step 6: Create public directory
echo "Creating public directory..."
mkdir -p public
cat > public/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>HTTP/3 Test Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f0f2f5; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .protocol { color: #007acc; font-weight: bold; font-size: 1.2em; }
        .http3 { color: #ff6b35; }
        button { background: #007acc; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; margin: 10px 0; }
        button:hover { background: #005a9e; }
        #result { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>HTTP/3 Test Server (VM)</h1>
        <p>This page is served over <span id="protocol" class="protocol">checking...</span></p>
        
        <h3>Tests:</h3>
        <button onclick="testAPI()">Test API Endpoint</button>
        <button onclick="testProtocol()">Check Protocol</button>
        
        <div id="result"></div>
        
        <h3>Command Line Tests:</h3>
        <pre style="background: #2d3748; color: #e2e8f0; padding: 15px; border-radius: 5px; overflow-x: auto;">
curl -v --http3-only -k https://192.168.15.10:443/api/test
curl -v --http2 -k https://192.168.15.10:443/api/test</pre>
    </div>
    
    <script>
        function updateProtocol() {
            const isHTTPS = window.location.protocol === 'https:';
            let protocolText = isHTTPS ? 'HTTPS' : 'HTTP';
            
            if (isHTTPS) {
                if (navigator.userAgent.includes('Chrome') || navigator.userAgent.includes('Edge')) {
                    protocolText = 'HTTPS (HTTP/3 capable browser)';
                }
            }
            
            document.getElementById('protocol').textContent = protocolText;
            if (protocolText.includes('HTTP/3')) {
                document.getElementById('protocol').className = 'protocol http3';
            }
        }
        
        async function testAPI() {
            showResult('Testing API endpoint...', 'info');
            try {
                const start = performance.now();
                const response = await fetch('/api/test');
                const end = performance.now();
                const data = await response.json();
                
                showResult('<strong>API Test Successful!</strong><br><strong>Response time:</strong> ' + (end - start).toFixed(2) + 'ms<br><strong>Protocol:</strong> ' + data.protocol, 'success');
            } catch (error) {
                showResult('<strong>API Test Failed:</strong><br>' + error.message, 'error');
            }
        }
        
        async function testProtocol() {
            showResult('Checking protocol support...', 'info');
            try {
                const response = await fetch('/api/test');
                const data = await response.json();
                const isHTTP3 = data.protocol.includes('HTTP/3');
                
                showResult('<strong>Protocol Check:</strong><br><strong>Current Protocol:</strong> ' + data.protocol + '<br><strong>HTTP/3 Active:</strong> ' + (isHTTP3 ? 'Yes!' : 'No (using fallback)'), isHTTP3 ? 'success' : 'info');
            } catch (error) {
                showResult('<strong>Protocol Check Failed:</strong><br>' + error.message, 'error');
            }
        }
        
        function showResult(message, type) {
            const resultDiv = document.getElementById('result');
            resultDiv.innerHTML = message;
            resultDiv.className = type || 'info';
        }
        
        updateProtocol();
        setInterval(updateProtocol, 5000);
    </script>
</body>
</html>
EOF

# Step 7: Build the application
echo "Building QUIC server..."
go build -o quic-server main.go

echo ""
echo "Setup complete!"
echo "==============="
echo ""
echo "To start the server:"
echo "  sudo ./quic-server"
echo ""
echo "To test:"
echo "  curl -vk https://192.168.15.10:443/"
echo "  curl -v --http3-only -k https://192.168.15.10:443/api/test"
echo ""
echo "Don't forget to:"
echo "1. Copy the rootCA.pem to client machines"
echo "2. Install it in browser certificate store"
echo "3. Start browser with QUIC flags"