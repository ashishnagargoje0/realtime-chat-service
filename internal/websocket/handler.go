package websocket

import (
	"net/http"
	"realtime-chat-service/internal/storage"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWSWithHistory sends last 10 messages to new clients
func ServeWSWithHistory(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	hub.Register <- client

	// Send last 10 messages
	messages, err := storage.GetLastMessages(10)
	if err == nil {
		for _, m := range messages {
			client.Send <- []byte(m.Username + ": " + m.Content)
		}
	}

	go client.WritePump()
	go client.ReadPump()
}
