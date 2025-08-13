package websocket

import (
	"log"
	"realtime-chat-service/internal/redis"
	"realtime-chat-service/internal/storage"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	Username string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Prepend username
		fullMessage := []byte(c.Username + ": " + string(message))

		// Broadcast to all clients
		c.Hub.Broadcast <- fullMessage

		// Publish to Redis channel
		redis.PublishMessage(c.Hub.RedisClient, c.Hub.Channel, string(fullMessage))

		// Save to Postgres
		err = storage.SaveMessage(c.Username, string(message))
		if err != nil {
			log.Println("Failed to save message:", err)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
