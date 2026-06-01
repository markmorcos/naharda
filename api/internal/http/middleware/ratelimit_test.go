package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimit(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	h := RateLimit(2, 100)(next)

	call := func(ip string) int {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = ip
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}

	if call("1.2.3.4:1") != http.StatusOK || call("1.2.3.4:2") != http.StatusOK {
		t.Fatal("first two requests should pass")
	}
	if got := call("1.2.3.4:3"); got != http.StatusTooManyRequests {
		t.Errorf("third request = %d, want 429", got)
	}
	// A different IP has its own bucket.
	if call("9.9.9.9:1") != http.StatusOK {
		t.Error("distinct IP should not be limited")
	}
}
