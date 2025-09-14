package cache

import (
	"context"
	"time"

	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

// StatusCache определяет интерфейс для кэширования статусов уведомлений
type StatusCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
}

// RedisCache реализует StatusCache используя Redis
type RedisCache struct {
	client   *redis.Client
	strategy retry.Strategy
}

// NewRedisCache создает новый Redis кэш
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client:   client,
		strategy: retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 2},
	}
}

// Get получает значение из Redis кэша
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.GetWithRetry(ctx, c.strategy, key)
}

// Set сохраняет значение в Redis кэш с TTL
func (c *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.setWithRetry(ctx, key, value, ttl)
}

func (c *RedisCache) setWithRetry(ctx context.Context, id, value string, ttl time.Duration) error {
	return retry.Do(func() error {
		return c.client.Client.Set(ctx, id, value, ttl).Err()
	}, c.strategy)
}
