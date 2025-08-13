package main

import (
	"fmt"
	"log"
	"net/http"

	"realtime-chat-service/config"
	"realtime-chat-service/internal/redis"
	"realtime-chat-service/internal/storage"
	"realtime-chat-service/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Postgres DB
	storage.InitDB()

	// Initialize Redis
	rdb := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)

	// Initialize WebSocket Hub with Redis
	hub := websocket.NewHub(rdb, cfg.RedisChannel)
	go hub.Run()

	// Subscribe to Redis channel to broadcast messages from other servers
	go redis.SubscribeMessages(rdb, cfg.RedisChannel, func(msg string) {
		hub.Broadcast <- []byte(msg)
	})

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWSWithHistory(hub, w, r) // sends last 10 messages to new clients
	})

	// Serve frontend (HTML + JS)
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	log.Printf("Server started on :%s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), nil))
}
