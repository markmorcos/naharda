package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/markmorcos/naharda/api/internal/domain"
)

// cbeURL is the Central Bank of Egypt published exchange-rates page (§4 canonical).
const cbeURL = "https://www.cbe.org.eg/en/economic-research/statistics/exchange-rates"

// cbeCurrencyNames maps CBE's English currency labels to our quote codes. Match
// is case-insensitive substring, so minor label changes ("US Dollar" vs
// "U.S. Dollar") still resolve.
var cbeCurrencyNames = map[string]string{
	"us dollar":      "USD",
	"euro":           "EUR",
	"saudi riyal":    "SAR",
	"uae dirham":     "AED",
	"kuwaiti dinar":  "KWD",
	"pound sterling": "GBP",
	"sterling":       "GBP",
}

var numberRe = regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)

// FetchCBE fetches and parses the CBE exchange-rates page into EGP-per-unit
// rates per tracked quote. It is the canonical `official` source; callers keep a
// reference source as a cross-check and fail-soft to last-good if CBE is
// unreachable or returns nothing parseable.
//
// NOTE: the CBE site sits behind a WAF that may reject datacenter IPs; in that
// case this returns an error and the ingest falls back to the cross-check
// source. The parser is tolerant of layout changes (it locates currencies by
// label and reads the numeric cells), and returns an error on zero rows so we
// never publish garbage.
func FetchCBE(ctx context.Context) (map[string]float64, domain.Source, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cbeURL, nil)
	if err != nil {
		return nil, domain.Source{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		return nil, domain.Source{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, domain.Source{}, fmt.Errorf("CBE GET: status %d", resp.StatusCode)
	}
	rates, err := parseCBE(resp.Body)
	if err != nil {
		return nil, domain.Source{}, err
	}
	src := domain.Source{
		Name:      "Central Bank of Egypt",
		URL:       cbeURL,
		FetchedAt: time.Now().UTC(),
	}
	return rates, src, nil
}

// parseCBE extracts EGP-per-unit rates from the CBE rates HTML. Separated from
// the HTTP fetch so it can be unit-tested against a fixture. It scans every
// table row, identifies the currency by its label, and uses the mid of the
// numeric cells (buy/sell) as the rate. Returns an error if nothing parses.
func parseCBE(r io.Reader) (map[string]float64, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	out := make(map[string]float64)
	doc.Find("tr").Each(func(_ int, row *goquery.Selection) {
		// Read each cell separately: concatenating the row text would merge
		// adjacent numbers (e.g. "48.50"+"48.60") and corrupt parsing.
		var cells []string
		row.Find("td, th").Each(func(_ int, cell *goquery.Selection) {
			cells = append(cells, strings.TrimSpace(cell.Text()))
		})
		rowText := strings.ToLower(strings.Join(cells, " "))
		code := ""
		for name, c := range cbeCurrencyNames {
			if strings.Contains(rowText, name) {
				code = c
				break
			}
		}
		if code == "" {
			return
		}
		if _, done := out[code]; done {
			return // first match wins (avoids duplicate label rows)
		}
		var vals []float64
		for _, cell := range cells {
			m := numberRe.FindString(strings.ReplaceAll(cell, ",", ""))
			if m == "" {
				continue
			}
			if v, err := strconv.ParseFloat(m, 64); err == nil && v > 0 {
				vals = append(vals, v)
			}
		}
		if len(vals) == 0 {
			return
		}
		// Use the mid of the first two numerics (buy/sell); fall back to the
		// single available value.
		var rate float64
		if len(vals) >= 2 {
			rate = (vals[0] + vals[1]) / 2
		} else {
			rate = vals[0]
		}
		out[code] = rate
	})
	if len(out) == 0 {
		return nil, fmt.Errorf("CBE parse: no currency rows found (layout change or WAF block)")
	}
	return out, nil
}
