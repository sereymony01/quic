package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// Check if certificates exist
	certFile := "localhost.pem"
	keyFile := "localhost-key.pem"

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		fmt.Println("Certificate files not found!")
		fmt.Println("Please run these commands first:")
		fmt.Println("   mkcert -install")
		fmt.Println("   mkcert localhost 127.0.0.1 ::1")
		fmt.Println("")
		fmt.Println("This will create:")
		fmt.Printf("   %s\n", certFile)
		fmt.Printf("   %s\n", keyFile)
		fmt.Println("\nPress any key to exit...")
		fmt.Scanln()
		return
	}

	// Create static directory if it doesn't exist
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		fmt.Println("Creating static directory...")
		os.Mkdir("./static", 0755)

		// Create a simple index.html
		indexHTML := "<!DOCTYPE html>\n<html>\n<head>\n    <title>HTTP/3 Test Server - Windows</title>\n    <style>\n        body { font-family: Arial, sans-serif; margin: 40px; background: #f0f2f5; }\n        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }\n        .protocol { color: #007acc; font-weight: bold; font-size: 1.2em; }\n        .http3 { color: #ff6b35; }\n        button { background: #007acc; color: white; border: none; padding: 10px 20px; border-radius: 5px; cursor: pointer; margin: 10px 0; }\n        button:hover { background: #005a9e; }\n        #result { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px; }\n        .success { border-left: 4px solid #28a745; }\n        .error { border-left: 4px solid #dc3545; }\n    </style>\n</head>\n<body>\n    <div class=\"container\">\n        <h1>HTTP/3 Test Server (Windows)</h1>\n        <p>This page is served over <span id=\"protocol\" class=\"protocol\">checking...</span></p>\n        \n        <h3>Tests:</h3>\n        <button onclick=\"testAPI()\">Test API Endpoint</button>\n        <button onclick=\"testProtocol()\">Check Protocol</button>\n        <button onclick=\"testMultiple()\">Multiple Requests</button>\n        \n        <div id=\"result\"></div>\n        \n        <h3>Command Line Tests:</h3>\n        <pre style=\"background: #2d3748; color: #e2e8f0; padding: 15px; border-radius: 5px; overflow-x: auto;\">\ncurl -v --http3-only -k https://localhost:9444/api/test\ncurl -v --http2 -k https://localhost:9444/api/test\ncurl -v http://localhost:8080/api/test</pre>\n    </div>\n    \n    <script>\n        function updateProtocol() {\n            const isHTTPS = window.location.protocol === 'https:';\n            let protocolText = isHTTPS ? 'HTTPS' : 'HTTP';\n            \n            if (isHTTPS) {\n                if (navigator.userAgent.includes('Chrome') || navigator.userAgent.includes('Edge')) {\n                    protocolText = 'HTTPS (HTTP/3 capable browser)';\n                }\n            }\n            \n            document.getElementById('protocol').textContent = protocolText;\n            if (protocolText.includes('HTTP/3')) {\n                document.getElementById('protocol').className = 'protocol http3';\n            }\n        }\n        \n        async function testAPI() {\n            showResult('Testing API endpoint...', 'info');\n            try {\n                const start = performance.now();\n                const response = await fetch('/api/test');\n                const end = performance.now();\n                const data = await response.json();\n                \n                showResult('<strong>API Test Successful!</strong><br><strong>Response time:</strong> ' + (end - start).toFixed(2) + 'ms<br><strong>Protocol:</strong> ' + data.protocol + '<br><strong>Response:</strong><br><pre>' + JSON.stringify(data, null, 2) + '</pre>', 'success');\n            } catch (error) {\n                showResult('<strong>API Test Failed:</strong><br>' + error.message, 'error');\n            }\n        }\n        \n        async function testProtocol() {\n            showResult('Checking protocol support...', 'info');\n            try {\n                const response = await fetch('/api/test');\n                const data = await response.json();\n                const isHTTP3 = data.protocol.includes('HTTP/3');\n                \n                showResult('<strong>Protocol Check:</strong><br><strong>Current Protocol:</strong> ' + data.protocol + '<br><strong>HTTP/3 Active:</strong> ' + (isHTTP3 ? 'Yes!' : 'No (using fallback)') + (isHTTP3 ? '' : '<br><em>Try refreshing the page a few times to activate HTTP/3</em>'), isHTTP3 ? 'success' : 'info');\n            } catch (error) {\n                showResult('<strong>Protocol Check Failed:</strong><br>' + error.message, 'error');\n            }\n        }\n        \n        async function testMultiple() {\n            showResult('Running multiple requests to test protocol negotiation...', 'info');\n            const results = [];\n            \n            for (let i = 0; i < 5; i++) {\n                try {\n                    const start = performance.now();\n                    const response = await fetch('/api/test?test=' + (i + 1));\n                    const end = performance.now();\n                    const data = await response.json();\n                    results.push({\n                        test: i + 1,\n                        protocol: data.protocol,\n                        time: (end - start).toFixed(2)\n                    });\n                } catch (error) {\n                    results.push({\n                        test: i + 1,\n                        error: error.message\n                    });\n                }\n                \n                await new Promise(resolve => setTimeout(resolve, 100));\n            }\n            \n            const http3Count = results.filter(r => r.protocol && r.protocol.includes('HTTP/3')).length;\n            let resultHTML = '<strong>Multiple Request Test Results:</strong><br><strong>HTTP/3 requests:</strong> ' + http3Count + '/5<br><br>';\n            \n            for (let r of results) {\n                resultHTML += '<strong>Test ' + r.test + ':</strong> ';\n                if (r.error) {\n                    resultHTML += 'Error: ' + r.error;\n                } else {\n                    resultHTML += r.protocol + ' (' + r.time + 'ms)';\n                }\n                resultHTML += '<br>';\n            }\n            \n            showResult(resultHTML, http3Count > 0 ? 'success' : 'info');\n        }\n        \n        function showResult(message, type) {\n            const resultDiv = document.getElementById('result');\n            resultDiv.innerHTML = message;\n            resultDiv.className = type || 'info';\n        }\n        \n        updateProtocol();\n        setInterval(updateProtocol, 5000);\n    </script>\n</body>\n</html>"

		err := os.WriteFile("./static/index.html", []byte(indexHTML), 0644)
		if err != nil {
			log.Printf("Warning: Could not create index.html: %v", err)
		} else {
			fmt.Println("Created static/index.html")
		}
	}

	mux := http.NewServeMux()

	// File server
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", fs)

	// Connection tracking for migration testing
	connectionMap := make(map[string][]string)

	// API endpoint
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		protocol := r.Proto
		if r.Proto == "HTTP/3.0" {
			protocol = "HTTP/3.0 🚀"
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, `{
  "protocol": "%s",
  "method": "%s",
  "remote_addr": "%s",
  "user_agent": "%s",
  "test_param": "%s",
  "timestamp": "%s",
  "headers_count": %d
}`, protocol, r.Method, r.RemoteAddr, r.UserAgent(),
			r.URL.Query().Get("test"), time.Now().Format("15:04:05"), len(r.Header))
	})

	// Connection migration test endpoint
	mux.HandleFunc("/api/migration", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session")
		if sessionID == "" {
			sessionID = fmt.Sprintf("session_%d", time.Now().Unix())
		}

		clientIP := r.RemoteAddr
		protocol := r.Proto
		if r.Proto == "HTTP/3.0" {
			protocol = "HTTP/3.0 🚀"
		}

		// Track IP changes for this session
		if connectionMap[sessionID] == nil {
			connectionMap[sessionID] = []string{}
		}

		// Check if this is a new IP for this session
		isMigration := false
		ipExists := false
		for _, ip := range connectionMap[sessionID] {
			if ip == clientIP {
				ipExists = true
				break
			}
		}

		if len(connectionMap[sessionID]) > 0 && !ipExists {
			// This is a new IP for this session - migration detected
			isMigration = true
		}

		connectionMap[sessionID] = append(connectionMap[sessionID], clientIP)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")

		// Build JSON array for IP history
		var ipHistoryJSON strings.Builder
		ipHistoryJSON.WriteString("[")
		for i, ip := range connectionMap[sessionID] {
			if i > 0 {
				ipHistoryJSON.WriteString(",")
			}
			ipHistoryJSON.WriteString(`"`)
			ipHistoryJSON.WriteString(ip)
			ipHistoryJSON.WriteString(`"`)
		}
		ipHistoryJSON.WriteString("]")

		fmt.Fprintf(w, `{
  "session_id": "%s",
  "protocol": "%s",
  "current_ip": "%s",
  "ip_history": %s,
  "migration_detected": %t,
  "connection_count": %d,
  "timestamp": "%s",
  "server_time": %d
}`, sessionID, protocol, clientIP,
			ipHistoryJSON.String(),
			isMigration, len(connectionMap[sessionID]),
			time.Now().Format("15:04:05"), time.Now().Unix())
	})

	// Logging middleware
	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protocol := r.Proto
		emoji := "📡"
		if r.Proto == "HTTP/3.0" {
			protocol = "HTTP/3.0 🚀"
			emoji = "🚀"
		}

		// Check if this is a migration endpoint call
		if strings.Contains(r.URL.Path, "/api/migration") {
			sessionID := r.URL.Query().Get("session")
			if sessionID != "" {
				log.Printf("%s %s %s %s (Protocol: %s, Session: %s)",
					emoji, r.RemoteAddr, r.Method, r.URL.Path, protocol, sessionID)
			} else {
				log.Printf("%s %s %s %s (Protocol: %s, NEW SESSION)",
					emoji, r.RemoteAddr, r.Method, r.URL.Path, protocol)
			}
		} else {
			log.Printf("%s %s %s %s (Protocol: %s)",
				emoji, r.RemoteAddr, r.Method, r.URL.Path, protocol)
		}

		// Set Alt-Svc header for HTTP/3 advertisement
		w.Header().Set("Alt-Svc", `h3=":9444"; ma=86400`)
		w.Header().Set("X-Server-Protocol", r.Proto)

		mux.ServeHTTP(w, r)
	})

	// TLS configuration optimized for connection migration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		NextProtos: []string{"h3", "h2", "http/1.1"},
		ClientAuth: tls.NoClientCert,
		// Enable session resumption to help with connection migration
		SessionTicketsDisabled: false,
		ClientSessionCache:     tls.NewLRUClientSessionCache(256),
	}

	// Start HTTP/1.1 & HTTP/2 server (TCP)
	go func() {
		tcpServer := &http.Server{
			Addr:         ":9444",
			Handler:      loggedMux,
			TLSConfig:    tlsConfig,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}

		log.Println("Starting HTTP/1.1 & HTTP/2 server (TCP) on :9444")
		if err := tcpServer.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Printf("TCP server error: %v", err)
		}
	}()

	// Start HTTP/1.1 server (no TLS) for testing
	go func() {
		httpServer := &http.Server{
			Addr:    ":8080",
			Handler: loggedMux,
		}
		log.Println("Starting HTTP/1.1 server (no TLS) on :8080")
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Give servers time to start
	time.Sleep(2 * time.Second)

	// Create HTTP/3 server (UDP)
	h3Server := &http3.Server{
		Addr:      ":9444",
		Handler:   loggedMux,
		TLSConfig: tlsConfig,
	}

	// Print startup information
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("HTTP/3 Server Ready on Windows!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Test URLs:")
	fmt.Println("   https://localhost:9444        (HTTP/3 + HTTP/2)")
	fmt.Println("   http://localhost:8080         (HTTP/1.1 only)")
	fmt.Println("   https://localhost:9444/api/test (API endpoint)")
	fmt.Println("")
	fmt.Println("Windows Testing Commands:")
	fmt.Println("   curl.exe -v --http3-only -k https://localhost:9444/api/test")
	fmt.Println("   curl.exe -v --http2 -k https://localhost:9444/api/test")
	fmt.Println("   curl.exe -v http://localhost:8080/api/test")
	fmt.Println("")
	fmt.Println("Browser Setup:")
	fmt.Println("   Chrome: chrome://flags/#enable-quic (enable)")
	fmt.Println("   Edge: edge://flags/#enable-quic (enable)")
	fmt.Println("   Then restart browser and visit https://localhost:9444")
	fmt.Println("")
	fmt.Println("Tips for Windows:")
	fmt.Println("   - Allow through Windows Firewall if prompted")
	fmt.Println("   - Look for rocket emoji in logs = HTTP/3 working!")
	fmt.Println("   - Try refreshing page multiple times")
	fmt.Println("   - HTTP/3 may take a few requests to activate")
	fmt.Println(strings.Repeat("=", 60))

	log.Println("Starting HTTP/3 server (UDP) on :9444...")
	err := h3Server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatal("HTTP/3 server failed to start:", err)
	}
}
