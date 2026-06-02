package sources

import (
	"strings"
	"testing"
)

// A fixture resembling the CBE rates table: currency label + buy/sell cells.
const cbeFixture = `
<html><body>
<table>
  <thead><tr><th>Currency</th><th>Buy</th><th>Sell</th></tr></thead>
  <tbody>
    <tr><td>US Dollar</td><td>48.50</td><td>48.60</td></tr>
    <tr><td>Euro</td><td>52.10</td><td>52.30</td></tr>
    <tr><td>Saudi Riyal</td><td>12.90</td><td>12.95</td></tr>
    <tr><td>UAE Dirham</td><td>13.20</td><td>13.25</td></tr>
    <tr><td>Kuwaiti Dinar</td><td>157.00</td><td>157.80</td></tr>
    <tr><td>Pound Sterling</td><td>61.40</td><td>61.70</td></tr>
  </tbody>
</table>
</body></html>`

func TestParseCBE(t *testing.T) {
	rates, err := parseCBE(strings.NewReader(cbeFixture))
	if err != nil {
		t.Fatalf("parseCBE: %v", err)
	}
	want := map[string]float64{
		"USD": 48.55,
		"EUR": 52.20,
		"SAR": 12.925,
		"AED": 13.225,
		"KWD": 157.40,
		"GBP": 61.55,
	}
	for code, exp := range want {
		got, ok := rates[code]
		if !ok {
			t.Errorf("missing %s", code)
			continue
		}
		if d := got - exp; d > 0.001 || d < -0.001 {
			t.Errorf("%s = %v, want %v", code, got, exp)
		}
	}
	if len(rates) != len(want) {
		t.Errorf("got %d rates, want %d", len(rates), len(want))
	}
}

func TestParseCBE_RejectionPage(t *testing.T) {
	// The CBE WAF returns a tiny rejection page with no rate rows; we must error
	// (so the ingest falls back to the cross-check) rather than publish garbage.
	const rejected = `<html><head><title>Request Rejected</title></head>
<body>The requested URL was rejected.</body></html>`
	if _, err := parseCBE(strings.NewReader(rejected)); err == nil {
		t.Fatal("expected error for a page with no currency rows")
	}
}
