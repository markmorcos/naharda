# Design — add-fx-official-and-gold-world

## FX — official (CBE)
`/v1/fx` returns the §4 two-market shape, but only `official` is populated here:
```
  { official: { usd, eur, … , source, fetched_at },
    parallel: { min, avg, max, n, sources[] }   ← present, EMPTY until add-sensitive-sources }
```
Source: Central Bank of Egypt (🟢, commercial-OK under attribution §12). Ingest cadence ~ a few
times/day (official rates move slowly). Cache `max-age=300` (§9.2); history immutable (`31536000s`).

## Gold — world-derived
`world_derived = spot_gold × FX(EGP) × karat_factor` for karats **18/21/24**. Tagged
`stream=world_derived`; the `egypt_retail` stream stays empty here and the two **never merge** (§4 —
they diverge by the local "masna3eya" premium). Cache `max-age=600`; history immutable.

## Data-quality machinery — now LIVE (§9.5, uses the add-bootstrap skeleton)
- **Outlier guard**: a value >5% (per-source `outlier_threshold`) off the trailing-1h average →
  written `pending_review` + alert; API serves **last-good** + `meta` staleness flag (fail-soft §2.6).
- **Cross-source disagreement**: when cross-checked FX sources differ > 2% (per-source
  `disagreement_threshold`), log, **prefer the `canonical` source**, and flag in `meta`.
- **Alerting**: slog WARN + `ALERT_WEBHOOK_URL` ping.
- **Manual override**: an operator can pin a value when an ingester is broken (precedence within window).

## Tables (immutable history — §2.4)
```
  fx_rates     (market='official', base='EGP', quote, value, source, fetched_at, pending_review)
  gold_prices  (stream='world_derived', karat, value_egp, source, fetched_at, pending_review)
```
Corrections are new rows, never mutations. History endpoints (`/v1/fx/history`, `/v1/gold/history`)
serve immutable rows with year-long cache.

## Provenance & attribution
Every row + response names its source (CBE, the spot/gold feed) and `fetched_at`; `meta.attribution`
is always present (§2.5, §2.11).

## Implementation note (sources)
The v1 build wires working free upstreams and attributes them honestly:
- **FX official** → `open.er-api.com` (exchangerate-api.com **market reference** rate). The `official`
  market slot is populated by this as a stand-in; **wiring the Central Bank of Egypt (CBE) as the
  canonical `official` source is a production follow-up** (CBE is 🟢/approved per §2, §12). The
  attribution string names the actual source used.
- **Gold spot** → `api.gold-api.com` (XAU USD/oz); `world_derived` = spot/31.1035 × USD-EGP × karat/24.
- Ingest cadence: FX `@every 1h`, gold `@every 15m`; both primed once on startup. The outlier guard
  needs ≥3 in-window samples before it can hold a value (no baseline → no false holds).
