# Carrier-Grade WebRTC Client

<!--
  NextJS PWA client for WebRTC signaling and peer connections.
  Connects to Auth and Signaling backend services.
  By:- Faisal Hanif | imfanee@gmail.com
-->

A Progressive Web App (PWA) that connects to the carrier-grade WebRTC backend for real-time video/audio calls.

## Prerequisites

- Node.js 18+
- Auth service running on port 8081
- Signaling service running on port 8080

## Setup

```bash
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Usage

1. Enter a **User ID** (e.g. `alice` or `bob`)
2. Enter a **Room ID** (e.g. `room1`) — both peers must use the same room
3. Click **Join Room**
4. When both peers are in the room, click **Start Call** on one side
5. Grant camera/microphone permissions when prompted

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_AUTH_URL` | `http://localhost:8081` | Auth service URL |
| `NEXT_PUBLIC_SIGNALING_URL` | `http://localhost:8080` | Signaling service URL |

## PWA

The app includes a `manifest.json` for installability. Add `icon-192.png` and `icon-512.png` to `public/` for full PWA support.

---

*By:- Faisal Hanif | imfanee@gmail.com*
