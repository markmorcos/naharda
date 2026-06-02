package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Parallel-market USD/EGP scrapers (register-parallel-fx-sources, §16 #1).
//
// Each source publishes a single black-market USD/EGP figure on an HTML page.
// The scrapers locate it by a specific element and validate it sits in a sane
// band — on a layout change they return an error (never a bogus number) so the
// ingest fail-soft skips the source rather than publishing garbage, mirroring
// cbe.go's zero-rows guard.

// Sane band for a parallel USD/EGP quote. Anything outside is treated as a
// parse/layout error. Wide enough to tolerate real moves, tight enough to
// reject page furniture (counts, years, percentages).
const (
	parallelMinRate = 20.0
	parallelMaxRate = 200.0
)

// fetchParallelBody GETs an HTML page with the honest UA (§9.6). The caller
// must close the returned body.
func fetchParallelBody(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	return resp.Body, nil
}

// parseRate extracts the first numeric token from s and checks it is within the
// parallel sane band. Returns an error otherwise (caller treats as "no quote").
func parseRate(s string) (float64, error) {
	m := numberRe.FindString(strings.ReplaceAll(strings.TrimSpace(s), ",", ""))
	if m == "" {
		return 0, fmt.Errorf("no number in %q", s)
	}
	v, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0, err
	}
	if v < parallelMinRate || v > parallelMaxRate {
		return 0, fmt.Errorf("rate %v outside sane band [%v,%v]", v, parallelMinRate, parallelMaxRate)
	}
	return v, nil
}
