# Design — add-dashboard-polish

## Count-up
A tiny island upgrades the existing live-update handler: instead of swapping the number instantly,
animate from the old to the new value with `requestAnimationFrame` over ~400ms, formatting with the
same tabular Plex-Mono figures (so width never changes → no CLS). `prefers-reduced-motion` → instant
swap (no animation). Applies to the hero on `/` and `/ar`.

## Stats strip
A small component fetched server-side from `/v1/stats` (SSR, edge-cached like the page), shown near
the API section: "X requests served · Y data points · FX updated Z". Localized (EN/ar-EG). Makes the
§10 public stats visible and adds a trust signal.

## 0-CLS fonts
Generate a metric-matched fallback `@font-face` (`size-adjust` + `ascent/descent-override`) for IBM
Plex Sans/Arabic (Fontaine or hand-computed) so the webfont swap shifts nothing; `preload` the hero
font. This is the last piece to lock Lighthouse CLS at 0 (a ranking factor — ties to add-seo-coverage).

## Decisions
1. **Enhance the existing SSE handler** for count-up (no new connection) — motion stays "= freshness".
2. Stats strip is **SSR + cached**, not a client fetch — keeps it cheap and SEO-visible.
