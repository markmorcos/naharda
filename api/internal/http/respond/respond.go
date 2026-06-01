// Package respond writes the standard JSON envelope, ETag, Cache-Control, and
// errors (project.md §9.1, §9.2).
package respond

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// JSON writes data in the standard envelope with the given Cache-Control max-age,
// sets an ETag, and returns 304 when If-None-Match matches.
func JSON(w http.ResponseWriter, r *http.Request, maxAge int, data any, meta domain.Meta) {
	if meta.CachedAt.IsZero() {
		meta.CachedAt = time.Now().UTC()
	}
	if meta.Tier == "" {
		meta.Tier = "anonymous"
	}
	if meta.Attribution == "" {
		meta.Attribution = domain.DefaultAttribution
	}
	if meta.Sources == nil {
		meta.Sources = []domain.Source{}
	}
	// Report data age as the gap between cache time and the newest source fetch.
	if meta.FreshnessSeconds == 0 {
		var newest time.Time
		for _, src := range meta.Sources {
			if src.FetchedAt.After(newest) {
				newest = src.FetchedAt
			}
		}
		if !newest.IsZero() {
			if age := int(meta.CachedAt.Sub(newest).Seconds()); age > 0 {
				meta.FreshnessSeconds = age
			}
		}
	}

	// ETag hashes the DATA payload only — not the volatile meta timestamps
	// (cached_at / fetched_at) — so it stays stable within the cache window and
	// conditional requests actually match (§9.1).
	dataBytes, err := json.Marshal(data)
	if err != nil {
		Error(w, http.StatusInternalServerError, "internal", "Failed to encode response.", 0)
		return
	}
	etag := etagFor(dataBytes)
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if etagMatches(r.Header.Get("If-None-Match"), etag) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	body, err := json.Marshal(domain.Envelope{Data: data, Meta: meta})
	if err != nil {
		Error(w, http.StatusInternalServerError, "internal", "Failed to encode response.", 0)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// Error writes the standard error envelope (English-only messages in v1).
func Error(w http.ResponseWriter, status int, code, message string, retryAfter int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(domain.ErrorEnvelope{
		Error: domain.ErrorBody{Code: code, Message: message, RetryAfterSeconds: retryAfter},
	})
}

func etagFor(b []byte) string {
	sum := sha256.Sum256(b)
	return `"` + hex.EncodeToString(sum[:16]) + `"`
}

// etagMatches compares If-None-Match against our ETag, tolerating weak
// validators (W/"...") that proxies like Cloudflare add when re-encoding, and
// the comma-separated list / "*" forms.
func etagMatches(ifNoneMatch, etag string) bool {
	if ifNoneMatch == "" {
		return false
	}
	want := strings.TrimPrefix(etag, "W/")
	for _, tok := range strings.Split(ifNoneMatch, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "*" || strings.TrimPrefix(tok, "W/") == want {
			return true
		}
	}
	return false
}
