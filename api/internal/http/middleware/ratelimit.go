package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/markmorcos/naharda/api/internal/http/respond"
)

// window tracks per-IP fixed-window counters for minute and day budgets.
type window struct {
	minStart time.Time
	minCount int
	dayStart time.Time
	dayCount int
}

type limiter struct {
	mu     sync.Mutex
	perMin int
	perDay int
	ips    map[string]*window
}

// RateLimit enforces per-IP minute and day budgets (§9.3). In-memory; a single
// replica in v1. The Authorization header is read elsewhere (Auth) but does not
// raise limits until v2.
func RateLimit(perMin, perDay int) func(http.Handler) http.Handler {
	l := &limiter{perMin: perMin, perDay: perDay, ips: make(map[string]*window)}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ok, retry := l.allow(r.RemoteAddr); !ok {
				w.Header().Set("Retry-After", strconv.Itoa(retry))
				respond.Error(w, http.StatusTooManyRequests, "rate_limited", "Rate limit exceeded.", retry)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (l *limiter) allow(ip string) (bool, int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	wd := l.ips[ip]
	if wd == nil {
		wd = &window{minStart: now, dayStart: now}
		l.ips[ip] = wd
	}
	if now.Sub(wd.minStart) >= time.Minute {
		wd.minStart, wd.minCount = now, 0
	}
	if now.Sub(wd.dayStart) >= 24*time.Hour {
		wd.dayStart, wd.dayCount = now, 0
	}
	if wd.minCount >= l.perMin {
		return false, secsUntil(wd.minStart.Add(time.Minute))
	}
	if wd.dayCount >= l.perDay {
		return false, secsUntil(wd.dayStart.Add(24 * time.Hour))
	}
	wd.minCount++
	wd.dayCount++
	return true, 0
}

func secsUntil(t time.Time) int {
	s := int(time.Until(t).Seconds()) + 1
	if s < 1 {
		s = 1
	}
	return s
}
