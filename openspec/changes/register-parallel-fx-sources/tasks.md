# Tasks — register-parallel-fx-sources

## Slice 1 — Scrapers
- [ ] `internal/sources/parallel_egcurrency.go`: `goquery` scraper of egcurrency.com black-market
      USD/EGP → single EGP-per-USD float, honest UA, sane-band validation, error (not 0) on miss.
- [ ] `internal/sources/parallel_blackmarketlive.go`: same for en.blackmarketlive.org/egp/usd/.
- [ ] `internal/sources/parallel_sarfegp.go`: same for sarfegp.com black-market page.
- [ ] Verify each selector against the live page (rates ~52–53 EGP as of 2026-06-02) before commit.

## Slice 2 — Register + verify
- [ ] Add the three to `RegisteredParallelSources` in `internal/sources/sensitive.go`.
- [ ] Confirm per-source attribution (name + URL + `fetched_at`) flows into `/v1/fx.meta.sources`.
- [ ] Unit-test parsing against saved HTML fixtures (sane-band + error-on-miss paths).
- [ ] Manual check with `SENSITIVE_SOURCES_ENABLED=true`:
      `curl /v1/fx | jq '.data.parallel'` → `{min,avg,max,n,sources[]}`, `n` ≥ 2.
- [ ] Confirm fail-soft: kill one source → `n` drops, response + `official` intact.
