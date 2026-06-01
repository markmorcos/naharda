// Package domain holds core types and invariants shared across the service.
package domain

import "time"

// DefaultAttribution is the baseline attribution string (project.md §9.1).
const DefaultAttribution = "Data via Naharda (naharda.com). Free tier non-commercial."

// Source records the provenance of a datum (project.md §2.5 — provenance always).
type Source struct {
	Name      string    `json:"name"`
	URL       string    `json:"url,omitempty"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Meta is the envelope metadata present on every successful response (§9.1).
type Meta struct {
	Tier             string    `json:"tier"`
	CachedAt         time.Time `json:"cached_at"`
	FreshnessSeconds int       `json:"freshness_seconds"`
	Sources          []Source  `json:"sources"`
	Attribution      string    `json:"attribution"`
	Stale            bool      `json:"stale,omitempty"` // set when serving a last-good value (§9.5)
}

// Envelope is the standard success response shape (§9.1).
type Envelope struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

// ErrorBody is the standard error payload (§9.1).
type ErrorBody struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	RetryAfterSeconds int    `json:"retry_after_seconds,omitempty"`
}

// ErrorEnvelope wraps ErrorBody under an "error" key.
type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}
