package sources

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const sarfEGPURL = "https://sarfegp.com/en/us-dollar-to-egp-black-market/"

// sarfEGPParallel scrapes sarfegp.com's black-market USD/EGP buy/sell.
type sarfEGPParallel struct{}

func (sarfEGPParallel) Name() string { return "sarfegp.com" }

func (s sarfEGPParallel) FetchUSD(ctx context.Context) (float64, error) {
	body, err := fetchParallelBody(ctx, sarfEGPURL)
	if err != nil {
		return 0, err
	}
	defer body.Close()
	return parseSarfEGP(body)
}

// parseSarfEGP reads the black-market buy/sell cells (`#usd-egp-buy`,
// `#usd-egp-sell`) and returns their mid — falling back to whichever single
// side is present. Errors if neither is a sane value. Separated for testing.
func parseSarfEGP(r io.Reader) (float64, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return 0, err
	}
	buy, buyErr := parseRate(doc.Find("#usd-egp-buy").First().Text())
	sell, sellErr := parseRate(doc.Find("#usd-egp-sell").First().Text())
	switch {
	case buyErr == nil && sellErr == nil:
		return (buy + sell) / 2, nil
	case buyErr == nil:
		return buy, nil
	case sellErr == nil:
		return sell, nil
	default:
		return 0, fmt.Errorf("sarfegp: no sane USD/EGP buy/sell rate (layout change)")
	}
}
