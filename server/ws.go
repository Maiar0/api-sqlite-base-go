package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//handle cors at conection level
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default
		return true
	},
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var clients = make(map[string]*Client)

func reader(conn *websocket.Conn, target any) error {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Read error: %v", err)
			return err
		}
		if err := json.Unmarshal(msg, target); err != nil {
			log.Printf("[WS] Unmarshal error: %v", err)
			return err
		}
		//add client to list
		//echo message back to client
		respond(conn, msgType, target, nil)
	}
}

func respond(conn *websocket.Conn, msgType int, data interface{}, errMsg error) {
	m := Message{
		Data: data,
	}
	if errMsg != nil {
		m.Err = errMsg.Error()
	}

	out, _ := json.Marshal(m) //ignoring error for now
	conn.WriteMessage(msgType, out)
}

func HandleEchoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
	}
	defer conn.Close()
	log.Printf("[WS]client conected from %s", r.RemoteAddr)
	var m Message
	reader(conn, &m)

	log.Printf("[WS] client disconnected from %s", r.RemoteAddr)
}
