// Package stream provides the live-update broadcaster (SSE) and the Postgres
// LISTEN/NOTIFY bridge (project.md add-live-updates). It is an additive
// side-channel — it never touches the cacheable /v1/* REST path.
package stream

import "sync"

// Hub fans out pre-formatted SSE frames to subscribed clients. Broadcasts are
// non-blocking (a slow client drops the frame rather than stalling the hub).
type Hub struct {
	mu       sync.Mutex
	clients  map[chan []byte]struct{}
	maxConns int
	done     chan struct{}
}

// NewHub creates a Hub. maxConns <= 0 means unlimited.
func NewHub(maxConns int) *Hub {
	return &Hub{
		clients:  make(map[chan []byte]struct{}),
		maxConns: maxConns,
		done:     make(chan struct{}),
	}
}

// Subscribe registers a client channel. Returns false if at capacity.
func (h *Hub) Subscribe(ch chan []byte) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.maxConns > 0 && len(h.clients) >= h.maxConns {
		return false
	}
	h.clients[ch] = struct{}{}
	return true
}

// Unsubscribe removes a client channel.
func (h *Hub) Unsubscribe(ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, ch)
}

// Broadcast delivers a frame to all clients without blocking on slow ones.
func (h *Hub) Broadcast(frame []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- frame:
		default: // slow client — drop this frame
		}
	}
}

// Count returns the number of connected clients.
func (h *Hub) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

// Done is closed by Shutdown so stream handlers drain.
func (h *Hub) Done() <-chan struct{} { return h.done }

// Shutdown signals all handlers to close (graceful drain). Safe to call once.
func (h *Hub) Shutdown() {
	select {
	case <-h.done:
		// already closed
	default:
		close(h.done)
	}
}
