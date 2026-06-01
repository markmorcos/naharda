package sources

import (
	"context"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// FXQuotes is the Tier-1 currency set (matches the dashboard SEO focus).
var FXQuotes = []string{"USD", "EUR", "SAR", "AED", "KWD", "GBP"}

// FetchOfficialFX returns EGP-per-unit rates for the tracked quote currencies.
//
// NOTE: the working upstream is exchangerate-api.com (market reference rates).
// Production must wire the Central Bank of Egypt (CBE) as the canonical
// `official` source (§4, §12). Attribution always names the source actually used.
func FetchOfficialFX(ctx context.Context) (map[string]float64, domain.Source, error) {
	var raw struct {
		Result string             `json:"result"`
		Rates  map[string]float64 `json:"rates"`
	}
	if err := getJSON(ctx, "https://open.er-api.com/v6/latest/USD", &raw); err != nil {
		return nil, domain.Source{}, err
	}
	egpPerUSD := raw.Rates["EGP"]
	out := make(map[string]float64, len(FXQuotes))
	if egpPerUSD > 0 {
		for _, q := range FXQuotes {
			if qPerUSD := raw.Rates[q]; qPerUSD > 0 {
				out[q] = egpPerUSD / qPerUSD // EGP per 1 unit of quote
			}
		}
	}
	src := domain.Source{
		Name:      "exchangerate-api.com",
		URL:       "https://www.exchangerate-api.com",
		FetchedAt: time.Now().UTC(),
	}
	return out, src, nil
}
