# Design — add-status-and-docs

Both surfaces live in the **existing `web/` Astro app** (one deployable, brand consistency, SEO),
read-only against the live API. No backend changes.

## Docs (`/docs`)

```
  /docs               overview: base URL https://api.naharda.com · no auth (IP-limited) · JSON only
  /docs/conventions   the envelope { data, meta } · error format · ETag/304 · Cache-Control table ·
                      rate limits (60/min, 1000/day) · attribution + free-tier non-commercial (§9)
  /docs/[resource]    per family: path, params, a curl example, a sample JSON response
                      (fx · gold · fuel · prayer-times · weather · aqi · calendar · cities · stats · stream)
```

- **Format:** hand-written Astro/MD pages (content collection), one per resource. `curl` examples
  use `api.naharda.com`; sample responses are **static** (committed) so docs are stable and don't
  depend on live data at render time. (Optionally a "live response" island could fetch a real sample,
  but static is the v1 default.)
- **Stream:** document `GET /v1/stream` (SSE) with an `EventSource` snippet — the dev-facing half of
  add-live-updates.
- **SEO:** docs are content too (rank for "naharda api", "egypt fx api", "egypt prayer times api").
  Proper titles/descriptions/canonical; linked from the homepage "This, as an API" section.

## Status (`/status`)

```
  /status
    ├ Components (live):  API · Database · each data family
    │     server-side checks: GET /healthz, /readyz, and a sample read per family → green/degraded/down
    ├ Freshness (from /v1/stats):  per-family last_updated, requests served, data-point counts
    └ History/incidents:  link to UptimeKuma (status.naharda.com history) — the §9.4 source of truth
```

- **Live checks run server-side** (the status pod calls the API), reduced to a green/amber/red per
  component. **Cached ~30–60s** (`Cache-Control` + a short in-memory/edge cache) so a burst of
  status views doesn't hammer the API — the page is itself a cacheable read.
- **UptimeKuma stays the uptime/incident source** (§9.4). `/status` is a branded *live snapshot* +
  freshness, with a clear link to UptimeKuma for 90-day history and incident timelines. We do NOT
  rebuild incident management.
- A `degraded` state (e.g. a 🟡 source down, or a family stale beyond its cache window) is shown
  per-component using the §9.5 freshness signals already in `meta`.

## Decisions

1. **One app, not a separate docs site** ✅ — `/docs` + `/status` in `web/` (brand, one deploy,
   shared tokens). A standalone Starlight site was considered and rejected for v1.
2. **UptimeKuma for history, branded `/status` for the live snapshot** ✅ — best of both; no
   reinventing incident tracking.
3. **Subdomain routing — to confirm:** simplest is `status.naharda.com` / `docs.naharda.com` as
   **Cloudflare redirects/aliases** to `naharda.com/status` and `/docs` (no chart/ingress changes).
   Alternative: add extra ingress hosts to the web service (needs chart multi-host support). Lean:
   Cloudflare alias.
4. **OpenAPI spec** — deferred. A machine-readable `/v1/openapi.json` (+ client generation) is a
   strong dev-first follow-up but out of scope here.

## Cross-cutting / non-breaking
- Read-only against existing endpoints; **no `/v1/*` change**. `/status` and `/docs` are cacheable
  (status ~60s, docs long) — consistent with §2.3.
- Both inherit the brand (Nile + gold, Plex, dark mode) from the existing `Base` layout.
