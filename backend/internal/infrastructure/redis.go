package infrastructure

import (
	"basekarya-backend/internal/config"
	"basekarya-backend/pkg/logger"
	"context"
	"time"

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

func (r *RedisClientProvider) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClientProvider) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClientProvider) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
