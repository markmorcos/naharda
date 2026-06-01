# add-dashboard

> The `web/` marketing front door: a feather-light, bilingual (EN + Egyptian-Arabic), SEO-first
> dashboard that answers _"what's happening in Egypt today?"_ and sells the API.
> Cites `project.md` §2 (Principles 9, 12), §5 (Architecture), §6 (Tech Stack), §7 (Repo),
> §8 (v1 scope), §9 (Standards), §10 (Hooks), §11 (Deployment), §12 (GDPR).

## Why

`project.md` §2.12 — _"The dashboard sells the API. Marketing surface first, product second."_ —
and §8's v1 success test: _"the dashboard ranks for Egyptian-currency queries."_ The single
highest-value query on the internet for this product — **"الدولار بكام النهاردة"** (the
parallel-market dollar rate) — has no good incumbent. This change builds the surface that captures
that intent and converts curious visitors into API users.

## What changes

A new `web/` deployable (Astro) plus the two small backend endpoints it depends on:

- **`web/`** — static-first, hybrid-SSR dashboard: hero (`/v1/today` answer), data grid, per-intent
  SEO landing pages, API pitch, email capture, status link.
- **`api/`** — `GET /v1/stats` (public aggregate metrics) and `POST /v1/signups` (email capture).
  Their backing tables (`usage_log`, `signups`) ship earlier in `add-bootstrap`; this change adds
  the **handlers**.
- **Deploy** — `web/Dockerfile` (distroless nodejs), `web/deployment.yaml` (chart 0.6.0),
  `.github/workflows/deploy-web.yaml`, and the `deploy-naharda-web` dispatch type in the
  `markmorcos/infrastructure` repo.

## Scope

In: brand/design system, Astro app, hybrid-SSR rendering, EN UI + ar-EG SEO pages, programmatic
intent-pages (FX-parallel + gold prioritized), structured data, self-hosted analytics, the two API
endpoints, the web deploy pipeline, README integration examples, public status page link.

## Non-goals (v1)

- Login / accounts / API-key UI (auth is v2 — §8).
- Full Arabic UI localization (v1 = English UI + Arabic *SEO landing pages* only).
- Fighting weather/prayer SERPs (shown on the dashboard, **not** chased for ranking — Tier 3).
- Confirmation emails / double opt-in (no Resend until v2 — §10); v1 is single opt-in.
- A standalone `docs.naharda.com` site (deferred; a thin in-app docs page only).
- Cloudflare Pages hosting (stay in k3s; revisit only if global TTFB is measured to matter).

## Acceptance criteria

- [ ] `naharda.com` live, serving the `/v1/today` snapshot with **live numbers present in the
      server HTML** (verifiable via "view source", not just after JS runs).
- [ ] Lighthouse ≥ 95 on Performance, SEO, Accessibility, Best-Practices (mobile).
- [ ] ar-EG intent pages exist for USD/EUR/SAR/AED/KWD/GBP + gold, with `hreflang` + JSON-LD.
- [ ] Email capture writes to `signups`; `/v1/stats` renders on the dashboard.
- [ ] Deploys via `deploy-naharda-web` through the infrastructure pipeline.
- [ ] Zero cookie banner (no-cookie analytics); fonts self-hosted (GDPR — §12).

## Dependencies

Built **after** `add-bootstrap` (schema + middleware + api deploy) and at least
`add-fx-official-and-gold-world` (so the headline numbers exist). `add-sensitive-sources` (parallel
FX) makes the wedge pages compelling but can land independently behind its flag.
