// Package main — Load shedding example for graceful degradation under load.
//
// Rejects new connections when system resources exceed thresholds.
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"net/http"
	"sync/atomic"
	"time"
)

// LoadShedder rejects requests when load exceeds a threshold.
type LoadShedder struct {
	activeConnections int64
	maxConnections    int64
	cpuUsagePercent  int64 // Simulated; in production use actual metrics
	maxCPUPercent    int64
}

// NewLoadShedder creates a load shedder with the given limits.
func NewLoadShedder(maxConnections, maxCPUPercent int64) *LoadShedder {
	return &LoadShedder{
		maxConnections: maxConnections,
		maxCPUPercent: maxCPUPercent,
	}
}

// Allow returns true if the request should be accepted.
func (l *LoadShedder) Allow() bool {
	active := atomic.LoadInt64(&l.activeConnections)
	cpu := atomic.LoadInt64(&l.cpuUsagePercent)
	if active >= l.maxConnections {
		return false
	}
	if cpu >= l.maxCPUPercent {
		return false
	}
	return true
}

// Acquire increments active connections; call Release when done.
func (l *LoadShedder) Acquire() bool {
	if !l.Allow() {
		return false
	}
	atomic.AddInt64(&l.activeConnections, 1)
	return true
}

// Release decrements active connections.
func (l *LoadShedder) Release() {
	atomic.AddInt64(&l.activeConnections, -1)
}

// SetCPUUsage updates the simulated CPU usage (for demo).
func (l *LoadShedder) SetCPUUsage(percent int64) {
	atomic.StoreInt64(&l.cpuUsagePercent, percent)
}

// Middleware wraps an HTTP handler with load shedding.
func (l *LoadShedder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Acquire() {
			http.Error(w, `{"error":"service overloaded"}`, http.StatusServiceUnavailable)
			return
		}
		defer l.Release()
		next.ServeHTTP(w, r)
	})
}

func main() {
	shedder := NewLoadShedder(5, 90)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":9090",
		Handler: shedder.Middleware(mux),
	}
	_ = server
	// server.ListenAndServe()
}
