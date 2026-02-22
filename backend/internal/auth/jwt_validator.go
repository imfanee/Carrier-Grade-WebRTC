// Package auth — JWT validation implementation for TokenValidator contract.
//
// By:- Faisal Hanif | imfanee@gmail.com

package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/faisalhanif/carrier-grade-webrtc/pkg/contracts"
)

// JWTValidator implements contracts.TokenValidator.
type JWTValidator struct {
	secretKey []byte
}

// NewJWTValidator creates a validator with the given secret.
func NewJWTValidator(secretKey string) *JWTValidator {
	return &JWTValidator{secretKey: []byte(secretKey)}
}

// Validate parses and validates a JWT, returning claims or an error.
func (j *JWTValidator) Validate(ctx context.Context, tokenString string) (*contracts.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	exp, _ := claims["exp"].(float64)
	sub, _ := claims["sub"].(string)
	sid, _ := claims["session_id"].(string)
	if sub == "" {
		return nil, errors.New("missing subject")
	}
	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}
	return &contracts.Claims{
		Subject:   sub,
		SessionID: sid,
		ExpiresAt: int64(exp),
	}, nil
}
