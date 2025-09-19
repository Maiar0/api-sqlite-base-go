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

	server.Run(mux, ":3000")
}
