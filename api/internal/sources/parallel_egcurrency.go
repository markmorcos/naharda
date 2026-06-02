package sources

import (
	"context"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

const egcurrencyURL = "https://egcurrency.com/en/currency/usd-to-egp/blackmarket"

// egcurrencyParallel scrapes egcurrency.com's black-market USD/EGP headline rate.
type egcurrencyParallel struct{}

func (egcurrencyParallel) Name() string { return "egcurrency.com" }

func (s egcurrencyParallel) FetchUSD(ctx context.Context) (float64, error) {
	body, err := fetchParallelBody(ctx, egcurrencyURL)
	if err != nil {
		return 0, err
	}
	defer body.Close()
	return parseEGCurrency(body)
}

// parseEGCurrency reads the prominent headline rate, which the site renders as
// the large `b.d-block` element. Separated from the fetch for unit testing.
func parseEGCurrency(r io.Reader) (float64, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return 0, err
	}
	if v, err := parseRate(doc.Find("b.d-block").First().Text()); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("egcurrency: no sane USD/EGP rate in b.d-block (layout change)")
}
