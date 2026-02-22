# Architecture Overview

<!--
  High-level architecture for the Carrier-Grade WebRTC reference implementation.
  Describes components, data flow, and design principles.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Design Principles

- **Security First** — Authentication, authorization, and encryption at every boundary
- **Failure Isolation** — Bounded failure domains with circuit breakers and load shedding
- **Observability** — Structured logs, metrics, and traces for operational visibility
- **Stateless Signaling** — Session state in Redis; services horizontally scalable

## Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CLIENT LAYER                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  NextJS PWA Client                                                  │    │
│  │  • WebRTC PeerConnection  • WebSocket Signaling  • Media Capture    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ HTTPS / WSS
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           EDGE / GATEWAY LAYER                              │
│  (Reverse Proxy / Load Balancer — TLS termination, rate limiting)           │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
┌──────────────────────────┐ ┌──────────────────┐ ┌──────────────────────────┐
│   SIGNALING SERVICE      │ │   AUTH SERVICE    │ │   HEALTH AGGREGATOR     │
│   • WebSocket handlers   │ │   • Token verify  │ │   • /health/live        │
│   • SDP relay            │ │   • Session check │ │   • /health/ready       │
│   • ICE candidate relay  │ │   • JWT validation│ │   • Dependency probes   │
└──────────────────────────┘ └──────────────────┘ └──────────────────────────┘
                    │                 │
                    └────────┬────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DATA LAYER                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Redis                                                              │    │
│  │  • Session store  • Offer/Answer cache  • TTL-based expiry          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Client** authenticates via Auth service, receives JWT.
2. **Client** opens WebSocket to Signaling service with JWT in header/query.
3. **Signaling** validates token with Auth, stores session in Redis.
4. **Client** sends SDP offer; Signaling relays to peer via WebSocket.
5. **Peer** sends SDP answer; Signaling relays back.
6. **ICE candidates** exchanged through Signaling; peers establish direct P2P connection.
7. **Media** flows peer-to-peer (no media through backend).

## Security Boundaries

- All client-facing endpoints require valid JWT.
- WebSocket connections are authenticated before room/peer operations.
- Redis stores only ephemeral session data; no PII persistence.
- TLS everywhere for transport security.

## Scalability

- Signaling and Auth are stateless; scale horizontally behind load balancer.
- Redis supports clustering for high availability.
- WebSocket affinity (sticky sessions) recommended for signaling.

---

*By:- Faisal Hanif | imfanee@gmail.com*
 
