# add-analytics

> Cookieless, self-hosted analytics for page views and button clicks. Cites `project.md` §2.3
> (lightweight), §3/§12 (no PII / GDPR — no cookie banner), §9.4 (observability).

## Why

We want stats — which pages get viewed, which buttons get clicked, rough geography — but the
dashboard HTML is now edge-cached (so server `usage_log` misses most pageviews) and we deliberately
don't store raw IPs (§3/§12). A cookieless, self-hosted analytics tool fills the gap without a
consent banner. This was specced in the dashboard design and deferred.

## What changes

- **Self-hosted Umami** in k3s (its own small DB, or a `umami` database on the existing Postgres),
  fronted by Cloudflare.
- **Tracking script** in `Base.astro` (cookieless, ~2KB), loaded from the Umami host.
- **Custom events** on key buttons via `data-umami-event`: subscribe, details-click per family,
  copy-curl, city-select, theme-toggle, docs/status nav.

## Scope

In: the Umami deployment + DB, the script tag, the event attributes, the deploy wiring.

## Non-goals

- Storing raw IPs (GDPR — §3/§12; Umami hashes for unique counts + geo only).
- Google Analytics or any cookie-based tracker (would need a banner).
- A custom analytics backend (Umami is the boring, durable choice).

## Acceptance criteria

- [ ] Page views + unique visitors (cookieless) + countries + referrers visible in Umami.
- [ ] Tracked button events fire (subscribe, details, copy-curl, city-select, theme).
- [ ] No cookies set; no consent banner; raw IPs not stored.

## Dependencies

Builds on `dashboard` (Base.astro) + the deploy pipeline. Reuses Postgres (or a small dedicated DB).
