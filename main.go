package main

import (
	"net/http"

	server "github.com/Maiar0/api-sqlite-base-go/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	})

	mux.HandleFunc("/ws/echo", server.HandleEchoWS)

	server.Run(mux, ":3000")
}
