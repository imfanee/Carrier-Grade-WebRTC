// Package main — Observability example: structured logging and metrics.
//
// Demonstrates patterns for logs, metrics, and tracing in carrier-grade services.
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogEntry represents a structured log entry (JSON).
type LogEntry struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	TraceID   string            `json:"trace_id,omitempty"`
	SpanID    string            `json:"span_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// StructuredLogger writes JSON logs for log aggregation (e.g. ELK, Loki).
type StructuredLogger struct {
	logger *log.Logger
}

// NewStructuredLogger creates a logger that outputs JSON.
func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// Info logs an info-level message with optional fields.
func (s *StructuredLogger) Info(message string, fields map[string]interface{}) {
	s.emit("INFO", message, "", "", fields)
}

// Error logs an error-level message.
func (s *StructuredLogger) Error(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	s.emit("ERROR", message, "", "", fields)
}

// WithTrace logs with trace/span IDs for distributed tracing.
func (s *StructuredLogger) WithTrace(traceID, spanID string) *StructuredLogger {
	return &StructuredLogger{logger: s.logger}
}

func (s *StructuredLogger) emit(level, message, traceID, spanID string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   message,
		TraceID:   traceID,
		SpanID:    spanID,
		Fields:    fields,
	}
	data, _ := json.Marshal(entry)
	s.logger.Println(string(data))
}

// MetricsCollector simulates Prometheus-style metrics.
type MetricsCollector struct {
	counters map[string]int64
}

// IncrementCounter increments a counter (e.g. requests_total).
func (m *MetricsCollector) IncrementCounter(name string, labels map[string]string) {
	// In production: use Prometheus client
	_ = name
	_ = labels
}

// ObserveHistogram records a histogram value (e.g. request_duration_seconds).
func (m *MetricsCollector) ObserveHistogram(name string, value float64, labels map[string]string) {
	// In production: use Prometheus client
	_ = name
	_ = value
	_ = labels
}

func main() {
	logger := NewStructuredLogger()
	logger.Info("service started", map[string]interface{}{
		"port":   8080,
		"env":    "development",
	})
	logger.Info("request handled", map[string]interface{}{
		"path":     "/ws/signal",
		"duration": 0.05,
		"status":   200,
	})
}
