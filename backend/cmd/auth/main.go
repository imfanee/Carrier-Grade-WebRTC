// Command auth — Auth stub service for token issuance and validation.
//
// Provides JWT issuance, validation endpoint, and health checks.
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/faisalhanif/carrier-grade-webrtc/internal/auth"
)

const defaultPort = "8081"
const defaultSecret = "carrier-grade-webrtc-secret-change-in-production"

func main() {
	port := getEnv("AUTH_PORT", defaultPort)
	secret := getEnv("AUTH_SECRET", defaultSecret)
	validator := auth.NewJWTValidator(secret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/token", handleIssueToken(secret))
	mux.HandleFunc("GET /auth/validate", handleValidate(validator))
	mux.HandleFunc("GET /health/live", handleLiveness)
	mux.HandleFunc("GET /health/ready", handleReadiness)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Auth service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func handleIssueToken(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID string `json:"userId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
			http.Error(w, `{"error":"userId required"}`, http.StatusBadRequest)
			return
		}
		claims := jwt.MapClaims{
			"sub":         req.UserID,
			"session_id":  generateSessionID(),
			"exp":         time.Now().Add(24 * time.Hour).Unix(),
			"iat":         time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": signed})
	}
}

func handleValidate(validator *auth.JWTValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("Authorization")
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}
		if token == "" {
			http.Error(w, `{"error":"token required"}`, http.StatusUnauthorized)
			return
		}
		claims, err := validator.Validate(r.Context(), token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sub":         claims.Subject,
			"session_id":  claims.SessionID,
			"expires_at":  claims.ExpiresAt,
		})
	}
}

func handleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {
	// Auth has no external dependencies in stub; always ready.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + randomHex(8)
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[rand.Intn(16)]
	}
	return string(b)
}
