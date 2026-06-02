package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/markmorcos/naharda/api/internal/http/respond"
	"github.com/markmorcos/naharda/api/internal/stream"
)

// Stream is the SSE live-update endpoint (project.md add-live-updates). It sends
// an initial snapshot, then an "update" event whenever ingest NOTIFYs a change,
// with a heartbeat to survive proxy idle timeouts. It is no-store and additive —
// the cacheable /v1/* endpoints are unaffected.
func (h *Handlers) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		respond.Error(w, http.StatusInternalServerError, "unsupported", "Streaming not supported.", 0)
		return
	}

	ch := make(chan []byte, 8)
	if !h.hub.Subscribe(ch) {
		respond.Error(w, http.StatusServiceUnavailable, "stream_full", "Too many live connections; retry shortly.", 5)
		return
	}
	defer h.hub.Unsubscribe(ch)

	hdr := w.Header()
	hdr.Set("Content-Type", "text/event-stream; charset=utf-8")
	hdr.Set("Cache-Control", "no-store")
	hdr.Set("Connection", "keep-alive")
	hdr.Set("X-Accel-Buffering", "no") // disable nginx response buffering for SSE
	w.WriteHeader(http.StatusOK)

	// Initial snapshot so the client has data immediately.
	if snap, err := stream.BuildSnapshot(r.Context(), h.store); err == nil {
		fmt.Fprintf(w, "event: snapshot\ndata: %s\n\n", snap)
		flusher.Flush()
	}

	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done(): // client disconnected
			return
		case <-h.hub.Done(): // server draining
			return
		case <-heartbeat.C:
			if _, err := w.Write([]byte(": heartbeat\n\n")); err != nil {
				return
			}
			flusher.Flush()
		case frame, open := <-ch:
			if !open {
				return
			}
			if _, err := w.Write(frame); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
