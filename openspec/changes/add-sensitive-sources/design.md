# Design — add-sensitive-sources

## The gate (non-negotiable)
`SENSITIVE_SOURCES_ENABLED` env flag, **default false**. Code, scrapers, and schema ship dark.
The flag is flipped **only after** a human signs off on the exact 2–3 sources (§12, §16 #1). Until
then `/v1/fx.parallel` and `/v1/gold.egypt_retail` are present-but-empty. This satisfies §8's
"🟡 sources scaffolded behind a feature flag pending human approval."

## Parallel FX — honest aggregate, never a single number (§4, Decision Log)
```
  parallel: { min, avg, max, n, sources[] }     ← always a range; n = source count
```
Scrape 2–3 approved sources (`PuerkitoBio/goquery`), aggregate into the range. If sources disagree
wildly, the spread itself is the signal — we publish the honest min/max, not a fabricated midpoint.
Cache `max-age=300` (matches official FX). History immutable.

## Egypt-retail gold — separate stream
`egypt_retail` includes the local "masna3eya" workmanship premium and therefore **diverges** from
`world_derived`; the two streams are **never merged** (§4). Tagged `stream=egypt_retail`. Cache 600s.

## Compliance posture (§12, §9.6)
- Honest `User-Agent` identifying Naharda + `abuse@naharda.com` contact link on every outbound request.
- **Low frequency** polling (these are "fresh enough", not real-time — §3).
- Per-response attribution naming each source; willingness to remove a source on request.
- Berlin jurisdiction mitigates regulatory risk around parallel-rate publication (§12).

## Graceful degradation (fail-soft §2.6)
A broken 🟡 scrape degrades **only** its field — `parallel` (or `egypt_retail`) goes empty/stale with
a `meta` flag while `official` / `world_derived` and the rest of the response stay intact. The
outlier guard (§9.5) still applies to each aggregated value.

## Tables
Extends `fx_rates` (`market='parallel'`, plus the aggregate inputs + per-source rows for `n`/spread)
and `gold_prices` (`stream='egypt_retail'`). Immutable history; corrections are new rows.
