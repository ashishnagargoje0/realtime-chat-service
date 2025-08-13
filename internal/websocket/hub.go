package websocket

import "github.com/redis/go-redis/v9"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	RedisClient *redis.Client // Redis client for Pub/Sub
	Channel     string        // Redis channel name
}

// NewHub creates a new Hub with optional Redis client and channel
func NewHub(rdb *redis.Client, channel string) *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		RedisClient: rdb,
		Channel:     channel,
	}
}

// Run listens for register, unregister, and broadcast events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
