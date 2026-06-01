# add-safe-endpoints

> The 🟢 safe-source families — no ToS/scraping risk. Cites `project.md` §2 (source posture),
> §4 (Domain Model), §9.2 (caching), §9.4 (observability).

## Why

These are the low-risk, commercial-use-OK families that prove data quality and give the dashboard
something to show on day one, with zero legal exposure (§12). They unblock `add-dashboard`'s grid.

## What changes

Read endpoints + their ingesters/compute and per-family data tables:

- **`/v1/cities`** — the 13 canonical cities (static; hardcoded lat/lon in `domain/cities.go`).
- **`/v1/calendar`** — Hijri ↔ Gregorian, computed locally (`hablullah/go-hijri`), no external call.
- **`/v1/prayer-times/{city}`** — daily salah times with explicit `method` code.
- **`/v1/weather/{city}`** + **`/v1/aqi/{city}`** — current conditions / AQI from Open-Meteo.
- **`/v1/fuel`** — pump prices (gasoline 80/92/95, diesel) via **manual entry** (🟢, see design).

## Scope

In: the six 🟢 endpoints above, their ingesters (or compute), per-family tables, cache policies,
and per-source attribution.

## Non-goals

- Any 🟡 source (parallel FX, retail gold — `add-sensitive-sources`).
- FX or gold (`add-fx-official-and-gold-world`).
- A fuel **scraper** — v1 fuel is manual-entry by design (changes ~1–4×/year).
- Coptic calendar (deferred per §8).

## Acceptance criteria

- [ ] All six endpoints return the standard envelope with correct `Cache-Control` per §9.2.
- [ ] `/v1/calendar` computes locally with no network dependency.
- [ ] Prayer-times responses include the explicit `method` code; past dates are immutable.
- [ ] Weather + AQI cite Open-Meteo; fuel cites Ministry of Petroleum/EGPC with `effective_from`.

## Dependencies

After `add-bootstrap` (middleware, schema substrate, `manual_override`, cities list).
