// Package contracts — TokenValidator interface for JWT validation.
//
// By:- Faisal Hanif | imfanee@gmail.com

package contracts

import "context"

// Claims represents validated JWT claims.
type Claims struct {
	Subject   string
	SessionID string
	ExpiresAt int64
}

// TokenValidator validates JWTs and returns claims or an error.
type TokenValidator interface {
	Validate(ctx context.Context, token string) (*Claims, error)
}
