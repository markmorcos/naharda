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
	"github.com/markmorcos/naharda/api/internal/stream"
)

// NewRouter assembles the chi router with the full middleware chain (§9).
func NewRouter(cfg config.Config, st *store.Store, logger *slog.Logger, hub *stream.Hub) http.Handler {
	r := chi.NewRouter()

	// Order matters: recover → real-ip → request-id → logging → cors → rate-limit → auth.
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)
	r.Use(mw.Logging(logger, st))
	r.Use(mw.CORS(cfg.CORSOrigins))
	r.Use(mw.RateLimit(cfg.RatePerMinute, cfg.RatePerDay))
	r.Use(mw.Auth) // reads Bearer; no-op until v2

	h := handlers.New(st, cfg.SensitiveEnabled, hub)

	// Health/readiness are unversioned and uncached (§9.4).
	r.Get("/healthz", h.Healthz)
	r.Get("/readyz", h.Readyz)

	// Versioned API surface.
	r.Route("/v1", func(r chi.Router) {
		// Safe (🟢) read endpoints.
		r.Get("/cities", h.Cities)
		r.Get("/calendar", h.Calendar)
		r.Get("/prayer-times/{city}", h.PrayerTimes)
		r.Get("/weather/{city}", h.Weather)
		r.Get("/aqi/{city}", h.AirQuality)
		r.Get("/fuel", h.Fuel)

		// FX + gold (official / world-derived; 🟡 fields empty until add-sensitive-sources).
		r.Get("/fx", h.FX)
		r.Get("/fx/history", h.FXHistory)
		r.Get("/gold", h.Gold)
		r.Get("/gold/history", h.GoldHistory)

		// Dashboard support: email capture + public stats (§10).
		r.Post("/signups", h.Signups)
		r.Delete("/signups", h.DeleteSignup) // GDPR erasure (add-privacy-gdpr)
		r.Get("/stats", h.Stats)

		// Live updates (SSE) — additive, no-store (add-live-updates).
		r.Get("/stream", h.Stream)

		// Machine-readable API description (add-openapi).
		r.Get("/openapi.json", h.OpenAPI)
	})

	return r
}
