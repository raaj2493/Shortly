package repository

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(addr , password string){
	RedisClient = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password, // empty if no password
        DB:       0,        // default database
    })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    fmt.Println("Redis connection established")
}