// Package httpapi builds the HTTP router and wires the middleware chain.
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/markmorcos/naharda/api/internal/config"
	"github.com/markmorcos/naharda/api/internal/http/handlers"
	mw "github.com/markmorcos/naharda/api/internal/http/middleware"
	"github.com/markmorcos/naharda/api/internal/store"
)

// NewRouter assembles the chi router with the full middleware chain (§9).
func NewRouter(cfg config.Config, st *store.Store, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	// Order matters: recover → real-ip → request-id → logging → cors → rate-limit → auth.
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)
	r.Use(mw.Logging(logger, st))
	r.Use(mw.CORS(cfg.CORSOrigins))
	r.Use(mw.RateLimit(cfg.RatePerMinute, cfg.RatePerDay))
	r.Use(mw.Auth) // reads Bearer; no-op until v2

	h := handlers.New(st)

	// Health/readiness are unversioned and uncached (§9.4).
	r.Get("/healthz", h.Healthz)
	r.Get("/readyz", h.Readyz)

	// Versioned API surface. Endpoints are added by later changes.
	r.Route("/v1", func(r chi.Router) {})

	return r
}
