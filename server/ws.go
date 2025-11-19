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
	Action string      `json:"action,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Err    string      `json:"err,omitempty"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var clients = make(map[string]*Client)

func reader(conn *websocket.Conn) error {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Read error: %v", err)
			return err
		}
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Printf("[WS] Unmarshal error: %v", err)
			Send(conn, Message{}, err, msgType)
			continue
		}
		//add client to list
		//echo message back to client
		Send(conn, m, nil, msgType)
	}
}

func Send(conn *websocket.Conn, m Message, errMsg error, msgType int) {
	if errMsg != nil {
		m.Err = errMsg.Error()
	}

	out, err := json.Marshal(m)
	if err != nil {
		log.Printf("[WS] Marshal error: %v", err)
		return
	}

	if err := conn.WriteMessage(msgType, out); err != nil {
		log.Printf("[WS] Write error: %v", err)
	}
}

func HandleEchoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
	}
	defer conn.Close()
	log.Printf("[WS]client conected from %s", r.RemoteAddr)

	reader(conn)

	log.Printf("[WS] client disconnected from %s", r.RemoteAddr)
}
