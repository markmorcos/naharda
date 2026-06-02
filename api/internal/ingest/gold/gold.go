// Package gold derives world gold prices from spot × FX × karat
// (project.md §4). The egypt_retail stream is added by add-sensitive-sources.
package gold

import (
	"context"
	"log/slog"
	"time"

	"github.com/markmorcos/naharda/api/internal/quality"
	"github.com/markmorcos/naharda/api/internal/sources"
	"github.com/markmorcos/naharda/api/internal/store"
)

var karats = []int{18, 21, 24}

const outlierThresholdPct = 5.0

// Run fetches spot gold, combines with the latest USD/EGP, and stores the
// world-derived EGP-per-gram price for each karat, applying the outlier guard (§9.5).
func Run(ctx context.Context, st *store.Store, alerter *quality.Alerter, log *slog.Logger) {
	spotUSD, src, err := sources.FetchSpotGoldUSD(ctx)
	if err != nil {
		log.Warn("gold spot fetch failed", "err", err)
		return
	}
	usdEGP, ok, err := st.LatestFXRate(ctx, "official", "USD")
	if err != nil || !ok || usdEGP <= 0 {
		log.Warn("gold derive skipped: no USD/EGP rate yet", "err", err)
		return
	}
	usdPerGram := spotUSD / sources.GramsPerTroyOunce
	for _, k := range karats {
		egpPerGram := usdPerGram * usdEGP * float64(k) / 24.0
		pending := false
		if avg, n, err := st.TrailingAvgGold(ctx, "world_derived", k, time.Hour); err == nil &&
			n >= 3 && quality.IsOutlier(egpPerGram, avg, outlierThresholdPct) {
			pending = true
			alerter.Alert(ctx, "gold world-derived outlier held", map[string]any{
				"karat": k, "value": egpPerGram, "trailing_avg": avg,
			})
		}
		if err := st.InsertGoldPrice(ctx, "world_derived", k, egpPerGram, src.Name+" × FX", pending); err != nil {
			log.Warn("gold insert failed", "karat", k, "err", err)
		}
	}
	_ = st.Notify(ctx, "naharda_updates", `{"family":"gold"}`)
	log.Info("gold ingest complete", "spot_usd_oz", spotUSD, "usd_egp", usdEGP)
}
