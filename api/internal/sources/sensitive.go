package sources

import "context"

// 🟡 SENSITIVE SOURCES (project.md §2 source posture, §12, §16 #1).
//
// The interfaces and aggregation engine ship now; the concrete source
// implementations are intentionally NOT registered until a human signs off on
// the exact scrape targets (ToS + legal sensitivity around parallel-rate
// publication). Until then the registries are empty and ingest is a no-op even
// when SENSITIVE_SOURCES_ENABLED=true. All implementations MUST use the honest
// User-Agent (getJSON / a goquery client) and degrade gracefully.

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

// RegisteredParallelSources is empty until sources are approved (§16 #1).
var RegisteredParallelSources = []ParallelFXSource{}

// RegisteredRetailGoldSources is empty until sources are approved (§16 #1).
var RegisteredRetailGoldSources = []RetailGoldSource{}
