package redis

import (
	"context"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedisClient() *redis.Client {
	redis_uri := os.Getenv("REDIS_URI")
	if redis_uri == "" {
		fmt.Printf("empty redis uri\n")
		return nil
	}

	options, err := redis.ParseURL(redis_uri)
	if err != nil {
		fmt.Printf("redis parse url failed\n")
		return nil
	}

	redisClient = redis.NewClient(options)

	// check connection
	if err = redisClient.Ping(context.Background()).Err(); err != nil {
		fmt.Printf("connect to redis failed\n")
		return nil
	}

	return redisClient
}
