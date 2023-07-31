package kv

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var rdb *redis.Client

func init() {
	// Create a new Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Replace with your Redis server address
		Password: "",               // Replace with your Redis password
		DB:       0,                // Replace with your Redis database number
	})
	// Ping the Redis server to check the connection
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

}

func Set(ctx context.Context, key, val string, expireTime time.Duration) {
	// Set a value in Redis
	err := rdb.Set(ctx, key, val, expireTime).Err()
	if err != nil {
		log.Println("Failed to set value in Redis:", err)
	}
}

func Get(ctx context.Context, key string) *string {
	// Get the value from Redis
	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("Key does not exist")
		return nil
	} else if err != nil {
		fmt.Println("Failed to get value from Redis:", err)
		return nil
	} else {
		return &value
	}
}
