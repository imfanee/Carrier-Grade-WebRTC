/**
 * WebRTC peer connection manager for media and data channels.
 * Handles SDP exchange, ICE candidates, and media streams.
 * By:- Faisal Hanif | imfanee@gmail.com
 */

export class WebRTCClient {
  private peerConnection: RTCPeerConnection | null = null;
  private localStream: MediaStream | null = null;
  private remoteStreamCallback: ((stream: MediaStream) => void) | null = null;

  async createPeerConnection(config?: RTCConfiguration): Promise<RTCPeerConnection> {
    const defaultConfig: RTCConfiguration = {
      iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
    };
    this.peerConnection = new RTCPeerConnection(config ?? defaultConfig);

    this.peerConnection.ontrack = (event) => {
      if (event.streams[0]) {
        this.remoteStreamCallback?.(event.streams[0]);
      }
    };

    return this.peerConnection;
  }

  async getLocalStream(): Promise<MediaStream> {
    if (this.localStream) {
      return this.localStream;
    }
    this.localStream = await navigator.mediaDevices.getUserMedia({
      audio: true,
      video: { width: 640, height: 480 },
    });
    return this.localStream;
  }

  async createOffer(): Promise<RTCSessionDescriptionInit> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not created');
    }
    const stream = await this.getLocalStream();
    stream.getTracks().forEach((track) => {
      this.peerConnection!.addTrack(track, stream);
    });
    const offer = await this.peerConnection.createOffer();
    await this.peerConnection.setLocalDescription(offer);
    return offer;
  }

  async createAnswer(): Promise<RTCSessionDescriptionInit> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not created');
    }
    const stream = await this.getLocalStream();
    stream.getTracks().forEach((track) => {
      this.peerConnection!.addTrack(track, stream);
    });
    const answer = await this.peerConnection.createAnswer();
    await this.peerConnection.setLocalDescription(answer);
    return answer;
  }

  async setRemoteDescription(sdp: RTCSessionDescriptionInit): Promise<void> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not created');
    }
    await this.peerConnection.setRemoteDescription(new RTCSessionDescription(sdp));
  }

  async addIceCandidate(candidate: RTCIceCandidateInit): Promise<void> {
    if (!this.peerConnection) return;
    await this.peerConnection.addIceCandidate(new RTCIceCandidate(candidate));
  }

  onIceCandidate(callback: (candidate: RTCIceCandidate) => void): void {
    if (!this.peerConnection) return;
    this.peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        callback(event.candidate);
      }
    };
  }

  onRemoteStream(callback: (stream: MediaStream) => void): void {
    this.remoteStreamCallback = callback;
  }

  close(): void {
    this.peerConnection?.close();
    this.peerConnection = null;
    this.localStream?.getTracks().forEach((t) => t.stop());
    this.localStream = null;
  }

  get connectionState(): string {
    return this.peerConnection?.connectionState ?? 'closed';
  }
}
