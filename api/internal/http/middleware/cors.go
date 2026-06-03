package middleware

import "net/http"

// CORS allows the configured origins ("*" allows all). v1 is read-only except
// the dashboard's email capture (POST /v1/signups) and GDPR erasure
// (DELETE /v1/signups).
func CORS(origins []string) func(http.Handler) http.Handler {
	allowAll := false
	set := make(map[string]bool, len(origins))
	for _, o := range origins {
		if o == "*" {
			allowAll = true
		}
		set[o] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" && (allowAll || set[origin]) {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, If-None-Match")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
