# Design — add-seo-coverage

## Intent pages — generated, not hand-written
A content collection of intent descriptors drives the pages (one template, EN + ar-EG):
```
  currency:  { quote: USD/EUR/SAR/AED/KWD/GBP }  → /{quote}-to-egp  + /black-market-dollar (parallel-first)
  gold:      { karat: 18/21/24 }                  → /gold-{karat}k   (21k already exists)
```
Each page: answer-first H1 (the query verbatim), live number (SSR, in HTML), official-vs-parallel
or karat explainer, a 30-day history sparkline (the moat — competitors can't fake it), FAQ, API CTA.
Per §1, the **parallel/black-market** pages are the priority (the wedge).

## Structured data
`FAQPage` (answer-box eligible) + `Dataset` (temporalCoverage + provider) + `BreadcrumbList` on each;
`WebSite` + `SearchAction` already on home. Stable `<title>` (number in H1/description, not title).

## Technical
- `@astrojs/sitemap` (split by lang, daily lastmod) + `robots.txt` (allow, sitemap link).
- Reciprocal `hreflang` (en · ar-EG · x-default) — already the pattern.
- Internal link graph: home ⇄ every intent page ⇄ related currencies + gold.
- **OG image** (stretch): satori-rendered card with the current number for social shares.

## Caching
Intent pages are data pages → `max-age=300, stale-while-revalidate` (the existing middleware).
