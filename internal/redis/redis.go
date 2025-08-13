package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

func SubscribeMessages(rdb *redis.Client, channel string, handler func(string)) {
	pubsub := rdb.Subscribe(ctx, channel)
	ch := pubsub.Channel()

	for msg := range ch {
		handler(msg.Payload)
	}
}

func PublishMessage(rdb *redis.Client, channel, message string) {
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		log.Println("Error publishing:", err)
	}
}
