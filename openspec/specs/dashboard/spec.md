# dashboard Specification

## Purpose
TBD - created by archiving change add-dashboard. Update Purpose after archive.
## Requirements
### Requirement: Live data SHALL be present in server-rendered HTML
The dashboard MUST inject current data values into the HTML on the server (hybrid SSR), so that
search-engine crawlers and no-JS clients see real, current numbers without executing JavaScript.
Rendered HTML MUST be edge-cacheable per the project caching policy (`max-age=300`,
`stale-while-revalidate`).

#### Scenario: Crawler reads the homepage source
- **WHEN** a client requests `naharda.com` and inspects the raw response body
- **THEN** the current FX and gold values are present in the markup before any JavaScript runs

#### Scenario: Stale-while-revalidate keeps the page fast and fresh
- **WHEN** the edge cache holds a copy older than `max-age` but within the SWR window
- **THEN** the cached HTML is served instantly and revalidated in the background via SSR

### Requirement: The dashboard SHALL meet a strict performance & lightness budget
Pages MUST score ≥ 95 on mobile Lighthouse (Performance, SEO, Accessibility, Best-Practices),
ship no UI framework on the critical path, self-host all fonts (subset per script), and introduce
no layout shift from fonts or motion.

#### Scenario: Audit a data page
- **WHEN** a mobile Lighthouse audit runs against any intent page
- **THEN** all four categories score ≥ 95 and Cumulative Layout Shift is ≤ 0.1

### Requirement: Intent pages SHALL exist per query, bilingual with structured data
The site MUST publish one page per high-value query intent (Tier-1 FX + gold prioritized) in both
English and Egyptian-Arabic (`ar-EG`), each answer-first, with reciprocal `hreflang`
(`en` / `ar-EG` / `x-default`), `dir="rtl"` on Arabic pages, and JSON-LD (`FAQPage`, `Dataset`,
`BreadcrumbList`).

#### Scenario: Arabic dollar-rate query
- **WHEN** a user in Egypt searches the dialect query for the dollar rate today
- **THEN** the `/ar/usd-to-egp` page is eligible to rank with its current rate, FAQ structured data,
  and a correct `hreflang` link to its English counterpart

### Requirement: The dashboard SHALL respect privacy by default
Analytics MUST be cookieless and self-hosted (no consent banner required), and every data surface
MUST display source attribution per the project standards.

#### Scenario: First visit
- **WHEN** a user loads any page for the first time
- **THEN** no cookie-consent banner is shown, no third-party analytics cookie is set, and source
  attribution is visible

