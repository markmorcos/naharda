# Tasks — add-fx-official-and-gold-world

After `add-bootstrap`.

## Slice 1 — Schema & data-quality wiring
- [ ] `fx_rates` + `gold_prices` tables (immutable history, `pending_review` column).
- [ ] Seed `sources` rows (CBE, spot/gold feed) with `canonical` + thresholds.
- [ ] Activate the outlier-guard + disagreement logic against the bootstrap skeleton.

## Slice 2 — Official FX
- [ ] CBE ingester (few-times/day); write `market=official` rows with provenance.
- [ ] `/v1/fx` handler — `official` populated, `parallel` aggregate present-but-empty; `max-age=300`.
- [ ] `/v1/fx/history` — immutable, year-long cache.

## Slice 3 — World-derived gold
- [ ] Gold-compute ingester: `spot × FX × karat` for 18/21/24; `stream=world_derived`.
- [ ] `/v1/gold` handler — `world_derived` populated, `egypt_retail` present-but-empty; `max-age=600`.
- [ ] `/v1/gold/history` — immutable, year-long cache.

## Slice 4 — Quality scenarios
- [ ] Verify >5% spike → `pending_review` + alert + serve last-good.
- [ ] Verify FX cross-source disagreement >2% → prefer canonical + `meta` flag.
