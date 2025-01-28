package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var redisClient *redis.Client

func SetUpRedisConnection(ctx context.Context) (*redis.Client, error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:" + os.Getenv("REDIS_PORT"),
	})

	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := redisClient.Ping(timeoutCtx).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	return redisClient, nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}
