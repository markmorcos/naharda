// Package fx ingests official EGP exchange rates and applies the data-quality
// guard (project.md §9.5).
package fx

import (
	"context"
	"log/slog"
	"time"

	"github.com/markmorcos/naharda/api/internal/quality"
	"github.com/markmorcos/naharda/api/internal/sources"
	"github.com/markmorcos/naharda/api/internal/store"
)

const outlierThresholdPct = 5.0

// Run fetches official FX, holds outliers as pending_review, and stores the rest.
func Run(ctx context.Context, st *store.Store, alerter *quality.Alerter, log *slog.Logger) {
	rates, src, err := sources.FetchOfficialFX(ctx)
	if err != nil {
		log.Warn("fx fetch failed", "err", err)
		return
	}
	held := 0
	for quote, value := range rates {
		pending := false
		// Outlier guard: only meaningful with a baseline (≥3 recent samples).
		if avg, n, err := st.TrailingAvgFX(ctx, "official", quote, time.Hour); err == nil &&
			n >= 3 && quality.IsOutlier(value, avg, outlierThresholdPct) {
			pending = true
			held++
			alerter.Alert(ctx, "fx official outlier held", map[string]any{
				"quote": quote, "value": value, "trailing_avg": avg,
			})
		}
		if err := st.InsertFXRate(ctx, "official", quote, value, src.Name, pending); err != nil {
			log.Warn("fx insert failed", "quote", quote, "err", err)
		}
	}
	log.Info("fx ingest complete", "quotes", len(rates), "held", held)
}
