package handlers

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"net/http"
)

// openapiDoc is the hand-maintained OpenAPI 3.1 description of the public API.
// A drift test (router package) asserts it stays in sync with the chi routes.
//
//go:embed openapi.json
var openapiDoc []byte

// openapiETag is a strong validator over the embedded doc (it never changes at
// runtime, so it can be computed once).
var openapiETag = func() string {
	sum := sha256.Sum256(openapiDoc)
	return `"` + hex.EncodeToString(sum[:16]) + `"`
}()

// OpenAPIJSON returns the embedded OpenAPI document. Exposed for the route/spec
// drift test (add-openapi).
func OpenAPIJSON() []byte { return openapiDoc }

// OpenAPI serves the embedded OpenAPI document. Long-cached with an ETag; the
// content is immutable for a given build.
func (h *Handlers) OpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("ETag", openapiETag)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if match := r.Header.Get("If-None-Match"); match != "" && match == openapiETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openapiDoc)
}
