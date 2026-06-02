# Tasks — add-status-and-docs

Web-only (no backend changes). Both surfaces in the existing `web/` Astro app.

## Slice 1 — Docs
- [ ] `/docs` overview (base URL, no-auth, JSON-only, rate limits) + `/docs/conventions`
      (envelope, error format, ETag/304, Cache-Control table, attribution/terms).
- [ ] Per-resource pages (content collection): fx, gold, fuel, prayer-times, weather, aqi, calendar,
      cities, stats, stream — each with path/params, a `curl` example, and a committed sample response.
- [ ] Document `GET /v1/stream` (SSE) with an `EventSource` snippet.
- [ ] Titles/descriptions/canonical for SEO; link from the homepage "This, as an API" section.

## Slice 2 — Status
- [ ] `/status` page: server-side checks (`/healthz`, `/readyz`, a sample read per family) → green/
      degraded/down per component, cached ~30–60s.
- [ ] Pull `/v1/stats` (requests served, per-family `last_updated`, data-point counts) into the page.
- [ ] Degraded detection from `meta` freshness (stale beyond cache window → degraded).
- [ ] Link to the existing UptimeKuma status/history.

## Slice 3 — Routing + ops
- [ ] Alias `docs.naharda.com` → `/docs` and `status.naharda.com` → `/status` (Cloudflare redirect/
      alias; or extra ingress hosts if the chart supports multi-host).
- [ ] Verify `/v1/*` unchanged; `/status` cacheable (~60s), `/docs` long-cached.
