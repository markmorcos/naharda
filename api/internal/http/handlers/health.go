// Package handlers holds HTTP handlers, one concern per file.
package handlers

import (
	"net/http"

	"github.com/markmorcos/naharda/api/internal/store"
	"github.com/markmorcos/naharda/api/internal/stream"
)

// Handlers carries shared dependencies for HTTP handlers.
type Handlers struct {
	store            *store.Store
	sensitiveEnabled bool        // 🟡 parallel FX / retail gold gate (§8)
	hub              *stream.Hub // live-update broadcaster (add-live-updates)
}

// New constructs Handlers.
func New(st *store.Store, sensitiveEnabled bool, hub *stream.Hub) *Handlers {
	return &Handlers{store: st, sensitiveEnabled: sensitiveEnabled, hub: hub}
}

// Healthz reports process liveness (§9.4).
func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
	writeStatus(w, http.StatusOK, `{"status":"ok"}`)
}

// Readyz reports readiness including database reachability (§9.4).
func (h *Handlers) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Ping(r.Context()); err != nil {
		writeStatus(w, http.StatusServiceUnavailable, `{"status":"not_ready"}`)
		return
	}
	writeStatus(w, http.StatusOK, `{"status":"ready"}`)
}

func writeStatus(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(body))
}
