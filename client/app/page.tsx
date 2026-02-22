/**
 * Main page — WebRTC call interface with signaling integration.
 * By:- Faisal Hanif | imfanee@gmail.com
 */

'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { SignalingClient } from '@/lib/signaling_client';
import { WebRTCClient } from '@/lib/webrtc_client';
import { fetchToken, getSignalingUrl } from '@/lib/auth_client';

export default function Home() {
  const [userId, setUserId] = useState('');
  const [roomId, setRoomId] = useState('');
  const [status, setStatus] = useState<'idle' | 'connecting' | 'joined' | 'in-call' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const signalingRef = useRef<SignalingClient | null>(null);
  const webrtcRef = useRef<WebRTCClient | null>(null);
  const remotePeerIdRef = useRef<string | null>(null);

  const connectAndJoin = useCallback(async () => {
    if (!userId.trim() || !roomId.trim()) {
      setErrorMessage('Enter userId and roomId');
      return;
    }
    setStatus('connecting');
    setErrorMessage('');
    try {
      const token = await fetchToken(userId.trim());
      const signaling = new SignalingClient(getSignalingUrl());
      signalingRef.current = signaling;

      await signaling.connect(token);
      signaling.joinRoom(roomId.trim());
      setStatus('joined');

      signaling.onMessage((msg) => {
        if (msg.type === 'joined') {
          // First peer in peers list is the remote peer (for 2-peer demo)
          const peers = msg.peers ?? [];
          if (peers.length > 0) {
            remotePeerIdRef.current = peers[0];
          }
        }
        if (msg.type === 'peer_joined' && msg.peerId) {
          remotePeerIdRef.current = msg.peerId;
        }
        if (msg.type === 'offer' && msg.peerId && msg.sdp) {
          handleOffer(msg.peerId, msg.sdp);
        }
        if (msg.type === 'answer' && msg.peerId && msg.sdp) {
          handleAnswer(msg.peerId, msg.sdp);
        }
        if (msg.type === 'ice-candidate' && msg.peerId && msg.candidate) {
          handleIceCandidate(msg.peerId, msg.candidate);
        }
      });
    } catch (err) {
      setStatus('error');
      setErrorMessage(err instanceof Error ? err.message : 'Connection failed');
    }
  }, [userId, roomId]);

  const attachLocalStream = async (webrtc: WebRTCClient) => {
    const stream = await webrtc.getLocalStream();
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }
  };

  const handleOffer = async (peerId: string, sdp: string) => {
    const webrtc = new WebRTCClient();
    webrtcRef.current = webrtc;
    await webrtc.createPeerConnection();
    await attachLocalStream(webrtc);
    webrtc.onRemoteStream((stream) => {
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = stream;
      }
      setStatus('in-call');
    });
    webrtc.onIceCandidate((candidate) => {
      signalingRef.current?.sendIceCandidate(peerId, candidate.toJSON());
    });
    await webrtc.setRemoteDescription({ type: 'offer', sdp });
    const answer = await webrtc.createAnswer();
    signalingRef.current?.sendAnswer(peerId, answer.sdp ?? '');
  };

  const handleAnswer = async (peerId: string, sdp: string) => {
    const webrtc = webrtcRef.current;
    if (webrtc) {
      await webrtc.setRemoteDescription({ type: 'answer', sdp });
    }
  };

  const handleIceCandidate = async (peerId: string, candidate: RTCIceCandidateInit) => {
    const webrtc = webrtcRef.current;
    if (webrtc) {
      await webrtc.addIceCandidate(candidate);
    }
  };

  const startCall = useCallback(async () => {
    const peerId = remotePeerIdRef.current;
    if (!peerId || !signalingRef.current?.isConnected) return;
    const webrtc = new WebRTCClient();
    webrtcRef.current = webrtc;
    await webrtc.createPeerConnection();
    await attachLocalStream(webrtc);
    webrtc.onRemoteStream((stream) => {
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = stream;
      }
      setStatus('in-call');
    });
    webrtc.onIceCandidate((candidate) => {
      signalingRef.current?.sendIceCandidate(peerId, candidate.toJSON());
    });
    const offer = await webrtc.createOffer();
    signalingRef.current?.sendOffer(peerId, offer.sdp ?? '');
  }, []);

  const disconnect = useCallback(() => {
    webrtcRef.current?.close();
    webrtcRef.current = null;
    signalingRef.current?.leaveRoom(roomId);
    signalingRef.current?.disconnect();
    signalingRef.current = null;
    remotePeerIdRef.current = null;
    if (localVideoRef.current) localVideoRef.current.srcObject = null;
    if (remoteVideoRef.current) remoteVideoRef.current.srcObject = null;
    setStatus('idle');
  }, [roomId]);

  useEffect(() => {
    return () => {
      signalingRef.current?.disconnect();
      webrtcRef.current?.close();
    };
  }, []);

  return (
    <main className="main">
      <header className="header">
        <h1>Carrier-Grade WebRTC</h1>
        <p className="subtitle">PWA Client</p>
      </header>

      <section className="card">
        <h2>Connect</h2>
        <div className="form">
          <input
            type="text"
            placeholder="User ID"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            disabled={status !== 'idle'}
          />
          <input
            type="text"
            placeholder="Room ID"
            value={roomId}
            onChange={(e) => setRoomId(e.target.value)}
            disabled={status !== 'idle'}
          />
          {status === 'idle' && (
            <button className="btn-primary" onClick={connectAndJoin}>
              Join Room
            </button>
          )}
          {status === 'joined' && (
            <button className="btn-primary" onClick={startCall}>
              Start Call
            </button>
          )}
          {(status === 'joined' || status === 'in-call') && (
            <button className="btn-secondary" onClick={disconnect}>
              Leave
            </button>
          )}
        </div>
        {errorMessage && <p className="error">{errorMessage}</p>}
        <p className="status">Status: {status}</p>
      </section>

      <section className="video-section">
        <div className="video-container">
          <video ref={localVideoRef} autoPlay muted playsInline className="video" />
          <span className="label">Local</span>
        </div>
        <div className="video-container">
          <video ref={remoteVideoRef} autoPlay playsInline className="video" />
          <span className="label">Remote</span>
        </div>
      </section>

      <footer className="footer">
        <p>By:- Faisal Hanif | imfanee@gmail.com</p>
      </footer>
    </main>
  );
}
