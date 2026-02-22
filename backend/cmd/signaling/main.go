// Command signaling — WebRTC signaling service over WebSocket.
//
// Handles SDP/ICE relay, room management, and JWT-authenticated connections.
// By:- Faisal Hanif | imfanee@gmail.com

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/faisalhanif/carrier-grade-webrtc/internal/auth"
	"github.com/faisalhanif/carrier-grade-webrtc/internal/cache"
	"github.com/faisalhanif/carrier-grade-webrtc/internal/signaling/hub"
	"github.com/faisalhanif/carrier-grade-webrtc/pkg/contracts"
)

const defaultPort = "8080"
const defaultRedisAddr = "localhost:6379"
const defaultSecret = "carrier-grade-webrtc-secret-change-in-production"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	port := getEnv("SIGNALING_PORT", defaultPort)
	redisAddr := getEnv("REDIS_ADDR", defaultRedisAddr)
	secret := getEnv("AUTH_SECRET", defaultSecret)

	var store contracts.SessionStore
	var redisStore *cache.RedisStore
	redisStore, err := cache.NewRedisStore(redisAddr, "", 0)
	if err != nil {
		log.Printf("Redis unavailable, using in-memory fallback: %v", err)
		store = &hub.NoopStore{}
		redisStore = nil
	} else {
		store = redisStore
	}

	validator := auth.NewJWTValidator(secret)
	signalHub := hub.NewSignalHub(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws/signal", handleWebSocket(signalHub, validator))
	mux.HandleFunc("GET /health/live", handleLiveness)
	mux.HandleFunc("GET /health/ready", handleReadiness(redisStore))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Printf("Signaling service listening on :%s", port)
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

func handleWebSocket(signalHub *hub.SignalHub, validator *auth.JWTValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}
		claims, err := validator.Validate(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		peerID := claims.Subject + "-" + claims.SessionID
		signalHub.Register(peerID, conn)
		defer signalHub.Unregister(peerID)

		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var msg hub.SignalMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				sendError(conn, "invalid message")
				continue
			}
			signalHub.HandleMessage(peerID, msg)
		}
	}
}

func sendError(conn *websocket.Conn, message string) {
	conn.WriteJSON(map[string]string{"type": "error", "message": message})
}

func handleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleReadiness(store *cache.RedisStore) http.HandlerFunc { // nil when Redis unavailable
	return func(w http.ResponseWriter, r *http.Request) {
		if store == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := store.Client().Ping(ctx).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("redis unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
