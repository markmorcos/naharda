# Tasks — add-fx-official-and-gold-world

After `add-bootstrap`.

## Slice 1 — Schema & data-quality wiring
- [x] `fx_rates` + `gold_prices` tables (immutable history, `pending_review` column).
- [x] Seed `sources` rows (CBE, spot/gold feed) with `canonical` + thresholds.
- [x] Activate the outlier-guard + disagreement logic against the bootstrap skeleton.

## Slice 2 — Official FX
- [x] CBE ingester (few-times/day); write `market=official` rows with provenance.
- [x] `/v1/fx` handler — `official` populated, `parallel` aggregate present-but-empty; `max-age=300`.
- [x] `/v1/fx/history` — immutable, year-long cache.

## Slice 3 — World-derived gold
- [x] Gold-compute ingester: `spot × FX × karat` for 18/21/24; `stream=world_derived`.
- [x] `/v1/gold` handler — `world_derived` populated, `egypt_retail` present-but-empty; `max-age=600`.
- [x] `/v1/gold/history` — immutable, year-long cache.

## Slice 4 — Quality scenarios
- [x] Verify >5% spike → `pending_review` + alert + serve last-good.
- [x] Verify FX cross-source disagreement >2% → prefer canonical + `meta` flag.
