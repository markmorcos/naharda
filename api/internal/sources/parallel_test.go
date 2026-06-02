package sources

import (
	"strings"
	"testing"
)

// Fixtures mirror the real markup each source uses for its headline rate (the
// element the scraper targets), plus a "layout changed" page with no usable
// value to exercise the error-not-garbage guard.

const egcurrencyFixture = `<html><body>
  <div class="converter">
    <b class="fs-5 lh-1 text-nowrap text-danger">52.78</b>
    <b class="d-block lh-1 mb-1 text-danger" style="font-size: clamp(4rem,5vw,8rem);">52.78</b>
  </div>
  <b class="fs-5">52.25</b>
</body></html>`

const blackMarketLiveFixture = `<html><body>
  <h1>USD to EGP Black Market Rate Today</h1>
  <p class="price-black wp-block-paragraph"> 53.27 </p>
  <table><tr><td class="has-text-align-center">Bank Price</td><td>51.93</td></tr></table>
</body></html>`

const sarfEGPFixture = `<html><body>
  <table>
    <tr><td id="usd-egp-buy">54.44</td><td id="usd-egp-sell">54.99</td></tr>
    <tr><td id="usd-egp-bank-buy">53.57</td><td id="usd-egp-bank-sell">53.67</td></tr>
  </table>
</body></html>`

// noRate has the wrong structure (a layout change / rejection): no targeted
// element carries a sane-band number.
const noRate = `<html><body><h1>Request Rejected</h1><p>The page changed.</p></body></html>`

func approx(got, want float64) bool { d := got - want; return d < 0.001 && d > -0.001 }

func TestParseEGCurrency(t *testing.T) {
	v, err := parseEGCurrency(strings.NewReader(egcurrencyFixture))
	if err != nil {
		t.Fatalf("parseEGCurrency: %v", err)
	}
	if !approx(v, 52.78) {
		t.Errorf("got %v, want 52.78", v)
	}
	if _, err := parseEGCurrency(strings.NewReader(noRate)); err == nil {
		t.Error("expected error on layout change")
	}
}

func TestParseBlackMarketLive(t *testing.T) {
	v, err := parseBlackMarketLive(strings.NewReader(blackMarketLiveFixture))
	if err != nil {
		t.Fatalf("parseBlackMarketLive: %v", err)
	}
	if !approx(v, 53.27) {
		t.Errorf("got %v, want 53.27", v)
	}
	if _, err := parseBlackMarketLive(strings.NewReader(noRate)); err == nil {
		t.Error("expected error on layout change")
	}
}

func TestParseSarfEGP(t *testing.T) {
	// Mid of buy (54.44) and sell (54.99).
	v, err := parseSarfEGP(strings.NewReader(sarfEGPFixture))
	if err != nil {
		t.Fatalf("parseSarfEGP: %v", err)
	}
	if !approx(v, 54.715) {
		t.Errorf("got %v, want 54.715", v)
	}
	if _, err := parseSarfEGP(strings.NewReader(noRate)); err == nil {
		t.Error("expected error on layout change")
	}
}

// parseRate must reject out-of-band numbers (page furniture) so a stray count
// or year never becomes a published rate.
func TestParseRateSaneBand(t *testing.T) {
	for _, bad := range []string{"", "abc", "2026", "0.5", "1000000"} {
		if _, err := parseRate(bad); err == nil {
			t.Errorf("parseRate(%q) = nil err, want rejection", bad)
		}
	}
	if v, err := parseRate(" 53.27 EGP"); err != nil || !approx(v, 53.27) {
		t.Errorf("parseRate(53.27) = %v, %v", v, err)
	}
}
