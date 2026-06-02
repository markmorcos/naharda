// Package sensitive ingests 🟡 parallel FX and Egypt-retail gold. It runs only
// when SENSITIVE_SOURCES_ENABLED is true AND sources are registered (§8, §16 #1).
// Each source failure degrades only its own field (fail-soft §2.6).
package sensitive

import (
	"context"
	"log/slog"
	"time"

	"github.com/markmorcos/naharda/api/internal/quality"
	"github.com/markmorcos/naharda/api/internal/sources"
	"github.com/markmorcos/naharda/api/internal/store"
)

// defaultParallelOutlierPct guards a source that has no row in the `sources`
// registry. Parallel is noisier than official, so the fallback is wider (§9.5);
// per-source values come from the registry's outlier_threshold.
const defaultParallelOutlierPct = 8.0

// ParallelFXRun fetches each approved parallel source and stores per-source
// USD quotes (the FX handler aggregates them into {min,avg,max,n,sources}).
func ParallelFXRun(ctx context.Context, st *store.Store, alerter *quality.Alerter, log *slog.Logger) {
	stored := 0
	for _, src := range sources.RegisteredParallelSources {
		val, err := src.FetchUSD(ctx)
		if err != nil {
			log.Warn("parallel fx source failed", "source", src.Name(), "err", err)
			continue // fail-soft: skip this source only
		}
		threshold := defaultParallelOutlierPct
		if th, ok, err := st.SourceThreshold(ctx, src.Name()); err == nil && ok {
			threshold = th // tunable per source from the registry
		}
		pending := false
		if avg, n, err := st.TrailingAvgFX(ctx, "parallel", "USD", time.Hour); err == nil &&
			n >= 3 && quality.IsOutlier(val, avg, threshold) {
			pending = true
			alerter.Alert(ctx, "parallel fx outlier held", map[string]any{
				"source": src.Name(), "value": val, "trailing_avg": avg,
			})
		}
		if err := st.InsertFXRate(ctx, "parallel", "USD", val, src.Name(), pending); err != nil {
			log.Warn("parallel fx insert failed", "source", src.Name(), "err", err)
			continue
		}
		stored++
	}
	log.Info("parallel fx ingest complete", "sources", len(sources.RegisteredParallelSources), "stored", stored)
}

// RetailGoldRun fetches each approved retail-gold source and stores per-karat
// egypt_retail prices (kept separate from world_derived, never merged §4).
func RetailGoldRun(ctx context.Context, st *store.Store, _ *quality.Alerter, log *slog.Logger) {
	for _, src := range sources.RegisteredRetailGoldSources {
		perKarat, err := src.FetchPerGram(ctx)
		if err != nil {
			log.Warn("retail gold source failed", "source", src.Name(), "err", err)
			continue
		}
		for karat, value := range perKarat {
			if err := st.InsertGoldPrice(ctx, "egypt_retail", karat, value, src.Name(), false); err != nil {
				log.Warn("retail gold insert failed", "source", src.Name(), "karat", karat, "err", err)
			}
		}
	}
	log.Info("retail gold ingest complete", "sources", len(sources.RegisteredRetailGoldSources))
}
