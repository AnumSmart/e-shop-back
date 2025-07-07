package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(addr, password string, db int) *RedisRepo {
	return &RedisRepo{
		redis.NewClient(
			&redis.Options{
				Addr:     addr,
				Password: password,
				DB:       db,
			},
		),
	}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisRepo) Exists(ctx context.Context, redisKey string) (bool, error) {
	result, err := r.client.Exists(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	// Возвращает true, если ключ существует (1), false если нет (0)
	return result == 1, nil
}
