// Package contracts — HealthChecker interface for dependency health probes.
//
// By:- Faisal Hanif | imfanee@gmail.com

package contracts

import "context"

// HealthStatus represents the health state of a component.
type HealthStatus struct {
	Healthy bool
	Reason  string
}

// HealthChecker checks the health of a dependency (e.g. Redis, Auth).
type HealthChecker interface {
	Check(ctx context.Context) HealthStatus
}
