# Tasks — register-parallel-fx-sources

## Slice 1 — Scrapers
- [x] `internal/sources/parallel_egcurrency.go`: `goquery` scraper of egcurrency.com black-market
      USD/EGP → single EGP-per-USD float, honest UA, sane-band validation, error (not 0) on miss.
- [x] `internal/sources/parallel_blackmarketlive.go`: same for en.blackmarketlive.org/egp/usd/.
- [x] `internal/sources/parallel_sarfegp.go`: same for sarfegp.com black-market page (buy/sell mid).
- [x] Verify each selector against the live page (2026-06-02: egcurrency 52.78, blackmarketlive
      53.27, sarfegp 54.72 — parsers run green against the live HTML).

## Slice 2 — Register + seed + threshold
- [x] Add the three to `RegisteredParallelSources` in `internal/sources/sensitive.go`.
- [x] Seed migration `00009_seed_parallel_fx_sources.sql`: insert the three into `sources`
      (`family='fx'`, `canonical=false`, `outlier_threshold=8.0`), idempotent + `Down`.
- [x] `ParallelFXRun` reads each source's `outlier_threshold` from the `sources` table
      (`store.SourceThreshold`), falling back to 8.0 when absent; hardcoded constant is now the fallback.

## Slice 3 — Verify
- [x] Per-source attribution (name + URL + `fetched_at`) flows into `/v1/fx.meta.sources` — handler
      appends each parallel quote's source; URLs added to `sourceURLs`.
- [x] Unit-test parsing against HTML fixtures (`parallel_test.go`): sane-band + error-on-miss paths.
- [ ] Deploy-time check with `SENSITIVE_SOURCES_ENABLED=true`:
      `curl /v1/fx | jq '.data.parallel'` → `{min,avg,max,n,sources[]}`, `n` ≥ 2. (Needs a DB +
      live network; can't run in CI. Fail-soft is covered by the per-source error path + unit tests.)
- [x] Fail-soft by construction: each scraper errors (never returns 0) on a layout change, and
      `ParallelFXRun` `continue`s past a failing source — `n` drops, `official` and the response stay intact.
