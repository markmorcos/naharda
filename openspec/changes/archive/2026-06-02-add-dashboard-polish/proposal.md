# add-dashboard-polish

> The remaining dashboard finish: number count-up, a public stats strip, and 0-CLS font fallback.
> Cites the dashboard design (motion = freshness; tabular figures; CWV), §10 (public stats).

## Why

The dashboard core shipped functional but left a few design items: the hero number "ticks" on live
update (pulse only — no count-up), `/v1/stats` exists but isn't surfaced on the dashboard, and the
`size-adjust` fallback font (for guaranteed 0-CLS) wasn't hand-tuned. These are the polish that
makes it feel finished and hits the CWV/SEO target.

## What changes

- **Count-up animation** on the hero number when it changes (rAF, tabular figures → no layout
  shift; `prefers-reduced-motion` respected) — the "feels alive" finish on top of the SSE update.
- **Public stats strip** on the dashboard from `/v1/stats` (requests served, data points,
  last-updated) — the §10 public stats, made visible.
- **0-CLS fonts**: `size-adjust`/`ascent-override` fallback @font-face (Fontaine) so the webfont
  swap causes no layout shift; preload the hero font.

## Scope

In: the count-up island, the stats strip component, the font-fallback tuning.

## Non-goals

- New data or endpoints (uses existing `/v1/stats` + the SSE stream).
- A full redesign — this is finishing, not re-theming.

## Acceptance criteria

- [ ] On a live update, the hero number counts up smoothly (no CLS); reduced-motion shows an instant swap.
- [ ] A stats strip shows live `/v1/stats` figures on the dashboard.
- [ ] Lighthouse CLS = 0 with the size-adjusted fallback; hero font preloaded.

## Dependencies

Builds on `dashboard` + `live-updates` (SSE) + `public-api` (`/v1/stats`).
