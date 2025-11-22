package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Maiar0/api-sqlite-base-go/auth"
	"github.com/Maiar0/api-sqlite-base-go/server"
)

func main() {
	loadDotEnv()
	//setup
	auth.InitJWTSecret()
	//setup handles
	mux := http.NewServeMux()
	auth.Register(mux)

	// Serve files inside ./tests/ under /tests/
	fileServer := http.FileServer(http.Dir("./tests"))
	mux.Handle("/tests/", http.StripPrefix("/tests/", fileServer))

	mux.HandleFunc("/ws/echo", server.HandleEchoWS)

	server.Run(mux, ":3000")
}

// loadDotEnv loads key=value pairs from a local .env file so JWT_SECRET is
// available during local development without exporting it manually.
func loadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("warning: could not read .env: %v", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("warning: skipping malformed .env line: %q", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}
		if err := os.Setenv(key, val); err != nil {
			log.Printf("warning: could not set env %s: %v", key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("warning: error reading .env: %v", err)
	}
}
