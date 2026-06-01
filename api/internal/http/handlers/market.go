package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
	"github.com/markmorcos/naharda/api/internal/http/respond"
	"github.com/markmorcos/naharda/api/internal/store"
)

// FX returns official EGP rates; the parallel aggregate is present but empty
// until add-sensitive-sources (§4).
func (h *Handlers) FX(w http.ResponseWriter, r *http.Request) {
	rows, err := h.store.LatestFXRates(r.Context(), "official")
	if err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "FX data unavailable.", 30)
		return
	}
	official := make(map[string]float64, len(rows))
	var latest time.Time
	source := ""
	for _, rt := range rows {
		official[strings.ToLower(rt.Quote)] = round2(rt.Value)
		if rt.FetchedAt.After(latest) {
			latest = rt.FetchedAt
		}
		source = rt.Source
	}
	parallel := emptyParallel()
	if h.sensitiveEnabled {
		if quotes, err := h.store.LatestParallelQuotes(r.Context(), "USD"); err == nil && len(quotes) > 0 {
			parallel = aggregateParallel(quotes)
		}
	}
	data := map[string]any{
		"base":     "EGP",
		"official": official,
		"parallel": parallel,
	}
	respond.JSON(w, r, 300, data, domain.Meta{
		Sources:     oneSource(source, latest),
		Attribution: "FX: exchangerate-api.com reference rate (CBE official wiring is a production follow-up). Free tier non-commercial.",
	})
}

// FXHistory returns immutable history for a quote (?quote=usd&limit=). 1y cache.
func (h *Handlers) FXHistory(w http.ResponseWriter, r *http.Request) {
	quote := strings.ToUpper(r.URL.Query().Get("quote"))
	if quote == "" {
		quote = "USD"
	}
	rows, err := h.store.FXHistory(r.Context(), "official", quote, parseLimit(r))
	if err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "FX history unavailable.", 30)
		return
	}
	respond.JSON(w, r, 31536000, map[string]any{"quote": quote, "history": rows}, domain.Meta{
		Attribution: "FX history via Naharda. Free tier non-commercial.",
	})
}

// Gold returns the world-derived stream; egypt_retail is present but empty (§4).
func (h *Handlers) Gold(w http.ResponseWriter, r *http.Request) {
	rows, err := h.store.LatestGoldPrices(r.Context(), "world_derived")
	if err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "Gold data unavailable.", 30)
		return
	}
	world := make([]map[string]any, 0, len(rows))
	var latest time.Time
	source := ""
	for _, g := range rows {
		world = append(world, map[string]any{"karat": g.Karat, "value_egp": round2(g.ValueEGP)})
		if g.FetchedAt.After(latest) {
			latest = g.FetchedAt
		}
		source = g.Source
	}
	retail := []map[string]any{}
	if h.sensitiveEnabled {
		if rows, err := h.store.LatestRetailGold(r.Context()); err == nil {
			for _, g := range rows {
				retail = append(retail, map[string]any{"karat": g.Karat, "value_egp": round2(g.ValueEGP)})
			}
		}
	}
	data := map[string]any{
		"unit":          "EGP per gram",
		"world_derived": world,
		"egypt_retail":  retail, // never merged with world_derived (§4)
	}
	respond.JSON(w, r, 600, data, domain.Meta{
		Sources:     oneSource(source, latest),
		Attribution: "Gold (world-derived): gold-api.com spot × FX. Free tier non-commercial.",
	})
}

// GoldHistory returns immutable history for a karat (?karat=21&limit=). 1y cache.
func (h *Handlers) GoldHistory(w http.ResponseWriter, r *http.Request) {
	karat := 21
	if k, err := strconv.Atoi(r.URL.Query().Get("karat")); err == nil && k > 0 {
		karat = k
	}
	rows, err := h.store.GoldHistory(r.Context(), "world_derived", karat, parseLimit(r))
	if err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "unavailable", "Gold history unavailable.", 30)
		return
	}
	respond.JSON(w, r, 31536000, map[string]any{"karat": karat, "history": rows}, domain.Meta{
		Attribution: "Gold history via Naharda. Free tier non-commercial.",
	})
}

func emptyParallel() map[string]any {
	return map[string]any{"min": nil, "avg": nil, "max": nil, "n": 0, "sources": []any{}}
}

// aggregateParallel reduces per-source quotes to {min,avg,max,n,sources} — never
// a single value (§4, Decision Log: honesty about uncertainty).
func aggregateParallel(quotes []store.FXRate) map[string]any {
	if len(quotes) == 0 {
		return emptyParallel()
	}
	min, max, sum := quotes[0].Value, quotes[0].Value, 0.0
	srcs := make([]string, 0, len(quotes))
	for _, q := range quotes {
		if q.Value < min {
			min = q.Value
		}
		if q.Value > max {
			max = q.Value
		}
		sum += q.Value
		srcs = append(srcs, q.Source)
	}
	return map[string]any{
		"min":     round2(min),
		"avg":     round2(sum / float64(len(quotes))),
		"max":     round2(max),
		"n":       len(quotes),
		"sources": srcs,
	}
}

var sourceURLs = map[string]string{
	"exchangerate-api.com": "https://www.exchangerate-api.com",
	"gold-api.com × FX":    "https://gold-api.com",
}

func oneSource(name string, at time.Time) []domain.Source {
	if name == "" {
		return []domain.Source{}
	}
	return []domain.Source{{Name: name, URL: sourceURLs[name], FetchedAt: at}}
}

func parseLimit(r *http.Request) int {
	limit := 100
	if n, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && n > 0 {
		limit = n
	}
	if limit > 1000 {
		limit = 1000
	}
	return limit
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
