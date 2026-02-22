// Package hub — WebSocket signaling hub for WebRTC SDP/ICE relay.
//
// Manages peers, rooms, and message routing between connected clients.
// By:- Faisal Hanif | imfanee@gmail.com

package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/faisalhanif/carrier-grade-webrtc/pkg/contracts"
)

// SignalMessage represents a signaling message (join, offer, answer, ice-candidate, leave).
type SignalMessage struct {
	Type     string            `json:"type"`
	RoomID   string            `json:"roomId,omitempty"`
	PeerID   string            `json:"peerId,omitempty"`
	Peers    []string          `json:"peers,omitempty"`
	SDP      string            `json:"sdp,omitempty"`
	Candidate json.RawMessage   `json:"candidate,omitempty"`
}

// SignalHub manages connected peers and room membership.
type SignalHub struct {
	peers     map[string]*Peer
	rooms     map[string]map[string]*Peer
	store     contracts.SessionStore
	mu        sync.RWMutex
}

// Peer represents a connected WebSocket client.
type Peer struct {
	ID       string
	RoomID   string
	Conn     *websocket.Conn
	Send     chan []byte
}

// NewSignalHub creates a new signaling hub.
func NewSignalHub(store contracts.SessionStore) *SignalHub {
	if store == nil {
		store = &NoopStore{}
	}
	return &SignalHub{
		peers: make(map[string]*Peer),
		rooms: make(map[string]map[string]*Peer),
		store: store,
	}
}

// Register adds a peer to the hub.
func (h *SignalHub) Register(peerID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	peer := &Peer{
		ID:   peerID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	h.peers[peerID] = peer
	go peer.writePump()
}

// Unregister removes a peer and cleans up room membership.
func (h *SignalHub) Unregister(peerID string) {
	h.mu.Lock()
	peer, ok := h.peers[peerID]
	if !ok {
		h.mu.Unlock()
		return
	}
	delete(h.peers, peerID)
	if peer.RoomID != "" {
		if room, exists := h.rooms[peer.RoomID]; exists {
			delete(room, peerID)
			if len(room) == 0 {
				delete(h.rooms, peer.RoomID)
			}
		}
	}
	h.mu.Unlock()
	close(peer.Send)
}

// HandleMessage processes incoming signaling messages.
func (h *SignalHub) HandleMessage(peerID string, msg SignalMessage) {
	h.mu.Lock()
	peer, ok := h.peers[peerID]
	if !ok {
		h.mu.Unlock()
		return
	}
	h.mu.Unlock()

	switch msg.Type {
	case "join":
		h.handleJoin(peer, msg.RoomID)
	case "offer":
		h.relayToPeer(peer, msg.PeerID, msg)
	case "answer":
		h.relayToPeer(peer, msg.PeerID, msg)
	case "ice-candidate":
		h.relayToPeer(peer, msg.PeerID, msg)
	case "leave":
		h.handleLeave(peer, msg.RoomID)
	default:
		h.sendToPeer(peer, SignalMessage{Type: "error", PeerID: msg.PeerID})
	}
}

func (h *SignalHub) handleJoin(peer *Peer, roomID string) {
	if roomID == "" {
		h.sendToPeer(peer, SignalMessage{Type: "error"})
		return
	}
	h.mu.Lock()
	if peer.RoomID != "" {
		if room, exists := h.rooms[peer.RoomID]; exists {
			delete(room, peer.ID)
			if len(room) == 0 {
				delete(h.rooms, peer.RoomID)
			}
		}
	}
	peer.RoomID = roomID
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[string]*Peer)
	}
	existingPeers := make([]string, 0, len(h.rooms[roomID]))
	for id := range h.rooms[roomID] {
		existingPeers = append(existingPeers, id)
	}
	h.rooms[roomID][peer.ID] = peer
	h.mu.Unlock()

	h.sendToPeer(peer, SignalMessage{Type: "joined", RoomID: roomID, PeerID: peer.ID, Peers: existingPeers})

	// Notify existing peers that a new peer joined
	for _, otherID := range existingPeers {
		h.mu.RLock()
		otherPeer := h.peers[otherID]
		h.mu.RUnlock()
		if otherPeer != nil {
			h.sendToPeer(otherPeer, SignalMessage{Type: "peer_joined", PeerID: peer.ID})
		}
	}
}

func (h *SignalHub) handleLeave(peer *Peer, roomID string) {
	h.mu.Lock()
	if peer.RoomID != "" {
		if room, exists := h.rooms[peer.RoomID]; exists {
			delete(room, peer.ID)
			if len(room) == 0 {
				delete(h.rooms, peer.RoomID)
			}
		}
		peer.RoomID = ""
	}
	h.mu.Unlock()
}

func (h *SignalHub) relayToPeer(from *Peer, toPeerID string, msg SignalMessage) {
	if toPeerID == "" {
		return
	}
	h.mu.RLock()
	toPeer, ok := h.peers[toPeerID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	msg.PeerID = from.ID
	h.sendToPeer(toPeer, msg)
}

func (h *SignalHub) sendToPeer(peer *Peer, msg SignalMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case peer.Send <- data:
	default:
		log.Printf("dropped message to peer %s", peer.ID)
	}
}

func (p *Peer) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-p.Send:
			if !ok {
				return
			}
			if err := p.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NoopStore is a no-op SessionStore for when Redis is unavailable.
type NoopStore struct{}

func (n *NoopStore) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error { return nil }
func (n *NoopStore) Get(_ context.Context, _ string) ([]byte, error)                  { return nil, nil }
func (n *NoopStore) Delete(_ context.Context, _ string) error                         { return nil }
