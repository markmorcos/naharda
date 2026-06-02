# Design — register-parallel-fx-sources

## Sources
Each scraper lives in `internal/sources/` and implements `ParallelFXSource`
(`Name() string`, `FetchUSD(ctx) (float64, error)`), mirroring `cbe.go`: honest `User-Agent` +
`abuse@naharda.com`, `Accept: text/html`, the shared `client`, and tolerant `goquery` parsing
that returns an error (never `0`) on a missing/garbage value so the engine fail-soft skips it.

| File | Site | Selector strategy (verify live before commit) |
|------|------|-----------------------------------------------|
| `parallel_egcurrency.go` | egcurrency.com/en/currency/usd-to-egp/blackmarket | Rate is a single prominent heading; locate the price node, regex the first number. |
| `parallel_blackmarketlive.go` | en.blackmarketlive.org/egp/usd/ | "1 USD = X EGP" row in the converter table; read the EGP cell. |
| `parallel_sarfegp.go` | sarfegp.com/en/us-dollar-to-egp-black-market/ | Single rate heading; regex the first number. |

Selectors are scrape-fragile by nature — each parser locates the value by nearby label/context
(not brittle nth-child paths where avoidable) and validates the number is in a sane band
(e.g. 20–120 EGP/USD) before returning, so a layout change yields an error rather than a bogus
quote.

## Wiring
```
internal/sources/sensitive.go:
  var RegisteredParallelSources = []ParallelFXSource{
      egcurrencyParallel{}, blackmarketLiveParallel{}, sarfEGPParallel{},
  }
```
No other change. `sensitive.ParallelFXRun()` already iterates the registry, applies the 8%
outlier guard against the 1h trailing avg, stores per-source `market='parallel'` quotes, and the
`/v1/fx` handler aggregates them at read time into `{min,avg,max,n,sources[]}`.

## Compliance (§2, §12)
- Honest UA + contact on every request (reuse `userAgent`); low frequency (`@every 30m` default).
- Per-response attribution already lists each contributing source by name + `fetched_at`.
- Remove-on-request: dropping a source is a one-line registry edit.

## Robustness
- Independent failure: `FetchUSD` errors are logged and skipped; `n` drops but the aggregate and
  `official` stay intact (§2.6). If all three fail, `parallel` goes empty/stale with a `meta` flag.
- Outlier hold: a quote >8% off the trailing avg is stored `pending_review=true` and excluded from
  the served aggregate (existing engine behavior) — guards against a single site spiking.

## Decisions
1. **Three independent sites, not the fex.to aggregator** — fex.to re-serves egcurrency for Egypt,
   so including both would inflate `n` with a correlated value and defeat the cross-source spread.
2. **Reject Investing.com** despite richer data — its ToS prohibits scraping and it actively blocks
   bots; the honest-posture (§12) rules it out.
3. **No engine/threshold changes** — the 8% parallel guard from `add-sensitive-sources` stands; this
   change is purely source registration.
