# Design — add-safe-endpoints

## Endpoints, sources & caching (§9.2)

| Endpoint | Source | Cache (max-age) | Notes |
|---|---|---|---|
| `/v1/cities` | static `domain/cities.go` | 86400s | 13 cities, hardcoded lat/lon |
| `/v1/calendar` | local compute (`go-hijri`) | 86400s | Hijri↔Gregorian, no external call |
| `/v1/prayer-times/{city}` | Aladhan API (free, commercial-OK) | today 3600s · past 31536000s (immutable) | explicit `method` code |
| `/v1/weather/{city}` | Open-Meteo (free, commercial-OK) | 600s | current conditions |
| `/v1/aqi/{city}` | Open-Meteo Air-Quality | 600s | AQI by city |
| `/v1/fuel` | **manual entry** → `manual_override` | 86400s | products 80/92/95 + diesel; `effective_from` |

## Fuel = 🟢, manual-first (resolves explore gap #1)
Egyptian pump prices are **government-set and announced** (Petroleum Products Pricing Committee /
EGPC) — public administrative facts, commercial-OK, **zero political sensitivity** (unlike parallel
FX, this *is* the state's own number). They change only ~1–4×/year on pre-announced dates. Building
a fragile scraper for a quarterly event is wasted effort, so v1 fuel is **hand-entered via the
`manual_override` table** with `effective_from`. Attribution: "Egyptian Ministry of Petroleum / EGPC."

## Ingest cadences
Weather/AQI: every ~10 min (matches 600s cache). Prayer-times: precompute the day's cities once
daily. Calendar: pure compute per request (cacheable). Fuel: event-driven (manual).

## Provenance & attribution (§2.5, §2.11)
Every stored value records `source` + `fetched_at`; each response's `meta.sources[]` and
`meta.attribution` name the upstream (Aladhan, Open-Meteo, EGPC). Open-Meteo and Aladhan permit
commercial redistribution under attribution (§12).

## Tables (per-family, created here)
`prayer_times` (city, date, method, fajr…isha), `weather_current` (city, fetched_at, fields),
`aqi_current` (city, fetched_at, aqi + pollutants), `fuel_prices` (product, value, effective_from,
source). Cities are static; calendar is computed (no table).
