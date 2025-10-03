package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go/http3"
)

func mustAbs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}

func main() {
	addr := ":443"
	publicDir := "./public"
	certFile := "localhost+3.pem"    // adjust
	keyFile := "localhost+3-key.pem" // adjust

	// Check cert + key
	if _, err := os.Stat(certFile); err != nil {
		log.Fatalf("cert %s missing: %v", certFile, err)
	}
	if _, err := os.Stat(keyFile); err != nil {
		log.Fatalf("key %s missing: %v", keyFile, err)
	}
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		log.Fatalf("public dir %s missing", publicDir)
	}

	// Always wrap response to inject Alt-Svc
	injectAltSvc := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This is the magic line browsers look for
			log.Println("Injecting Alt-Svc header")
			w.Header().Add("Alt-Svc", `h3=":443"; ma=86400`)
			h.ServeHTTP(w, r)
		})
	}

	// File server
	fileServer := http.FileServer(http.Dir(publicDir))
	mux := http.NewServeMux()
	mux.Handle("/", injectAltSvc(fileServer))

	// API test
	mux.Handle("/api/test", injectAltSvc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"protocol": "%s"}`, r.Proto)
		w.Write([]byte(resp))
	})))

	log.Printf("Starting HTTP/3 server on %s, serving %s", addr, mustAbs(publicDir))

	// Start HTTP/3
	if err := http3.ListenAndServeTLS(addr, certFile, keyFile, mux); err != nil {
		log.Fatal(err)
	}
}
