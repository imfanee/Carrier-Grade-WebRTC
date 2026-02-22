// Package main — Non-blocking telemetry module example.
//
// Demonstrates async metrics and tracing without blocking request paths.
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// TelemetryEvent represents a single telemetry event (metric, span, log).
type TelemetryEvent struct {
	Type      string
	Name      string
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
}

// NonBlockingTelemetry buffers events and flushes asynchronously.
type NonBlockingTelemetry struct {
	buffer   chan TelemetryEvent
	flushFn  func(events []TelemetryEvent)
	stopCh   chan struct{}
	wg       sync.WaitGroup
	batchSize int
	interval  time.Duration
}

// NewNonBlockingTelemetry creates a telemetry collector that never blocks callers.
func NewNonBlockingTelemetry(bufferSize, batchSize int, flushInterval time.Duration, flushFn func([]TelemetryEvent)) *NonBlockingTelemetry {
	t := &NonBlockingTelemetry{
		buffer:   make(chan TelemetryEvent, bufferSize),
		flushFn:  flushFn,
		stopCh:   make(chan struct{}),
		batchSize: batchSize,
		interval:  flushInterval,
	}
	t.wg.Add(1)
	go t.flushLoop()
	return t
}

// Record enqueues an event without blocking. Drops if buffer is full.
func (t *NonBlockingTelemetry) Record(event TelemetryEvent) {
	event.Timestamp = time.Now()
	select {
	case t.buffer <- event:
	default:
		// Buffer full; drop to avoid blocking (or use metrics for drops)
	}
}

// RecordMetric is a convenience for counter/gauge metrics.
func (t *NonBlockingTelemetry) RecordMetric(name string, value float64, labels map[string]string) {
	t.Record(TelemetryEvent{
		Type:   "metric",
		Name:   name,
		Value:  value,
		Labels: labels,
	})
}

// RecordSpan records a span (start/end) for tracing.
func (t *NonBlockingTelemetry) RecordSpan(name string, duration time.Duration, labels map[string]string) {
	t.Record(TelemetryEvent{
		Type:   "span",
		Name:   name,
		Value:  duration.Seconds(),
		Labels: labels,
	})
}

// Shutdown gracefully stops the telemetry and flushes remaining events.
func (t *NonBlockingTelemetry) Shutdown(ctx context.Context) error {
	close(t.stopCh)
	t.wg.Wait()
	return nil
}

func (t *NonBlockingTelemetry) flushLoop() {
	defer t.wg.Done()
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	batch := make([]TelemetryEvent, 0, t.batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		events := make([]TelemetryEvent, len(batch))
		copy(events, batch)
		batch = batch[:0]
		go t.flushFn(events) // Flush in goroutine to avoid blocking loop
	}

	for {
		select {
		case <-t.stopCh:
			for len(batch) < cap(batch) {
				select {
				case e := <-t.buffer:
					batch = append(batch, e)
				default:
					goto done
				}
			}
		done:
			flush()
			return
		case e := <-t.buffer:
			batch = append(batch, e)
			if len(batch) >= t.batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func main() {
	telemetry := NewNonBlockingTelemetry(1000, 10, 2*time.Second, func(events []TelemetryEvent) {
		for _, e := range events {
			log.Printf("telemetry: %s %s %.2f %v", e.Type, e.Name, e.Value, e.Labels)
		}
	})
	defer telemetry.Shutdown(context.Background())

	// Simulate request handling
	start := time.Now()
	telemetry.RecordMetric("request_total", 1, map[string]string{"path": "/ws/signal"})
	time.Sleep(50 * time.Millisecond)
	telemetry.RecordSpan("handle_request", time.Since(start), map[string]string{"path": "/ws/signal"})

	time.Sleep(3 * time.Second) // Allow flush
}
