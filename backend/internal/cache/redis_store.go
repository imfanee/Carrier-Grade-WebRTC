// Package cache — Redis-backed implementation of SessionStore.
//
// By:- Faisal Hanif | imfanee@gmail.com

package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/faisalhanif/carrier-grade-webrtc/pkg/contracts"
)

// RedisStore implements contracts.SessionStore using Redis.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a Redis-backed session store.
func NewRedisStore(addr string, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

// Set stores a value with TTL.
func (r *RedisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key.
func (r *RedisStore) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// Delete removes a key.
func (r *RedisStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Client returns the underlying Redis client for health checks.
func (r *RedisStore) Client() *redis.Client {
	return r.client
}

// Ensure RedisStore implements contracts.SessionStore.
var _ contracts.SessionStore = (*RedisStore)(nil)
