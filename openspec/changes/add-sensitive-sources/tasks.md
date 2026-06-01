# Tasks — add-sensitive-sources

After `add-fx-official-and-gold-world`. **Ships dark; flag stays off until human sign-off.**

## Slice 0 — Human gate (blocking)
- [ ] `TODO(ask)`: shortlist 2–3 parallel-FX sources + retail-gold sources; get explicit sign-off (§16 #1).
- [ ] Record approved sources in `sources` (with thresholds + attribution) — flag still off.

## Slice 1 — Flag & scaffolding
- [ ] `SENSITIVE_SOURCES_ENABLED` env flag (default false); endpoints expose empty 🟡 fields when off.
- [ ] Compliance: honest `User-Agent` + `abuse@naharda.com`; low-frequency scheduler entries.

## Slice 2 — Parallel FX
- [ ] `goquery` scrapers for approved sources; aggregate to `{ min, avg, max, n, sources[] }`.
- [ ] Populate `/v1/fx.parallel` (never a single value); `max-age=300`; immutable history.
- [ ] Graceful degradation: failed scrape empties only `parallel` with a `meta` flag.

## Slice 3 — Egypt-retail gold
- [ ] Scraper for retail gold (masna3eya premium); `stream=egypt_retail`, never merged.
- [ ] Populate `/v1/gold.egypt_retail`; `max-age=600`; immutable history.

## Slice 4 — Verify gated behavior
- [ ] Flag off → fields empty, nothing scraped. Flag on (post-approval) → ranges populate.
- [ ] Each aggregated value still passes the outlier guard.
