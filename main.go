package main

import (
	"net/http"

	"github.com/Maiar0/api-sqlite-base-go/auth"
	"github.com/Maiar0/api-sqlite-base-go/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	})
	auth.Register(mux)
	// Serve files inside ./tests/ under /tests/
	fileServer := http.FileServer(http.Dir("./tests"))
	mux.Handle("/tests/", http.StripPrefix("/tests/", fileServer))

	mux.HandleFunc("/ws/echo", server.HandleEchoWS)

	server.Run(mux, ":3000")
}
