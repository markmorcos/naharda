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
`sensitive.ParallelFXRun()` already iterates the registry, stores per-source `market='parallel'`
quotes, and the `/v1/fx` handler aggregates them at read time into `{min,avg,max,n,sources[]}`.

## DB seed + per-source threshold (symmetry with CBE)
- **Seed migration** (`migrations/000NN_seed_parallel_fx_sources.sql`): insert the three into
  `sources` with `family='fx'`, `canonical=false`, and an `outlier_threshold` (default 8.0 — wider
  than official's 5.0, per §9.5). Idempotent `ON CONFLICT DO NOTHING`, with a matching `Down`.
- **Threshold lookup**: `ParallelFXRun` reads each source's `outlier_threshold` from the `sources`
  table (a `store.SourceThreshold(ctx, name)` lookup) instead of the hardcoded
  `outlierThresholdPct = 8.0` constant. Falls back to 8.0 if the row is missing, so behavior is
  unchanged when a source isn't seeded. This makes the guard tunable per source without a redeploy
  and keeps parallel symmetric with how official sources live in the DB.

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
3. **Threshold moves to the DB** — the hardcoded 8% becomes a per-source `outlier_threshold` column
   (default 8.0), so parallel sources are tunable in the registry like official ones. Aggregation,
   fail-soft, and the `/v1/fx` envelope are otherwise unchanged.
