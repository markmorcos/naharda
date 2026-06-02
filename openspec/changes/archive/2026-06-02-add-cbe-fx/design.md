# Design — add-cbe-fx

## Source
`internal/sources/cbe.go`: fetch the CBE published exchange-rates page and parse the per-currency
buy/sell (or mid) with `goquery`. Honest `User-Agent` + `abuse@naharda.com` (§9.6). Map CBE's
currency labels to our quote set (USD/EUR/SAR/AED/KWD/GBP). CBE publishes EGP-per-unit-of-foreign,
matching our `value` semantics.

## Wiring
```
  fx ingest (official):
    cbe   = sources.FetchCBE(ctx)         # canonical
    ref   = sources.FetchOfficialFX(ctx)  # exchangerate-api, cross-check (kept)
    per quote:
      if |cbe - ref| / cbe > disagreement_threshold → log + meta flag; SERVE CBE (canonical)
      outlier guard as today; store market='official', source='CBE'
  seed sources: CBE (canonical=true), exchangerate-api (canonical=false)
```
Attribution becomes "FX: Central Bank of Egypt" (CBE permits redistribution under attribution §12).

## Robustness
- CBE structure may change (it's a scrape) — guard parsing, alert on zero-rows, fail-soft to
  last-good (§2.6). The exchangerate-api cross-check is the safety net if CBE breaks.
- Cadence: a few times/day (CBE moves slowly); pairs with `add-fx-cadence` if more frequent.

## Decisions
1. Keep exchangerate-api as a **cross-check**, not removed — gives the disagreement guard a second
   opinion and a fallback if the CBE scrape breaks.
