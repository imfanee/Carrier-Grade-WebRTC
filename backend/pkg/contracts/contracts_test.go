// Package contracts — Tests for contract types and interfaces.
//
// By:- Faisal Hanif | imfanee@gmail.com

package contracts

import (
	"context"
	"testing"
	"time"
)

func TestClaimsFields(t *testing.T) {
	claims := &Claims{
		Subject:   "user-1",
		SessionID: "sess-abc",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	if claims.Subject != "user-1" {
		t.Errorf("expected Subject user-1, got %s", claims.Subject)
	}
	if claims.SessionID != "sess-abc" {
		t.Errorf("expected SessionID sess-abc, got %s", claims.SessionID)
	}
}

func TestHealthStatus(t *testing.T) {
	healthy := HealthStatus{Healthy: true, Reason: ""}
	if !healthy.Healthy {
		t.Error("expected healthy to be true")
	}
	unhealthy := HealthStatus{Healthy: false, Reason: "redis down"}
	if unhealthy.Healthy {
		t.Error("expected unhealthy to be false")
	}
}

func TestSessionStoreInterface(t *testing.T) {
	// Compile-time check that noop implementations satisfy the interface.
	var _ SessionStore = (*noopSessionStore)(nil)
}

type noopSessionStore struct{}

func (n *noopSessionStore) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return nil
}
func (n *noopSessionStore) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, nil
}
func (n *noopSessionStore) Delete(_ context.Context, _ string) error {
	return nil
}
