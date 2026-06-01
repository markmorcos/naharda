# Tasks — add-dashboard

Vertical slices, top to bottom. Each maps to a small PR. Built **after** `add-bootstrap` and the
FX/gold ingest changes (so headline numbers exist).

## Slice 1 — Foundation (`web/` scaffold + design tokens)
- [ ] Astro app in `web/` with Cloudflare-fronted hybrid rendering (Node adapter, k3s container).
- [ ] Open Props + design-token CSS (`--canvas/--ink/--brand/--accent/--up/--down/--official/--parallel`),
      light + dark via token swap, `prefers-color-scheme` aware + toggle island.
- [ ] Self-hosted IBM Plex (Sans / Sans Arabic / Mono), subset-per-script, woff2, preload,
      `size-adjust` fallback @font-face (Fontaine) → 0 CLS. Verify Lighthouse CLS = 0.
- [ ] Base layout: logical properties (RTL-ready), `naharda•` wordmark + gold-dot favicon.

## Slice 2 — Hero + Today grid (the `/v1/today` answer)
- [ ] SSR fetch `http://naharda-api.naharda.svc/v1/today`; inject numbers **into server HTML**.
- [ ] Set `Cache-Control: public, max-age=300, stale-while-revalidate=86400`.
- [ ] Hero: big tabular Plex-Mono number + delta colour + "updated N min ago" + sources count.
- [ ] Data grid cards: FX (official + parallel range), gold (3 karats × 2 streams), fuel, prayer,
      weather, AQI. `--official` solid vs `--parallel` dashed-range treatment.
- [ ] Live-tick island (rAF count-up); gold-dot pulse on refresh; `prefers-reduced-motion` honored.
- [ ] Verify: numbers present in "view source" (not JS-only); Lighthouse ≥95 mobile.

## Slice 3 — Programmatic intent pages (the SEO engine)
- [ ] Content collection → generate Tier-1 pages: `/usd-to-egp`, `/black-market-dollar`,
      `/eur|sar|aed|kwd|gbp-to-egp`, `/gold-price-egypt`, `/gold-{18,21,24}k`.
- [ ] Per page: answer-first H1, official-vs-parallel explainer, 30-day history sparkline, FAQ, API CTA.
- [ ] Tier-2: `/prayer-times/{13 cities}`. Tier-3 shown on dashboard only (no dedicated SEO push).
- [ ] JSON-LD: FAQPage + Dataset + BreadcrumbList per page; WebSite+SearchAction+Organization on home.
- [ ] `@astrojs/sitemap`, robots.txt, self-canonical, stable `<title>` (number in H1/description only).

## Slice 4 — Bilingual (ar-EG dialect SEO pages)
- [ ] Astro i18n routing: `/ar/...` with `<html lang="ar-EG" dir="rtl">`.
- [ ] Egyptian-dialect content for Tier-1 intent pages (H1s as people actually search).
- [ ] `hreflang` en · ar-EG · x-default, reciprocal; verify RTL layout + bidi number wrapping.

## Slice 5 — Email capture + stats (+ backend handlers)
- [ ] `api/`: `POST /v1/signups` (single opt-in, consent, honeypot, IP rate-limit, no email send).
- [ ] `api/`: `GET /v1/stats` (aggregate, PII-free, reads `usage_log`); add `300s` to §9.2 cache table.
- [ ] Web: email-capture form + consent checkbox + link to privacy policy + deletion-request path.
- [ ] Web: render `/v1/stats` (public counters) on the dashboard.

## Slice 6 — Conversion surfaces + content
- [ ] API-pitch section: copyable `curl` snippet + free-tier CTA on every intent page.
- [ ] README with integration examples; thin in-app docs page (no `docs.naharda.com` in v1).
- [ ] Privacy policy + GDPR deletion flow; footer attribution (§9 — required on every surface).

## Slice 7 — Analytics, deploy, status
- [ ] Self-hosted Umami/Plausible in-cluster; wire pageviews (no cookies, no banner).
- [ ] `web/Dockerfile` (distroless nodejs22), `web/deployment.yaml` (chart 0.6.0, ~25m/64Mi→250m/128Mi).
- [ ] `.github/workflows/deploy-web.yaml` + add `deploy-naharda-web` dispatch type in
      `markmorcos/infrastructure`.
- [ ] Cloudflare: DNS + cache rules for `naharda.com`; public status page link (`status.naharda.com`).
- [ ] Final gate: Lighthouse ≥95 (Perf/SEO/A11y/BP) mobile; numbers in HTML; zero cookie banner;
      fonts self-hosted.
