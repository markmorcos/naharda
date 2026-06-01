package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/markmorcos/naharda/api/internal/domain"
	"github.com/markmorcos/naharda/api/internal/http/respond"
)

// Signups stores an email with single opt-in consent (§10, §12). No email is
// sent in v1. A populated honeypot field is silently accepted but not stored.
func (h *Handlers) Signups(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email   string `json:"email"`
		Consent bool   `json:"consent"`
		Locale  string `json:"locale"`
		Website string `json:"website"` // honeypot — must be empty
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&body); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid_body", "Malformed request body.", 0)
		return
	}
	if body.Website != "" { // bot: pretend success, store nothing
		writeOK(w, http.StatusOK)
		return
	}
	addr, err := mail.ParseAddress(strings.TrimSpace(body.Email))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid_email", "A valid email is required.", 0)
		return
	}
	if !body.Consent {
		respond.Error(w, http.StatusBadRequest, "consent_required", "Consent is required to subscribe.", 0)
		return
	}
	if err := h.store.InsertSignup(r.Context(), strings.ToLower(addr.Address), true, body.Locale); err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "Could not save signup.", 30)
		return
	}
	writeOK(w, http.StatusCreated)
}

// Stats returns public aggregate metrics (§10). 5-minute cache.
func (h *Handlers) Stats(w http.ResponseWriter, r *http.Request) {
	st, err := h.store.GetStats(r.Context())
	if err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "Stats unavailable.", 30)
		return
	}
	respond.JSON(w, r, 300, st, domain.Meta{
		Attribution: "Public stats via Naharda. Free tier non-commercial.",
	})
}

func writeOK(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"ok":true}`))
}
