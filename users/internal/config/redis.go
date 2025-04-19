package config

import (
	"context"
	"fmt"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis () {
	fmt.Println("REDIS_ADDR",os.Getenv("REDIS_ADDR"), "REDIS_PASWORD", os.Getenv("REDIS_PASSWORD"))
	Client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Unable to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis correctly")	
}