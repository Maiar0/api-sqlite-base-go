package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Maiar0/api-sqlite-base-go/auth"
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

var (
	clients   = make(map[string]*Client)
	clientsMu sync.RWMutex
)

func reader(c *Client) error {
	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Read error: %v", err)
			return err
		}
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Printf("[WS] Unmarshal error: %v", err)
			Send(c.Conn, Message{}, err, msgType)
			continue
		}
		//c.ID is the users UUID
		//TODO::logic? route on m.Action( join game, send move, etc)

		Send(c.Conn, m, nil, msgType)
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
	//authenticate
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok || claims == nil { //TODO:: This shouldnt happen?
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userUUID := claims.UserUUID

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
	}
	log.Printf("[WS] client %s conected from %s", userUUID, r.RemoteAddr)
	//create client and store
	client := &Client{
		ID:   userUUID,
		Conn: conn,
	}

	clientsMu.Lock()
	clients[userUUID] = client
	clientsMu.Unlock()

	//clena up user and connection on exit
	defer func() {
		clientsMu.Lock()
		delete(clients, userUUID)
		clientsMu.Unlock()
		conn.Close()
		log.Printf("[WS] client %s disconnected", userUUID)
	}()
	//create connected message
	connectedMsg := Message{
		Action: "connected",
		Data: map[string]string{
			"uuid": userUUID,
		},
	}
	//send connected message
	Send(conn, connectedMsg, nil, websocket.TextMessage)

	if err := reader(client); err != nil {
		//read logged errors already
		return
	}

}
