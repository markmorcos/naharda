# Tasks — add-sensitive-sources

After `add-fx-official-and-gold-world`. **Machinery ships now; flag stays OFF and the concrete
scrape sources stay unregistered until human sign-off (§16 #1).** Reuses `fx_rates` (market=parallel)
and `gold_prices` (stream=egypt_retail) — no new tables.

## Slice 0 — Human gate (BLOCKING — not done)
- [ ] `TODO(ask)`: shortlist 2–3 parallel-FX sources + retail-gold sources; get explicit sign-off (§16 #1).
- [ ] Implement the approved `ParallelFXSource` / `RetailGoldSource` scrapers (goquery) and register them.
- [ ] Record approved sources in the `sources` table (thresholds + attribution).
- [ ] Flip `SENSITIVE_SOURCES_ENABLED=true` in deployment once the above is signed off.

## Slice 1 — Flag & scaffolding
- [x] `SENSITIVE_SOURCES_ENABLED` env flag (default false); endpoints expose empty 🟡 fields when off.
- [x] Compliance: honest `User-Agent` + `abuse@naharda.com` (shared sources client); gated scheduler entries.
- [x] Pluggable `ParallelFXSource` / `RetailGoldSource` interfaces + empty registries (await approval).

## Slice 2 — Parallel FX
- [x] Aggregate per-source quotes → `{ min, avg, max, n, sources[] }` (never a single value §4).
- [x] Populate `/v1/fx.parallel` when enabled; `max-age=300`; immutable history (shared fx_rates).
- [x] Graceful degradation: a failed source is skipped; aggregate uses whatever succeeded (§2.6).
- [ ] (gated) concrete goquery scrapers — part of Slice 0.

## Slice 3 — Egypt-retail gold
- [x] Populate `/v1/gold.egypt_retail` per karat when enabled; never merged with world_derived (§4).
- [ ] (gated) concrete retail scraper — part of Slice 0.

## Slice 4 — Verify gated behavior
- [x] Flag off → fields empty, nothing ingested. Flag on + rows present → aggregate/range populate.
      (Verified with synthetic source rows: parallel {min,avg,max,n,sources}, retail per-karat.)
