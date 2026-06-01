// Package middleware holds the HTTP middleware chain (project.md §9).
package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/markmorcos/naharda/api/internal/store"
)

// Logging emits a structured JSON log per request and best-effort writes to
// usage_log (§9.4). IP and key are hashed; the raw key is never logged.
func Logging(logger *slog.Logger, st *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(ww, r)

			endpoint := r.Method + " " + r.URL.Path
			ipHash := hashShort(r.RemoteAddr)
			keyHash := ""
			if k := bearer(r); k != "" {
				keyHash = hashShort(k)
			}

			logger.Info("request",
				"endpoint", endpoint,
				"ip_hash", ipHash,
				"key_hash", emptyToNil(keyHash),
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"dur_ms", time.Since(start).Milliseconds(),
			)

			go st.LogUsage(endpoint, ipHash, keyHash, ww.Status(), ww.BytesWritten())
		})
	}
}

func hashShort(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:8])
}

func emptyToNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}
	return ""
}
