# System Topology

<!--
  Network and deployment topology for the Carrier-Grade WebRTC reference.
  Describes physical/logical layout and connectivity.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Logical Topology

```
                         ┌─────────────────────────────────────────┐
                         │              INTERNET / CDN              │
                         └─────────────────────────────────────────┘
                                          │
                                          │ HTTPS / WSS
                                          ▼
                         ┌─────────────────────────────────────────┐
                         │         EDGE (Reverse Proxy / LB)        │
                         │  • TLS termination                       │
                         │  • Rate limiting                         │
                         │  • WebSocket upgrade                     │
                         └─────────────────────────────────────────┘
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ▼                     ▼                     ▼
         ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
         │  Signaling Pod 1  │  │  Signaling Pod 2  │  │  Auth Pod        │
         │  :8080            │  │  :8080            │  │  :8081           │
         └──────────────────┘  └──────────────────┘  └──────────────────┘
                    │                     │                     │
                    └─────────────────────┼─────────────────────┘
                                          │
                                          │ TCP 6379
                                          ▼
                         ┌─────────────────────────────────────────┐
                         │         Redis (Primary + Replica)        │
                         │  • Session keys  • Offer/Answer cache    │
                         └─────────────────────────────────────────┘
```

## Deployment Zones

| Zone | Components | Purpose |
|------|------------|---------|
| **DMZ / Edge** | Reverse proxy, load balancer | TLS, DDoS mitigation, routing |
| **Application** | Signaling, Auth | Business logic, stateless |
| **Data** | Redis | Ephemeral state, cache |

## Network Segments

- **Public** — Client ↔ Edge (HTTPS/WSS)
- **Private** — Edge ↔ App services (internal network)
- **Data** — App ↔ Redis (private, no public exposure)

## Port Mapping (Reference)

| Service | Port | Protocol |
|---------|------|----------|
| Signaling | 8080 | HTTP, WebSocket |
| Auth | 8081 | HTTP |
| Redis | 6379 | TCP |
| Client (dev) | 3000 | HTTP |

## Peer-to-Peer Media Path

WebRTC media (audio/video) flows directly between browsers after signaling. The backend never touches media streams.

```
  Client A ◄────────────────────────────► Client B
              Direct P2P (UDP/TCP)
              (STUN/TURN if NAT traversal needed)
```

---

*By:- Faisal Hanif | imfanee@gmail.com*
