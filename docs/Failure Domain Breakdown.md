# Failure Domain Breakdown

<!--
  Failure domain analysis and resilience strategy for carrier-grade WebRTC.
  Identifies single points of failure and mitigation approaches.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Failure Domains

### Domain 1: Client / Network

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Client disconnect | Single user loses session | Reconnect with same token; session in Redis with TTL |
| Network partition | Users cannot reach backend | Retry with exponential backoff; offline PWA state |
| Browser crash | Call drops | Reconnection flow; peer notified via signaling |

### Domain 2: Edge / Load Balancer

| Failure | Impact | Mitigation |
|---------|--------|------------|
| LB node down | Partial traffic loss | Multiple LB nodes; health checks |
| TLS cert expiry | All HTTPS fails | Automated cert renewal (e.g. Let's Encrypt) |
| DDoS | Service unavailable | Rate limiting, WAF, CDN absorption |

### Domain 3: Signaling Service

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Pod crash | Active WebSockets drop | Multiple pods; LB redistributes |
| OOM | Pod killed | Resource limits; load shedding |
| Bug / panic | Request failures | Circuit breaker to Auth/Redis; graceful shutdown |

### Domain 4: Auth Service

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Auth down | New connections fail | Circuit breaker; cached validation (short TTL) |
| Token validation slow | Latency spike | Timeout; circuit open after threshold |

### Domain 5: Redis

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Redis down | No session persistence | Fallback to in-memory (degraded); Redis Sentinel/Cluster |
| Memory full | Eviction, session loss | Maxmemory policy; monitoring |
| Network partition | Split-brain risk | Redis Sentinel; quorum-based failover |

## Cascading Failure Prevention

1. **Circuit Breaker** — Stop calling Auth/Redis when failure rate exceeds threshold.
2. **Load Shedding** — Reject new connections when CPU/memory exceeds limit.
3. **Timeouts** — All outbound calls have bounded timeout.
4. **Graceful Degradation** — Signaling can operate with reduced features if Redis is unavailable (e.g. no multi-room).

## Recovery Order

1. Redis (data layer)
2. Auth (security layer)
3. Signaling (application layer)
4. Edge (traffic layer)

---

*By:- Faisal Hanif | imfanee@gmail.com*
