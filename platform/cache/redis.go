package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/savioruz/roastgithub-api/pkg/utils"
	"os"
	"strconv"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

// NewRedisConnection initializes a Redis connection and returns a RedisClient instance
func NewRedisConnection() (*RedisClient, error) {
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))

	redisUrl, err := utils.ConnectionURLBuilder("redis")
	if err != nil {
		return nil, err
	}

	opt := &redis.Options{
		Addr:     redisUrl,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}

	// Initialize the Redis client
	client := redis.NewClient(opt)

	return &RedisClient{client: client}, nil
}

// Set stores a key-value pair in Redis with an expiration time
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
