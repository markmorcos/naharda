# register-parallel-fx-sources

> The human sign-off `add-sensitive-sources` was waiting on: register concrete parallel
> USD/EGP sources so the `parallel` aggregate can go live. Cites §4 (parallel = aggregate
> range), §2 (source posture, honesty), §9.5 (outlier/disagreement), §12 (Legal), §16 #1
> (the source-approval gate this change closes).

## Why

`add-sensitive-sources` shipped the parallel-FX machinery — interface, aggregation, flag,
graceful degradation — but left `RegisteredParallelSources` deliberately empty pending human
sign-off on the exact scrape targets (§16 #1). The parallel/black-market dollar is the moat
(§1); nothing is published until real sources are wired. This change closes that gate with
three independent, low-frequency HTML sources.

## What changes

- **Three `ParallelFXSource` scrapers** (`goquery`, honest UA), each yielding one EGP-per-USD
  quote, registered in `RegisteredParallelSources`:
  - `egcurrency.com` (`/en/currency/usd-to-egp/blackmarket`)
  - `blackmarketlive.org` (`/egp/usd/`)
  - `sarfegp.com` (`/en/us-dollar-to-egp-black-market/`)
- **Source independence**: the three are independent publishers, so `n` ≥ 2 and `{min,avg,max}`
  is meaningful. `api-bank.fex.to` is **excluded** — it re-serves egcurrency and would
  double-count.
- **No engine changes**: the 8% outlier guard, aggregation, attribution, and fail-soft already
  exist. This change only registers sources + adds the scrapers.
- **Flag stays off in config by default**; turning `SENSITIVE_SOURCES_ENABLED=true` is the
  operator's deploy-time decision (this is the sign-off).

## Scope

In: the three scraper implementations, their registration, per-source attribution metadata,
verification that each selector matches the live page.

## Non-goals

- Changing the `/v1/fx` envelope or the aggregation/outlier engine (already shipped).
- Retail-gold sources (`RegisteredRetailGoldSources` stays empty — separate sign-off).
- Investing.com USD/EGPp (ToS forbids scraping — explicitly rejected).
- High-frequency scraping (§3 — keep the `@every 30m` default).

## Acceptance criteria

- [ ] With the flag on, `/v1/fx.parallel` returns `{min, avg, max, n, sources[]}` with `n` ≥ 2.
- [ ] Each scraper carries the honest UA + contact and fail-soft: one source failing degrades
      only `n`/that quote, never the response (§2.6).
- [ ] Sources are independent (no fex.to/egcurrency double-count); attribution names each site.
- [ ] An out-of-band quote (>8% off the trailing avg) is held for review, not published.

## Dependencies

After `add-sensitive-sources` (the parallel machinery + flag). This change **is** the §16 #1
human source sign-off; merging it records that approval.
