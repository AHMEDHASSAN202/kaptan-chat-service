package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client
func NewRedisClient(fullUrl *string) *RedisClient {
	options, err := redis.ParseURL(*fullUrl)
	if err != nil {
		panic(err)
	}
	//// Configure TLS to skip certificate verification
	//options.TLSConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//}

	rdb := redis.NewClient(options)

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

func (r *RedisClient) Get(key string, dest interface{}) error {
	jsonData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonData), dest)
}

// Delete removes a key from the Redis store
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}

// Lock attempts to acquire a lock for a given key with an expiration time
func (r *RedisClient) Lock(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Unlock releases a lock by deleting the key
func (r *RedisClient) Unlock(key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis client connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// GetAllKeys retrieves all keys matching a given pattern
func (r *RedisClient) GetAllKeys(pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}
