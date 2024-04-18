package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Redis struct {
	Client *redis.Client
}

func New() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	slog.Info("redis success configured")

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("oh no. redis not send PONG (" + pong + ")")
		panic(err)
	}

	return &Redis{Client: client}, nil
}
