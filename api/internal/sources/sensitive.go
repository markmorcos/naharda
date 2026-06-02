package sources

import "context"

// 🟡 SENSITIVE SOURCES (project.md §2 source posture, §12, §16 #1).
//
// The interfaces and aggregation engine ship now; concrete implementations are
// registered only after a human signs off on the exact scrape targets (ToS +
// legal sensitivity around parallel-rate publication). All implementations use
// the honest User-Agent and degrade gracefully (a failed/garbage fetch yields
// an error, never a bogus number). Ingest stays a no-op unless
// SENSITIVE_SOURCES_ENABLED=true AND a registry below is non-empty.

// ParallelFXSource yields a single parallel-market EGP-per-USD quote.
type ParallelFXSource interface {
	Name() string
	FetchUSD(ctx context.Context) (float64, error)
}

// RetailGoldSource yields Egypt-retail EGP-per-gram prices by karat (incl. the
// masna3eya premium — kept separate from world_derived, never merged §4).
type RetailGoldSource interface {
	Name() string
	FetchPerGram(ctx context.Context) (map[int]float64, error)
}

// RegisteredParallelSources holds the approved parallel USD/EGP sources
// (register-parallel-fx-sources, §16 #1). All three are independent publishers
// — no aggregator that re-serves another, so n ≥ 2 and {min,avg,max} reflects a
// genuine cross-source spread.
var RegisteredParallelSources = []ParallelFXSource{
	egcurrencyParallel{},
	blackMarketLiveParallel{},
	sarfEGPParallel{},
}

// RegisteredRetailGoldSources is empty until sources are approved (§16 #1).
var RegisteredRetailGoldSources = []RetailGoldSource{}
