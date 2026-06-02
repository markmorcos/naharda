package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/markmorcos/naharda/api/internal/config"
	"github.com/markmorcos/naharda/api/internal/http/handlers"
	"github.com/markmorcos/naharda/api/internal/store"
)

// TestOpenAPINoDrift asserts the hand-maintained OpenAPI document and the actual
// chi routes describe exactly the same path set — in both directions — so the
// spec can't silently drift from the API (add-openapi).
func TestOpenAPINoDrift(t *testing.T) {
	st, err := store.New(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(config.Config{}, st, slog.Default(), nil)
	routes, ok := router.(chi.Routes)
	if !ok {
		t.Fatal("router is not a chi.Routes")
	}

	// Routes from chi (skip the SSE/health-internal ones? No — the doc covers all).
	routePaths := map[string]bool{}
	err = chi.Walk(routes, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routePaths[normalize(route)] = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Paths from the OpenAPI doc.
	var doc struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(handlers.OpenAPIJSON(), &doc); err != nil {
		t.Fatalf("openapi.json is not valid JSON: %v", err)
	}
	docPaths := map[string]bool{}
	for p := range doc.Paths {
		docPaths[normalize(p)] = true
	}

	for p := range routePaths {
		if !docPaths[p] {
			t.Errorf("route %q is not documented in openapi.json", p)
		}
	}
	for p := range docPaths {
		if !routePaths[p] {
			t.Errorf("openapi.json documents %q but it is not a route", p)
		}
	}
	if t.Failed() {
		t.Logf("routes: %s", sortedKeys(routePaths))
		t.Logf("doc:    %s", sortedKeys(docPaths))
	}
}

// normalize strips any chi trailing slash so the two sources compare cleanly.
func normalize(p string) string {
	if len(p) > 1 {
		p = strings.TrimSuffix(p, "/")
	}
	return p
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
