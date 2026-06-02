# dashboard

## ADDED Requirements

### Requirement: The dashboard SHALL publish a discoverable, complete intent-page surface
The dashboard MUST serve a sitemap and robots.txt, and MUST provide a Tier-1 intent page per
high-value query (USD/EUR/SAR/AED/KWD/GBP to EGP, black-market dollar, gold 18/21/24k) in EN and
ar-EG, each answer-first with the live value in server HTML and FAQ/Dataset/BreadcrumbList JSON-LD,
with reciprocal hreflang.

#### Scenario: Crawler discovers the surface
- **WHEN** a crawler fetches `/sitemap-index.xml` and `/robots.txt`
- **THEN** every public page is listed with its language alternates and the sitemap is linked from robots

#### Scenario: A currency intent page
- **WHEN** a user opens `/eur-to-egp` (or its `/ar/` counterpart)
- **THEN** the page answers the query first with the live rate in the server HTML, plus FAQ/Dataset/Breadcrumb structured data and a link to its language alternate
