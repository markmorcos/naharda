package sources

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const blackMarketLiveURL = "https://en.blackmarketlive.org/egp/usd/"

// blackMarketLiveParallel scrapes en.blackmarketlive.org's USD/EGP rate.
type blackMarketLiveParallel struct{}

func (blackMarketLiveParallel) Name() string { return "blackmarketlive.org" }

func (s blackMarketLiveParallel) FetchUSD(ctx context.Context) (float64, error) {
	body, err := fetchParallelBody(ctx, blackMarketLiveURL)
	if err != nil {
		return 0, err
	}
	defer body.Close()
	return parseBlackMarketLive(body)
}

// parseBlackMarketLive reads the semantic `.price-black` element carrying the
// black-market figure. Separated from the fetch for unit testing.
func parseBlackMarketLive(r io.Reader) (float64, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return 0, err
	}
	if v, err := parseRate(doc.Find(".price-black").First().Text()); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("blackmarketlive: no sane USD/EGP rate in .price-black (layout change)")
}
