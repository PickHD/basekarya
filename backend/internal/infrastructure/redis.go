package infrastructure

import (
	"context"
	"hris-backend/internal/config"
	"hris-backend/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClientProvider struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) *RedisClientProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Errorw("Failed to connect to redis:", err)
	}

	logger.Info("Connected to Redis")

	return &RedisClientProvider{Client: rdb}
}

func (r *RedisClientProvider) Close() error {
	return r.Client.Close()
}

func (r *RedisClientProvider) GetClient() *redis.Client {
	return r.Client
}
