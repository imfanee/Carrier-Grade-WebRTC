/**
 * WebSocket signaling client for WebRTC SDP/ICE relay.
 * Connects to backend signaling service with JWT authentication.
 * By:- Faisal Hanif | imfanee@gmail.com
 */

export interface SignalMessage {
  type: string;
  roomId?: string;
  peerId?: string;
  peers?: string[];
  sdp?: string;
  candidate?: RTCIceCandidateInit;
}

export type SignalMessageHandler = (msg: SignalMessage) => void;

export class SignalingClient {
  private webSocket: WebSocket | null = null;
  private messageHandler: SignalMessageHandler | null = null;
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl.replace(/^http/, 'ws');
  }

  connect(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const url = `${this.baseUrl}/ws/signal?token=${encodeURIComponent(token)}`;
      this.webSocket = new WebSocket(url);

      this.webSocket.onopen = () => resolve();
      this.webSocket.onerror = () => reject(new Error('WebSocket connection failed'));
      this.webSocket.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data) as SignalMessage;
          this.messageHandler?.(msg);
        } catch {
          // Ignore parse errors
        }
      };
    });
  }

  onMessage(handler: SignalMessageHandler): void {
    this.messageHandler = handler;
  }

  send(message: SignalMessage): void {
    if (this.webSocket?.readyState === WebSocket.OPEN) {
      this.webSocket.send(JSON.stringify(message));
    }
  }

  joinRoom(roomId: string): void {
    this.send({ type: 'join', roomId });
  }

  sendOffer(peerId: string, sdp: string): void {
    this.send({ type: 'offer', peerId, sdp });
  }

  sendAnswer(peerId: string, sdp: string): void {
    this.send({ type: 'answer', peerId, sdp });
  }

  sendIceCandidate(peerId: string, candidate: RTCIceCandidateInit): void {
    this.send({ type: 'ice-candidate', peerId, candidate });
  }

  leaveRoom(roomId: string): void {
    this.send({ type: 'leave', roomId });
  }

  disconnect(): void {
    if (this.webSocket) {
      this.webSocket.close();
      this.webSocket = null;
    }
  }

  get isConnected(): boolean {
    return this.webSocket?.readyState === WebSocket.OPEN;
  }
}
