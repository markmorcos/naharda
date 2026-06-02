# add-seo-coverage

> Complete the programmatic SEO surface so the dashboard can rank — sitemap, robots, and the full
> Tier-1 intent pages. Cites `project.md` §1 (the parallel-FX wedge), §8 (v1 success test: "ranks
> for Egyptian-currency queries"), and the dashboard design's SEO plan.

## Why

§8's v1 success test is that the dashboard **ranks for Egyptian-currency queries**. The dashboard
core shipped only 2 intent pages (USD, gold) and no sitemap/robots — Google can't discover or rank
the surface. This adds the rest of the Tier-1 intent pages (the actual ranking engine) plus the
technical SEO basics.

## What changes

- **Sitemap + robots**: `@astrojs/sitemap` (per-lang, daily lastmod) and `robots.txt`.
- **Tier-1 intent pages** (EN + ar-EG), each answer-first with live number, FAQ + Dataset +
  BreadcrumbList JSON-LD, history sparkline, API CTA:
  `/eur-to-egp` `/sar-to-egp` `/aed-to-egp` `/kwd-to-egp` `/gbp-to-egp` `/black-market-dollar`
  `/gold-18k` `/gold-24k` (+ `/ar/*`).
- **BreadcrumbList** JSON-LD on all intent pages; internal link graph from home/intent pages.
- **Dynamic OG image** with the current number (satori) — stretch.

## Scope

In: sitemap, robots, the listed intent pages (bilingual) generated from a content collection,
breadcrumb JSON-LD, internal linking. Tier-3 (weather/fuel) stay off the SEO battlefield (design).

## Non-goals

- Chasing weather/prayer SERPs (Tier-3 — shown, not ranked).
- Backlink/off-page SEO (out of scope for code).

## Acceptance criteria

- [ ] `sitemap-index.xml` + `robots.txt` served; every public page listed with `hreflang`.
- [ ] All listed Tier-1 intent pages exist EN + ar-EG with live data + FAQ/Dataset/Breadcrumb JSON-LD.
- [ ] Lighthouse SEO ≥ 95 on the new pages; bidirectional EN↔AR hreflang.

## Dependencies

Builds on `dashboard` (the intent-page pattern, Base, RateCard) and `fx`/`gold` history endpoints.
