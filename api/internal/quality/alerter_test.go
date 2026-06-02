package quality

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestFormatAlert(t *testing.T) {
	got := formatAlert("fx official outlier held", map[string]any{
		"quote": "USD", "value": 999.0, "trailing_avg": 52.22,
	})
	// keys sorted: quote, trailing_avg, value
	want := "fx official outlier held — quote=USD · trailing_avg=52.22 · value=999"
	if got != want {
		t.Errorf("formatAlert = %q, want %q", got, want)
	}
	if formatAlert("just msg", nil) != "just msg" {
		t.Error("empty fields should return the message unchanged")
	}
}

func TestAlerterNoConfigIsNoop(t *testing.T) {
	a := NewAlerter("", "", "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	// No channels configured → logs only, no panic, no network.
	a.Alert(context.Background(), "x", map[string]any{"k": 1})
}
