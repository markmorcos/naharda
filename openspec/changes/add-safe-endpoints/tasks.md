# Tasks — add-safe-endpoints

After `add-bootstrap`. One family per slice.

## Slice 1 — Cities & calendar (no external deps)
- [ ] `/v1/cities` handler over `domain/cities.go` (13 cities).
- [ ] `/v1/calendar` Hijri↔Gregorian via `go-hijri`; cache 86400s.

## Slice 2 — Prayer times
- [ ] `prayer_times` table; Aladhan ingester with explicit `method`; daily precompute for 13 cities.
- [ ] `/v1/prayer-times/{city}`; today 3600s, past immutable; past dates never recomputed.

## Slice 3 — Weather & AQI
- [ ] `weather_current` + `aqi_current` tables; Open-Meteo ingesters (~10-min cadence).
- [ ] `/v1/weather/{city}` + `/v1/aqi/{city}`; cache 600s; Open-Meteo attribution.

## Slice 4 — Fuel (manual)
- [ ] `fuel_prices` table; seed current prices via `manual_override` with `effective_from`.
- [ ] `/v1/fuel` handler (80/92/95 + diesel); cache 86400s; EGPC attribution.
- [ ] Document the manual-update procedure (run on the next announced price change).
