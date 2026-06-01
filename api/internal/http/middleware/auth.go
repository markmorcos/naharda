package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const apiKeyCtxKey ctxKey = "naharda.apiKey"

// Auth reads an optional "Authorization: Bearer <key>" and stashes it in the
// request context. It is a no-op for access control in v1 (§9.3); per-key
// quotas activate in v2.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
			key := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
			if key != "" {
				r = r.WithContext(context.WithValue(r.Context(), apiKeyCtxKey, key))
			}
		}
		next.ServeHTTP(w, r)
	})
}

// APIKey returns the bearer key from the request context, if present.
func APIKey(r *http.Request) string {
	if v, ok := r.Context().Value(apiKeyCtxKey).(string); ok {
		return v
	}
	return ""
}
