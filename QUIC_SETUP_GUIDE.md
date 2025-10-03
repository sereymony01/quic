# QUIC/HTTP3 Server Setup Guide

This guide walks you through setting up a QUIC/HTTP3 server with proper SSL certificate trust chain.

## Step 1: Generate SSL Certificates

On the server machine (where your Go application runs), generate certificates that include all necessary IP addresses:

```bash
mkcert localhost 127.0.0.1 ::1 192.168.0.XXX
```

Replace `192.168.0.XXX` with your actual server IP address. This command will generate:
- `localhost+3.pem` (certificate file)
- `localhost+3-key.pem` (private key file)

## Step 2: Capture Root CA

The mkcert tool creates a root CA that needs to be trusted by client machines.

### On Linux Server:
The root CA is located at: `~/.local/share/mkcert/rootCA.pem`

Copy this file to your client machines for the next step.

### On Windows Server:
The root CA is typically located at: `%LOCALAPPDATA%\mkcert\rootCA.pem`

## Step 3: Install Root CA on Client Machines

### Linux Client (Chrome):

1. Copy the `rootCA.pem` from server to client
2. Install it as a trusted certificate:
   ```bash
   sudo cp rootCA.pem /usr/local/share/ca-certificates/mkcert-rootCA.crt
   sudo update-ca-certificates
   ```
3. Verify installation:
   ```bash
   ls -l /etc/ssl/certs | grep mkcert
   ```

### Windows Client (Chrome/Edge):

1. Copy the `rootCA.pem` from server to client
2. Install the certificate:
   - Press `Win+R`, type `mmc`, press Enter
   - Go to File → Add/Remove Snap-in → choose Certificates → click Add → select Computer account
   - Expand Trusted Root Certification Authorities → right-click Certificates → All Tasks → Import
   - Import your `rootCA.pem` file
   - Finish and restart Chrome/Edge

## Step 4: Test QUIC/HTTP3 Connection

### Command Line Test:
After setting up certificates, test the connection:

```bash
curl -vk https://192.168.0.XXX:9444/
```

You should see in the response headers:
```
< alt-svc: h3=":443";
```

This indicates that HTTP/3 is available.

### Advanced curl testing:
```bash
# Test HTTP/3 specifically
curl -v --http3-only -k https://192.168.0.XXX:9444/api/test

# Test HTTP/2 fallback
curl -v --http2 -k https://192.168.0.XXX:9444/api/test

# Test regular HTTP
curl -v http://192.168.0.XXX:8080/api/test
```

## Step 5: Browser Setup

### Linux Chrome:
Start Chrome with QUIC enabled for your server:
```bash
google-chrome --enable-quic --origin-to-force-quic-on=192.168.0.XXX:9444
```

### Windows Chrome:
Start Chrome with QUIC enabled:
```cmd
"C:\Program Files\Google\Chrome\Application\chrome.exe" --enable-quic --origin-to-force-quic-on=192.168.0.XXX:9444
```

### Alternative Browser Setup:
You can also enable QUIC through browser flags:
- Chrome: Navigate to `chrome://flags/#enable-quic` and enable it
- Edge: Navigate to `edge://flags/#enable-quic` and enable it
- Restart the browser after enabling

## Step 6: Verify HTTP/3 is Working

1. Open your browser and navigate to: `https://192.168.0.XXX:9444`
2. Open Developer Tools (F12) → Network tab
3. Refresh the page and look for "h3" in the Protocol column
4. You should see 🚀 emojis in your server logs indicating HTTP/3 connections

## Troubleshooting

### Common Issues:

1. **Certificate Trust Issues:**
   - Make sure the rootCA.pem is properly installed on client machines
   - Restart browser after installing certificates
   - Check Windows Certificate Store or Linux certificate store

2. **Firewall Issues:**
   - Ensure ports 9444 (HTTPS/HTTP3) and 8080 (HTTP) are open
   - Allow the application through Windows Firewall if prompted
   - Check any network firewalls between client and server

3. **HTTP/3 Not Activating:**
   - Try refreshing the page multiple times
   - HTTP/3 may take several requests to activate
   - Check browser flags are enabled
   - Verify Alt-Svc header is present in responses

4. **IP Address Issues:**
   - Make sure you replace `192.168.0.XXX` with your actual server IP
   - The IP in the certificate must match the IP you're connecting to
   - Test with `localhost` first if on the same machine

### Server Logs:
Watch your Go application logs for:
- 🚀 emoji = HTTP/3 connection
- 📡 emoji = HTTP/2 or HTTP/1.1 connection

### Testing Commands Summary:
```bash
# Test basic connectivity
curl -vk https://192.168.0.XXX:9444/

# Test API endpoint with different protocols
curl -v --http3-only -k https://192.168.0.XXX:9444/api/test
curl -v --http2 -k https://192.168.0.XXX:9444/api/test
curl -v http://192.168.0.XXX:8080/api/test

# Test connection migration
curl -v https://192.168.0.XXX:9444/api/migration?session=test1
```

## Security Notes

- The certificates generated are for development/testing purposes
- For production, use proper CA-signed certificates
- The `mkcert` tool is designed for local development environments
- Always use HTTPS in production environments

## Next Steps

Once everything is working:
1. Your Go server will be accessible via HTTP/3 on port 9444
2. HTTP/2 and HTTP/1.1 fallbacks are available on the same port
3. Test the connection migration features using the `/api/migration` endpoint
4. Monitor logs for protocol usage (look for 🚀 emojis for HTTP/3)