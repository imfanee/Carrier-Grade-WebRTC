# Architecture Diagram

<!--
  Visual architecture representation for the Carrier-Grade WebRTC reference.
  By:- Faisal Hanif | imfanee@gmail.com
-->

## Component Architecture (Mermaid)

```mermaid
flowchart TB
    subgraph Client["Client Layer"]
        PWA[NextJS PWA Client]
    end

    subgraph Edge["Edge Layer"]
        LB[Load Balancer / Reverse Proxy]
    end

    subgraph App["Application Layer"]
        SIG[Signaling Service<br/>:8080]
        AUTH[Auth Service<br/>:8081]
    end

    subgraph Data["Data Layer"]
        REDIS[(Redis<br/>:6379)]
    end

    PWA -->|HTTPS / WSS| LB
    LB --> SIG
    LB --> AUTH
    SIG -->|Validate Token| AUTH
    SIG -->|Session / State| REDIS
```

## Sequence: Join and Call

```mermaid
sequenceDiagram
    participant A as Client A
    participant Auth as Auth Service
    participant Sig as Signaling Service
    participant B as Client B

    A->>Auth: POST /auth/token {userId}
    Auth-->>A: {token}
    A->>Sig: WS /ws/signal?token=...
    Sig->>Auth: validate token
    Auth-->>Sig: claims
    A->>Sig: {type: join, roomId}
    Sig-->>A: {type: joined, peers: []}

    B->>Auth: POST /auth/token {userId}
    Auth-->>B: {token}
    B->>Sig: WS /ws/signal?token=...
    B->>Sig: {type: join, roomId}
    Sig-->>B: {type: joined, peers: [A]}
    Sig-->>A: {type: peer_joined, peerId: B}

    A->>Sig: {type: offer, peerId: B, sdp}
    Sig-->>B: {type: offer, peerId: A, sdp}
    B->>Sig: {type: answer, peerId: A, sdp}
    Sig-->>A: {type: answer, peerId: B, sdp}
    A->>Sig: {type: ice-candidate, ...}
    Sig-->>B: {type: ice-candidate, ...}
    B->>Sig: {type: ice-candidate, ...}
    Sig-->>A: {type: ice-candidate, ...}

    Note over A,B: P2P media established
```

## Data Flow

```mermaid
flowchart LR
    subgraph Signaling["Signaling Path"]
        C1[Client 1] -->|SDP/ICE| SIG[Signaling]
        SIG -->|Relay| C2[Client 2]
        C2 -->|SDP/ICE| SIG
        SIG -->|Relay| C1
    end

    subgraph Media["Media Path - Direct P2P"]
        C1 -.->|RTP/RTCP| C2
    end
```

---

*By:- Faisal Hanif | imfanee@gmail.com*
