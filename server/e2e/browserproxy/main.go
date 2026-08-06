package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	firstBasePath  = "/instance-a"
	secondBasePath = "/instance-b"
	firstTarget    = "http://127.0.0.1:4174"
	secondTarget   = "http://127.0.0.1:4175"
)

func main() {
	firstURL := mustParseURL(firstTarget)
	secondURL := mustParseURL(secondTarget)
	firstProxy := httputil.NewSingleHostReverseProxy(firstURL)
	secondProxy := httputil.NewSingleHostReverseProxy(secondURL)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(firstURL, secondURL))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case hasPathPrefix(r.URL.Path, firstBasePath):
			firstProxy.ServeHTTP(w, r)
		case hasPathPrefix(r.URL.Path, secondBasePath):
			secondProxy.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	server := &http.Server{
		Addr:              "127.0.0.1:" + envOrDefault("TEST_APP_PORT", "4173"),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("browser test proxy: %v", err)
	}
}

func healthHandler(targets ...*url.URL) http.HandlerFunc {
	client := &http.Client{Timeout: time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, target := range targets {
			request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, target.String()+"/health", nil)
			if err != nil {
				http.Error(w, "invalid health request", http.StatusInternalServerError)
				return
			}
			response, err := client.Do(request)
			if err != nil {
				http.Error(w, "test instance unavailable", http.StatusServiceUnavailable)
				return
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				log.Printf("close health response: %v", closeErr)
			}
			if response.StatusCode != http.StatusOK {
				http.Error(w, "test instance unhealthy", http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func hasPathPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func mustParseURL(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		panic(fmt.Sprintf("parse browser test target: %v", err))
	}
	return parsed
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
