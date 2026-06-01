package sources

import (
	"context"
	"time"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// GramsPerTroyOunce converts troy ounces to grams.
const GramsPerTroyOunce = 31.1034768

// FetchSpotGoldUSD returns the spot gold price in USD per troy ounce.
func FetchSpotGoldUSD(ctx context.Context) (float64, domain.Source, error) {
	var raw struct {
		Price float64 `json:"price"`
	}
	if err := getJSON(ctx, "https://api.gold-api.com/price/XAU", &raw); err != nil {
		return 0, domain.Source{}, err
	}
	src := domain.Source{Name: "gold-api.com", URL: "https://gold-api.com", FetchedAt: time.Now().UTC()}
	return raw.Price, src, nil
}
