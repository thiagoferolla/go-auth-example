package cache

import (
	"os"
	"time"

	"github.com/go-redis/redis"
)

type RedisProvider struct {
	RedisClient *redis.Client
}

func NewRedisProvider() *RedisProvider {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":6379",
	})

	return &RedisProvider{redisClient}
}

func (provider RedisProvider) Get(key string) (string, error) {
	return provider.RedisClient.Get(key).Result()
}

func (provider RedisProvider) Set(key string, value string) error {
	return provider.RedisClient.Set(key, value, 0).Err()
}

func (provider RedisProvider) SetEx(key string, value string, expiration int) error {
	return provider.RedisClient.Set(key, value, time.Duration(expiration)).Err()
}
