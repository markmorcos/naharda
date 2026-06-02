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

## Fetch model (refined during implementation)
For the **current-read** families (weather, AQI, prayer-times) v1 fetches **on-demand** from the
upstream and lets the edge cache absorb load (§2.3 "cache at edge, compute rarely") — no cron
ingesters or per-family tables. The Cache-Control windows do the heavy lifting (weather/AQI 600s;
prayer today 3600s, past immutable 1y). Outbound requests carry an honest `User-Agent` + contact
link (§9.6). Calendar is pure local compute. Fuel reads `domain.DefaultFuelPrices`, overridable
per product via the bootstrap `manual_override` table. Stored, *polled* ingest is reserved for the
genuinely-aggregated FX/gold families (add-fx-official-and-gold-world).

## Provenance & attribution (§2.5, §2.11)
Each response's `meta.sources[]` and `meta.attribution` name the upstream (Aladhan, Open-Meteo,
EGPC) with a `fetched_at`. Open-Meteo and Aladhan permit commercial redistribution under
attribution (§12).

## Tables
None added in this change. Prayer/weather/AQI are on-demand; cities + calendar are code/compute;
fuel uses code defaults optionally overridden via the existing `manual_override` table.
