// Package quality holds the data-quality skeleton: outlier detection and
// alerting (project.md §9.5). The full guard logic is wired up in the
// add-fx-official-and-gold-world change.
package quality

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// Alerter emits data-quality alerts to slog and an optional webhook (ntfy/Slack).
type Alerter struct {
	webhookURL string
	client     *http.Client
	log        *slog.Logger
}

// NewAlerter constructs an Alerter. An empty webhookURL logs only.
func NewAlerter(webhookURL string, log *slog.Logger) *Alerter {
	return &Alerter{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
		log:        log,
	}
}

// Alert logs a WARN and best-effort posts to the webhook.
func (a *Alerter) Alert(ctx context.Context, msg string, fields map[string]any) {
	a.log.Warn("data-quality alert", "msg", msg, "fields", fields)
	if a.webhookURL == "" {
		return
	}
	body, err := json.Marshal(map[string]any{"text": msg, "fields": fields})
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.webhookURL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		a.log.Warn("alert webhook failed", "err", err)
		return
	}
	_ = resp.Body.Close()
}

// IsOutlier reports whether value deviates from trailingAvg by more than
// thresholdPct percent. A zero average is never an outlier (no baseline yet).
func IsOutlier(value, trailingAvg, thresholdPct float64) bool {
	if trailingAvg == 0 {
		return false
	}
	dev := (value - trailingAvg) / trailingAvg
	if dev < 0 {
		dev = -dev
	}
	return dev*100 > thresholdPct
}
