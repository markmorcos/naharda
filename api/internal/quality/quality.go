// Package quality holds the data-quality guard: outlier detection and alerting
// (project.md §9.5). Alerts are self-contained — field values are folded into
// the message text — and delivered to Telegram and/or a generic webhook.
package quality

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Alerter emits data-quality alerts to slog and any configured channel.
type Alerter struct {
	webhookURL string
	tgToken    string
	tgChatID   string
	client     *http.Client
	log        *slog.Logger
}

// NewAlerter constructs an Alerter. Any destination may be empty: with none
// configured it logs only; with both, it delivers to both.
func NewAlerter(webhookURL, telegramToken, telegramChatID string, log *slog.Logger) *Alerter {
	return &Alerter{
		webhookURL: webhookURL,
		tgToken:    telegramToken,
		tgChatID:   telegramChatID,
		client:     &http.Client{Timeout: 5 * time.Second},
		log:        log,
	}
}

// Alert logs a WARN and best-effort delivers a self-contained message to every
// configured channel. Delivery failures are logged, never blocking the caller.
func (a *Alerter) Alert(ctx context.Context, msg string, fields map[string]any) {
	a.log.Warn("data-quality alert", "msg", msg, "fields", fields)
	text := formatAlert(msg, fields)
	if a.tgToken != "" && a.tgChatID != "" {
		a.post(ctx, "https://api.telegram.org/bot"+a.tgToken+"/sendMessage",
			map[string]any{"chat_id": a.tgChatID, "text": text}, "telegram")
	}
	if a.webhookURL != "" {
		a.post(ctx, a.webhookURL, map[string]any{"text": text}, "webhook")
	}
}

// formatAlert folds the fields into one human-readable line (keys sorted for
// stable output) so any plain-text channel conveys the detail.
func formatAlert(msg string, fields map[string]any) string {
	if len(fields) == 0 {
		return msg
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, fields[k]))
	}
	return msg + " — " + strings.Join(parts, " · ")
}

func (a *Alerter) post(ctx context.Context, url string, payload map[string]any, channel string) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		// The error embeds the request URL, which for Telegram contains the bot
		// token — redact it before logging (§9.6: never log secrets).
		msg := err.Error()
		if a.tgToken != "" {
			msg = strings.ReplaceAll(msg, a.tgToken, "***")
		}
		a.log.Warn("alert delivery failed", "channel", channel, "err", msg)
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
