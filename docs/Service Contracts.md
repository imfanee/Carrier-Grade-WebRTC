# Service Contracts

<!--
  API and interface contracts for the Carrier-Grade WebRTC microservices.
  Defines REST, WebSocket, and Go interface boundaries.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## REST Endpoints

### Auth Service

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/token` | Issue JWT for valid credentials |
| GET | `/auth/validate` | Validate JWT; returns claims or 401 |
| GET | `/health/live` | Liveness probe |
| GET | `/health/ready` | Readiness (e.g. no dependencies) |

### Signaling Service

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health/live` | Liveness probe |
| GET | `/health/ready` | Readiness (Redis, Auth connectivity) |
| WS | `/ws/signal` | WebSocket signaling (query: `?token=<jwt>`) |

## WebSocket Signaling Protocol

All messages are JSON. Direction: Client → Server (C2S) or Server → Client (S2C).

| Type | Direction | Payload | Description |
|------|-----------|---------|-------------|
| `join` | C2S | `{ "roomId": string }` | Join a room |
| `joined` | S2C | `{ "roomId": string, "peerId": string }` | Confirmation |
| `offer` | C2S | `{ "peerId": string, "sdp": string }` | SDP offer |
| `offer` | S2C | `{ "peerId": string, "sdp": string }` | Relay offer to peer |
| `answer` | C2S | `{ "peerId": string, "sdp": string }` | SDP answer |
| `answer` | S2C | `{ "peerId": string, "sdp": string }` | Relay answer |
| `ice-candidate` | C2S/S2C | `{ "peerId": string, "candidate": object }` | ICE candidate |
| `leave` | C2S | `{ "roomId": string }` | Leave room |
| `error` | S2C | `{ "code": string, "message": string }` | Error notification |

## Go Interface Definitions

See `backend/pkg/contracts/` for the canonical definitions. Summary:

```go
// SessionStore — Redis-backed session persistence
type SessionStore interface {
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Get(ctx context.Context, key string) ([]byte, error)
    Delete(ctx context.Context, key string) error
}

// TokenValidator — JWT validation
type TokenValidator interface {
    Validate(ctx context.Context, token string) (*Claims, error)
}

// HealthChecker — Dependency health
type HealthChecker interface {
    Check(ctx context.Context) HealthStatus
}
```

---

*By:- Faisal Hanif | imfanee@gmail.com*
