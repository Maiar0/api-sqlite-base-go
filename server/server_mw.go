package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
)

func Run(mux *http.ServeMux, port string) {
	if port == "" {
		port = ":3000"
	}
	log.Printf("Server running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, simpleLogger(corsMiddleware(mux))))
}

func simpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Detect if this request *intends* to do a WebSocket upgrade
		isWS := strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") &&
			strings.EqualFold(r.Header.Get("Upgrade"), "websocket")

		if isWS {
			log.Printf("[REQ][WS? yes] %s %s (Origin=%q)", r.Method, r.URL.Path, r.Header.Get("Origin"))
		} else {
			log.Printf("[REQ] %s %s", r.Method, r.URL.Path)
		}

		// Call downstream handler
		next.ServeHTTP(w, r)

		dur := time.Since(start)
		if isWS {
			// For WS, this logs when the connection closes
			log.Printf("[DONE][WS] %s %s (%v)", r.Method, r.URL.Path, dur)
		} else {
			log.Printf("[DONE] %s %s (%v)", r.Method, r.URL.Path, dur)
		}
	})
}

func corsMiddleware(mux http.Handler) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://app.example.com",
			"http://localhost:5173", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true, // if you need cookies/auth
		MaxAge:           600,  // seconds
	})

	return c.Handler(mux)
}
