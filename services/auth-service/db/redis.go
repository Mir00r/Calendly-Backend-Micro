package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection established.")
}
