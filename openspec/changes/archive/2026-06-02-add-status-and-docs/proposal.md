# add-status-and-docs

> Two dev-facing surfaces the v1 dashboard core deferred: a public **API docs** page and a public
> **status** page. Cites `project.md` §1 (dev-first API), §8 (v1 done-when: docs + status page),
> §9.1/§9.2/§9.3 (envelope, caching, rate limits — the contract to document), §9.4 (UptimeKuma +
> `status.naharda.com`), §11 (`docs.naharda.com`, `status.naharda.com` hosts).

## Why

§8's v1 "done" list explicitly includes a **public status page** and **integration docs**, and §11
reserves `docs.naharda.com` + `status.naharda.com`. The product is a *dev-first* API (§1) — without
docs, the dashboard sells something a developer can't actually adopt; without a status page, B2B
trust has nothing to point at. Both are additive web surfaces; neither requires backend changes.

## What changes

- **Docs** — a `/docs` section in the existing `web/` app: overview (base URL, no-auth, rate
  limits), the response envelope + error format + ETag/caching reference, a per-endpoint reference
  (FX, gold, fuel, prayer, weather, AQI, calendar, cities, stats, **stream**) with `curl` + sample
  JSON, and attribution/ToS. `docs.naharda.com` aliases to it.
- **Status** — a branded `/status` page in `web/`: live component health (server-side `/healthz`,
  `/readyz` + a sample read per family), `/v1/stats` (requests served, per-family freshness), and a
  link to the existing **UptimeKuma** for historical uptime/incidents. `status.naharda.com` aliases
  to it.

## Scope

In: the `/docs` content + the `/status` page in `web/`, the subdomain aliasing, brief server-side
caching of the status checks, and SEO metadata for the docs pages.

## Non-goals

- A separate docs **site/framework** (e.g. standalone Starlight) — keep one `web/` deployable.
- A machine-readable **OpenAPI spec** / "try it" console — documented as a follow-up, not v1.
- Rebuilding uptime history / incident management — UptimeKuma remains the source of truth (§9.4).
- Any backend/API change (docs + status read existing endpoints only).

## Acceptance criteria

- [ ] `/docs` lists every public endpoint with a working `curl` example + a sample response, plus the
      envelope, error, rate-limit and caching reference; reachable at `docs.naharda.com`.
- [ ] `/status` shows live per-component health + `/v1/stats` freshness, and links to UptimeKuma;
      reachable at `status.naharda.com`.
- [ ] The status checks are cached briefly so a burst of status pageviews doesn't hammer the API.
- [ ] Docs pages carry proper titles/descriptions/canonical and are crawlable (dev-discovery SEO).
- [ ] No change to any `/v1/*` endpoint.

## Dependencies

Builds on shipped capabilities: `public-api` / `api-core` (the contract being documented), `dashboard`
(the `web/` app + brand), and `add-dashboard`'s deploy pipeline. UptimeKuma already exists (§9.4).
