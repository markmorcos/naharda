# Tasks — add-dashboard

Status: a **deployable core is built and verified end-to-end** (SSR home + today-grid pulling live
`/v1` data, 2 EN + 2 AR intent pages with JSON-LD + hreflang, email capture, brand system, deploy
pipeline). Remaining items are coverage/polish over the same proven patterns.

## Slice 1 — Foundation (`web/` scaffold + design tokens)
- [x] Astro app (`output: server`, Node adapter, k3s container).
- [x] Design-token CSS (Nile + gold, light/dark via `data-theme` + `prefers-color-scheme`), theme toggle.
- [x] Self-hosted IBM Plex (Sans / Sans Arabic / Mono) via `@fontsource` (GDPR-clean).
- [ ] `size-adjust` fallback @font-face for 0-CLS (Fontaine) — not yet hand-tuned.
- [x] Base layout: logical properties (RTL-ready), `naharda•` wordmark + gold-dot pulse + favicon.

## Slice 2 — Hero + Today grid
- [x] SSR fetch of `/v1` endpoints; live numbers injected into server HTML (verified in view-source).
- [ ] `Cache-Control: max-age=300, stale-while-revalidate` on data routes — not yet set.
- [x] Hero (big tabular Plex-Mono number) + grid cards (FX, gold, fuel, weather, prayer, calendar).
- [x] Gold-dot pulse; fail-soft per card (down endpoint → "—", page intact).
- [ ] Number count-up island — not yet (pulse only).

## Slice 3 — Programmatic intent pages
- [x] `/usd-to-egp` + `/gold-price-egypt` (answer-first H1, explainer, API CTA).
- [x] JSON-LD `FAQPage` on intent pages; `WebSite` on home.
- [ ] Remaining Tier-1 pages (eur/sar/aed/kwd/gbp, black-market-dollar, gold-18/24k) — same pattern.
- [ ] `@astrojs/sitemap`, robots.txt, BreadcrumbList — not yet.

## Slice 4 — Bilingual (ar-EG dialect)
- [x] Astro i18n routing `/ar`, `<html lang="ar-EG" dir="rtl">`.
- [x] Dialect pages: `/ar` home + `/ar/usd-to-egp` (H1s as people search).
- [x] Reciprocal `hreflang` en · ar-EG · x-default on every page.
- [ ] Translate the remaining intent pages.

## Slice 5 — Email capture + stats (+ backend handlers)
- [x] `POST /v1/signups` (single opt-in, consent, honeypot, no email send) — built + tested.
- [x] `GET /v1/stats` (aggregate, PII-free, reads usage_log; 300s) — built + tested.
- [x] Web email-capture form + consent + honeypot → posts to the API.
- [ ] Render `/v1/stats` counters on the dashboard — endpoint ready, UI not wired.

## Slice 6 — Conversion surfaces + content
- [x] API-pitch section with copyable `curl` snippet.
- [ ] README integration examples; thin in-app docs page.
- [ ] Privacy policy + GDPR deletion flow + footer attribution (footer present; policy page not).

## Slice 7 — Analytics, deploy, status
- [ ] Self-hosted Umami/Plausible (cookieless) — not yet.
- [x] `web/Dockerfile` (distroless nodejs24), `web/deployment.yaml` (services schema), `deploy-web.yaml`.
- [x] `deploy-naharda-web` dispatch type registered (infrastructure PR #9).
- [ ] Cloudflare DNS/cache for naharda.com; status link; Lighthouse ≥95 gate — deploy-time.
