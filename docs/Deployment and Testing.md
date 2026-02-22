# Deployment and Testing Guide

<!--
  Step-by-step deployment and validation for the Carrier-Grade WebRTC reference.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Prerequisites

- Go 1.21+
- Node.js 18+
- Docker and Docker Compose (for Redis)
- Two browser tabs or devices for testing calls

## Step 1: Start Redis

```bash
docker-compose up -d redis
```

Verify Redis is running:

```bash
docker-compose ps
redis-cli ping
# Should return PONG
```

## Step 2: Start Auth Service

```bash
cd backend
go mod download
go run ./cmd/auth
```

Auth listens on `http://localhost:8081`. Verify:

```bash
curl -X POST http://localhost:8081/auth/token \
  -H "Content-Type: application/json" \
  -d '{"userId":"alice"}'
# Should return {"token":"eyJ..."}
```

## Step 3: Start Signaling Service

In a new terminal:

```bash
cd backend
go run ./cmd/signaling
```

Signaling listens on `http://localhost:8080`. Verify:

```bash
curl http://localhost:8080/health/live
# Should return ok

curl http://localhost:8080/health/ready
# Should return ok (or "redis unavailable" if Redis is down)
```

## Step 4: Start Client

In a new terminal:

```bash
cd client
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Step 5: Test a Call

1. **Tab 1**: User ID `alice`, Room ID `room1` → Join Room
2. **Tab 2**: User ID `bob`, Room ID `room1` → Join Room
3. **Tab 1** (or Tab 2): Click **Start Call**
4. Grant camera/microphone permissions in both tabs
5. Verify video/audio flows between peers

## Step 6: Health Checks

| Endpoint | Service | Expected |
|----------|---------|----------|
| `GET /health/live` | Auth (8081) | `ok` |
| `GET /health/ready` | Auth (8081) | `ok` |
| `GET /health/live` | Signaling (8080) | `ok` |
| `GET /health/ready` | Signaling (8080) | `ok` (or `redis unavailable`) |

## Environment Variables

### Auth

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_PORT` | `8081` | HTTP port |
| `AUTH_SECRET` | (hardcoded) | JWT signing secret — **change in production** |

### Signaling

| Variable | Default | Description |
|----------|---------|-------------|
| `SIGNALING_PORT` | `8080` | HTTP/WebSocket port |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `AUTH_SECRET` | (hardcoded) | Must match Auth service |

## Production Considerations

1. **TLS** — Use a reverse proxy (nginx, Caddy) for TLS termination
2. **Secrets** — Store `AUTH_SECRET` in a secrets manager
3. **Redis** — Use Redis Sentinel or Cluster for HA
4. **Scaling** — Run multiple Signaling pods; use sticky sessions for WebSocket affinity
5. **CORS** — Restrict origins in production

## Troubleshooting

| Issue | Check |
|-------|-------|
| "Failed to fetch token" | Auth service running? Correct port? |
| "WebSocket connection failed" | Signaling running? Token valid? |
| No video | Camera/mic permissions? Same room? |
| "redis unavailable" | Redis running? `REDIS_ADDR` correct? |

---

*By:- Faisal Hanif | imfanee@gmail.com*
