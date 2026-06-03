package stream

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/markmorcos/naharda/api/internal/store"
)

// BuildSnapshot returns the current headline numbers (official FX + world-derived
// gold) as a compact JSON payload for the stream's snapshot/update events.
func BuildSnapshot(ctx context.Context, st *store.Store) ([]byte, error) {
	official := map[string]float64{}
	if fx, err := st.LatestFXRates(ctx, "official"); err == nil {
		for _, r := range fx {
			official[strings.ToLower(r.Quote)] = round2(r.Value)
		}
	}
	gold := []map[string]any{}
	if rows, err := st.LatestGoldPrices(ctx, "world_derived"); err == nil {
		for _, g := range rows {
			gold = append(gold, map[string]any{"karat": g.Karat, "value_egp": round2(g.ValueEGP)})
		}
	}
	return json.Marshal(map[string]any{
		"fx":   map[string]any{"official": official},
		"gold": map[string]any{"world_derived": gold},
		"ts":   time.Now().UTC().Format(time.RFC3339),
	})
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
