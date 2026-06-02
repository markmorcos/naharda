# Tasks — add-cbe-fx

## Slice 1 — CBE source
- [x] `internal/sources/cbe.go`: goquery scraper of the CBE rates page → EGP-per-unit per quote, honest UA.
- [x] Robust parsing (guard structure changes; zero-rows → alert + skip, don't publish garbage).

## Slice 2 — Wire as canonical
- [x] Seed `sources`: CBE `canonical=true`, exchangerate-api `canonical=false` (cross-check).
- [x] FX ingest: store official from CBE; keep exchangerate-api as the disagreement cross-check.
- [x] Disagreement >threshold → prefer CBE + flag in `meta`; outlier guard unchanged; fail-soft.
- [x] Attribution → "Central Bank of Egypt"; verify `/v1/fx` provenance.
