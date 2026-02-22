// Package contracts — SessionStore interface for Redis-backed session persistence.
//
// By:- Faisal Hanif | imfanee@gmail.com

package contracts

import (
	"context"
	"time"
)

// SessionStore persists session and signaling state with TTL.
type SessionStore interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
