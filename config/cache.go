package config

import (
	"context"
	"user_api/utils/logger"

	"github.com/go-redis/redis/v8"
)

func NewCache(ctx context.Context, cfg Configuration) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		logger.GetLogger(ctx).Errorf("failed to connect redis: %v", err)
		return nil, err
	}

	return rdb, nil
}
