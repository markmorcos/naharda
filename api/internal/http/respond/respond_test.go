package respond

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/markmorcos/naharda/api/internal/domain"
)

func TestJSONEnvelopeAndConditional(t *testing.T) {
	data := map[string]string{"k": "v"}

	r := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	JSON(w, r, 300, data, domain.Meta{})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=300" {
		t.Errorf("cache-control = %q", cc)
	}
	etag := w.Header().Get("ETag")
	if etag == "" {
		t.Fatal("missing ETag")
	}
	body := w.Body.String()
	if !strings.Contains(body, `"data"`) || !strings.Contains(body, `"attribution"`) {
		t.Errorf("envelope missing fields: %s", body)
	}

	// Same data + matching If-None-Match → 304.
	r2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	r2.Header.Set("If-None-Match", etag)
	w2 := httptest.NewRecorder()
	JSON(w2, r2, 300, data, domain.Meta{})
	if w2.Code != http.StatusNotModified {
		t.Errorf("conditional GET status = %d, want 304", w2.Code)
	}
}

func TestETagMatches(t *testing.T) {
	etag := `"abc123"`
	cases := []struct {
		inm  string
		want bool
	}{
		{`"abc123"`, true},
		{`W/"abc123"`, true},      // weak validator (Cloudflare re-encode)
		{`"x", W/"abc123"`, true}, // comma list
		{`*`, true},               // wildcard
		{`"nope"`, false},
		{``, false},
	}
	for _, c := range cases {
		if got := etagMatches(c.inm, etag); got != c.want {
			t.Errorf("etagMatches(%q,%q)=%v want %v", c.inm, etag, got, c.want)
		}
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusTooManyRequests, "rate_limited", "slow down", 30)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "rate_limited") {
		t.Errorf("body = %s", w.Body.String())
	}
}
