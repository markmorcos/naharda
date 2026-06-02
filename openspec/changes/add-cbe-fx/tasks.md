# Tasks — add-cbe-fx

## Slice 1 — CBE source
- [ ] `internal/sources/cbe.go`: goquery scraper of the CBE rates page → EGP-per-unit per quote, honest UA.
- [ ] Robust parsing (guard structure changes; zero-rows → alert + skip, don't publish garbage).

## Slice 2 — Wire as canonical
- [ ] Seed `sources`: CBE `canonical=true`, exchangerate-api `canonical=false` (cross-check).
- [ ] FX ingest: store official from CBE; keep exchangerate-api as the disagreement cross-check.
- [ ] Disagreement >threshold → prefer CBE + flag in `meta`; outlier guard unchanged; fail-soft.
- [ ] Attribution → "Central Bank of Egypt"; verify `/v1/fx` provenance.
