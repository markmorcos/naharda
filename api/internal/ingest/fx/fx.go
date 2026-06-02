// Package fx ingests official EGP exchange rates and applies the data-quality
// guard (project.md §9.5).
package fx

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
	"github.com/markmorcos/naharda/api/internal/quality"
	"github.com/markmorcos/naharda/api/internal/sources"
	"github.com/markmorcos/naharda/api/internal/store"
)

const (
	outlierThresholdPct = 5.0
	// disagreementPct flags cross-source spread between the canonical (CBE) and
	// the reference (exchangerate-api) values (§9.5).
	disagreementPct = 2.0
)

// Run sources official FX from the Central Bank of Egypt (canonical) with
// exchangerate-api as a cross-check and fallback, holds outliers as
// pending_review, and stores the rest. Fail-soft: if neither source returns
// data it logs/alerts and serves nothing new (handlers fall back to last-good).
func Run(ctx context.Context, st *store.Store, alerter *quality.Alerter, log *slog.Logger) {
	cbe, cbeSrc, cbeErr := sources.FetchCBE(ctx)
	ref, refSrc, refErr := sources.FetchOfficialFX(ctx)

	var (
		rates      map[string]float64
		src        domain.Source
		crossCheck map[string]float64 // the other source, for disagreement flagging
	)
	switch {
	case cbeErr == nil && len(cbe) > 0:
		rates, src, crossCheck = cbe, cbeSrc, ref // serve CBE; cross-check vs reference
	case refErr == nil && len(ref) > 0:
		rates, src = ref, refSrc // CBE unavailable → serve reference, alert
		alerter.Alert(ctx, "CBE unavailable; serving reference FX", map[string]any{"cbe_err": errString(cbeErr)})
		log.Warn("CBE fx unavailable; using reference cross-check", "err", cbeErr)
	default:
		log.Warn("fx fetch failed (no source)", "cbe_err", cbeErr, "ref_err", refErr)
		alerter.Alert(ctx, "fx ingest: no source available", map[string]any{
			"cbe_err": errString(cbeErr), "ref_err": errString(refErr),
		})
		return
	}

	held := 0
	for quote, value := range rates {
		// Cross-check: flag (don't suppress) disagreement beyond the threshold;
		// the canonical value is still served.
		if crossCheck != nil {
			if other, ok := crossCheck[quote]; ok && other > 0 && value > 0 {
				if diff := math.Abs(value-other) / value * 100; diff > disagreementPct {
					alerter.Alert(ctx, "fx cross-check disagreement", map[string]any{
						"quote": quote, "canonical": value, "reference": other, "diff_pct": diff,
					})
					log.Warn("fx cross-check disagreement", "quote", quote,
						"canonical", value, "reference", other, "diff_pct", diff)
				}
			}
		}

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
	_ = st.Notify(ctx, "naharda_updates", `{"family":"fx"}`)
	log.Info("fx ingest complete", "source", src.Name, "quotes", len(rates), "held", held)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
