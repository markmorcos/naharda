# add-fx-official-and-gold-world

> Official EGP exchange rates (CBE) and world-derived gold — the 🟢 half of the headline numbers,
> plus the point where the data-quality machinery (§9.5) goes live. Cites §4 (Domain Model),
> §9.2 (caching), §9.5 (data quality), §2 (source posture).

## Why

The dashboard's hero answers "الدولار بكام" and "سعر الذهب". The **official** rate and the
**world-derived** gold price are safe, commercial-OK, and computable now — they make the headline
real while the 🟡 parallel/retail wedge is still pending source sign-off.

## What changes

- **`/v1/fx`** — EGP rates with `official` populated (the `parallel` aggregate stays present but
  empty until `add-sensitive-sources`). Source: Central Bank of Egypt.
- **`/v1/gold`** — the `world_derived` stream only: computed `spot × FX × karat` for karats 18/21/24.
  The `egypt_retail` stream stays empty until `add-sensitive-sources`; **streams never merge** (§4).
- **Data-quality machinery live** — outlier guard (>5% off trailing-1h), cross-source disagreement
  (>2% FX), canonical-source preference, held→`pending_review`→serve-last-good, alerting.
- **Tables** — `fx_rates`, `gold_prices` (immutable history).

## Scope

In: official FX ingest + endpoint, world-derived gold compute + endpoint, the applied data-quality
guard, immutable history tables, history endpoints (immutable cache).

## Non-goals

- Parallel FX and Egypt-retail gold (🟡 — `add-sensitive-sources`).
- Merging the two gold streams (forbidden — §4).
- Trading signals/predictions (§3).

## Acceptance criteria

- [ ] `/v1/fx` returns `official` with provenance; the `parallel` aggregate is present but empty.
- [ ] `/v1/gold` returns `world_derived` for 18/21/24k; `egypt_retail` present but empty.
- [ ] An injected >5% spike is held `pending_review`, alerts, and the endpoint serves last-good.
- [ ] FX cross-check disagreement >2% prefers the canonical source and flags it in `meta`.
- [ ] History endpoints return immutable rows with year-long cache.

## Dependencies

After `add-bootstrap` (envelope, `sources`, `manual_override`, guard skeleton).
