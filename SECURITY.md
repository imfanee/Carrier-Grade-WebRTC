# Security

<!--
  Security practices and threat model for the Carrier-Grade WebRTC reference.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Overview

This reference implementation follows security-first design. Below are the practices and assumptions for deployment.

---

## Auth Secret Handling

- **JWT signing secret** (`AUTH_SECRET`) must be:
  - At least 32 bytes (256 bits) of cryptographically random data
  - Stored in a secrets manager (e.g. HashiCorp Vault, AWS Secrets Manager) in production
  - Never committed to version control or logged
- Auth and Signaling services must share the same secret; rotate via coordinated deployment.
- In development, the default secret is acceptable only for local use.

---

## WebSocket Authentication

- All WebSocket connections to `/ws/signal` require a valid JWT in the query string (`?token=...`).
- Tokens are validated before the WebSocket upgrade; invalid tokens receive `401 Unauthorized`.
- No signaling operations occur without a valid token.
- Tokens should have short expiry (e.g. 24h for demo; consider 1h or less in production).

---

## Rate Limiting at Edge

- Rate limiting is **not** implemented in the application layer.
- It must be applied at the edge (reverse proxy, load balancer, or API gateway):
  - Per-IP connection limits for WebSocket upgrades
  - Per-IP request limits for `/auth/token` and `/auth/validate`
  - DDoS mitigation (e.g. Cloudflare, AWS Shield)
- Recommended: 100–500 WebSocket connections per IP per minute, configurable by environment.

---

## Data Handling Statement

- **Ephemeral session store only.** Redis holds:
  - Session keys and room membership
  - Offer/answer SDP blobs (short TTL)
- **No PII persistence.** User IDs and room IDs are not stored long-term; they exist only for the duration of a session.
- **No media through backend.** Audio and video flow peer-to-peer; the backend never sees or stores media.
- Session data in Redis should use TTL (e.g. 24h max) and be evicted when sessions end.

---

## Threat Model (Lite)

| Threat | Mitigation |
|-------|------------|
| Token theft | Short-lived JWTs; HTTPS only; no token in URLs in logs |
| Unauthorized signaling | JWT required for all WebSocket connections |
| Session hijacking | Token bound to session; rotate on sensitive actions |
| DoS / connection exhaustion | Rate limiting at edge; load shedding (see examples) |
| Secret leakage | Secrets in vault; no defaults in production |
| Data exfiltration | No PII in store; TLS everywhere |

---

## Reporting Security Issues

If you discover a security vulnerability, please email **imfanee@gmail.com** rather than opening a public issue.

---

*By:- Faisal Hanif | imfanee@gmail.com*
