// helpers/redis_handler.go
package helpers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisHandler defines methods for Redis operations
type RedisHandler interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}

// RedisHelper is the concrete implementation of RedisHandler
type RedisHelper struct {
	client *redis.Client
}

// NewRedisHelper creates a new RedisHelper
func NewRedisHelper(client *redis.Client) RedisHandler {
	return &RedisHelper{
		client: client,
	}
}

func (rh *RedisHelper) Get(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := rh.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return err
	}

	return nil
}

func (rh *RedisHelper) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rh.client.Set(ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rh *RedisHelper) Delete(key string) error {
	ctx := context.Background()
	err := rh.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
