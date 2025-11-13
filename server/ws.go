package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	//handle cors at conection level
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

func HandleEchoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
	}
	defer conn.Close()
	log.Printf("[WS]client conected from %s", r.RemoteAddr)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Read error: %v", err)
			break
		}
		log.Printf("[WS] Received message: %s", string(msg))
		//echo msg
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Printf("[WS] Write error: %v", err)
			break
		}
	}
	log.Printf("[WS] client disconnected from %s", r.RemoteAddr)
}
