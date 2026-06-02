package stream

import (
	"testing"
	"time"
)

func TestHubCapBroadcastUnsubscribe(t *testing.T) {
	h := NewHub(2)
	a, b, c := make(chan []byte, 1), make(chan []byte, 1), make(chan []byte, 1)

	if !h.Subscribe(a) || !h.Subscribe(b) {
		t.Fatal("first two subscriptions should succeed")
	}
	if h.Subscribe(c) {
		t.Error("third subscription should be rejected (cap=2)")
	}
	if h.Count() != 2 {
		t.Errorf("count = %d, want 2", h.Count())
	}

	h.Broadcast([]byte("x"))
	if string(<-a) != "x" || string(<-b) != "x" {
		t.Error("clients did not receive broadcast")
	}

	h.Unsubscribe(a)
	if h.Count() != 1 {
		t.Errorf("count after unsubscribe = %d, want 1", h.Count())
	}
}

func TestHubBroadcastDoesNotBlockOnSlowClient(t *testing.T) {
	h := NewHub(0)            // unlimited
	slow := make(chan []byte) // unbuffered, no reader
	h.Subscribe(slow)

	done := make(chan struct{})
	go func() { h.Broadcast([]byte("y")); close(done) }()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Broadcast blocked on a slow client")
	}
}

func TestHubShutdownIdempotent(t *testing.T) {
	h := NewHub(0)
	h.Shutdown()
	h.Shutdown() // must not panic
	select {
	case <-h.Done():
	default:
		t.Error("Done() should be closed after Shutdown")
	}
}
