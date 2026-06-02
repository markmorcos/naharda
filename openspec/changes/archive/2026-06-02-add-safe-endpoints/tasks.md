# Tasks — add-safe-endpoints

After `add-bootstrap`. One family per slice. (Current-read families fetch on-demand + edge-cache —
see design.md "Fetch model".)

## Slice 1 — Cities & calendar (no external deps)
- [x] `/v1/cities` handler over `domain/cities.go` (13 cities); cache 86400s.
- [x] `/v1/calendar` Hijri↔Gregorian via `go-hijri` (Umm al-Qura), optional `?date=`; cache 86400s.

## Slice 2 — Prayer times
- [x] Aladhan client with explicit Egyptian method (code 5); `?date=` support.
- [x] `/v1/prayer-times/{city}` on-demand; today 3600s, past immutable 1y; unknown city → 404.

## Slice 3 — Weather & AQI
- [x] Open-Meteo clients (weather + air-quality) with honest UA.
- [x] `/v1/weather/{city}` + `/v1/aqi/{city}` on-demand; cache 600s; Open-Meteo attribution.

## Slice 4 — Fuel (manual)
- [x] `domain.DefaultFuelPrices` seed (80/92/95 + diesel) with `effective_from`.
- [x] `/v1/fuel` handler prefers `manual_override` per product; cache 86400s; EGPC attribution.
- [x] Manual-update procedure documented (operator records new prices via manual_override).
