# Tasks — add-seo-coverage

## Slice 1 — Sitemap + robots
- [x] `@astrojs/sitemap` (per-lang, daily lastmod); `robots.txt` (allow + sitemap link).

## Slice 2 — Tier-1 intent pages (generated)
- [x] Content collection of intent descriptors (currencies + karats); one EN + one ar-EG template.
- [x] Generate `/eur|sar|aed|kwd|gbp-to-egp`, `/black-market-dollar`, `/gold-18k`, `/gold-24k` (+ `/ar/*`).
- [x] Each: answer-first H1, live number, explainer, 30-day sparkline, API CTA.

## Slice 3 — Structured data + linking
- [x] FAQPage + Dataset + BreadcrumbList JSON-LD per page.
- [x] Internal link graph (home ⇄ intent pages ⇄ related); reciprocal hreflang.
- [ ] Lighthouse SEO ≥95 check. (verify post-deploy)

## Slice 4 — OG images (stretch)
- [ ] Dynamic OG image (satori) with the current number per intent page. (deferred — stretch)
