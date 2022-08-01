package initiator

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sso/platform/logger"
)

func InitCache(url string, log logger.Logger) *redis.Client {
	// TODO implement
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to parse redis url: %v", err))
	}

	client := redis.NewClient(opts)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to ping redis: %v", err))
	}

	return client
}
