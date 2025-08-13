package config

import (
	"log"
	"os"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisChannel  string
	AppPort       string
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisChannel:  getEnv("REDIS_CHANNEL", "chatroom"),
		AppPort:       getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	log.Printf("Using default value for %s: %s", key, defaultVal)
	return defaultVal
}
