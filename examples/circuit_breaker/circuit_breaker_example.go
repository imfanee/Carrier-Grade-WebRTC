// Package main — Circuit breaker implementation example.
//
// Prevents cascading failures when calling downstream services (e.g. Auth, Redis).
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sony/gobreaker"
)

// DownstreamClient simulates a client that can fail.
type DownstreamClient struct {
	failRate float64 // 0–1, probability of failure
}

func (d *DownstreamClient) Call(ctx context.Context) error {
	// Simulate occasional failures
	if d.failRate > 0 && time.Now().UnixNano()%100 < int64(d.failRate*100) {
		return fmt.Errorf("downstream unavailable")
	}
	return nil
}

func main() {
	settings := gobreaker.Settings{
		Name:        "auth-service",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("circuit %s: %s -> %s", name, from, to)
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	client := &DownstreamClient{failRate: 0.5}

	for i := 0; i < 10; i++ {
		result, err := cb.Execute(func() (interface{}, error) {
			return nil, client.Call(context.Background())
		})
		if err != nil {
			if err == gobreaker.ErrOpenState {
				log.Printf("request %d: circuit open, skipping", i+1)
			} else {
				log.Printf("request %d: %v", i+1, err)
			}
			continue
		}
		_ = result
		log.Printf("request %d: success", i+1)
		time.Sleep(100 * time.Millisecond)
	}
}
