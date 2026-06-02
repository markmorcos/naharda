# Tasks — add-dashboard-polish

## Slice 1 — Count-up
- [x] Upgrade the live-update island: rAF count-up old→new (~400ms), tabular figures (no CLS).
- [x] `prefers-reduced-motion` → instant swap. Apply on `/` and `/ar`.

## Slice 2 — Public stats strip
- [x] Component fetching `/v1/stats` (SSR, cached): requests served · data points · last-updated.
- [x] Place near the API section; localized EN/ar-EG.

## Slice 3 — 0-CLS fonts
- [x] Metric-matched fallback @font-face (size-adjust/overrides) for Plex Sans/Arabic/Mono.
- [ ] Verify Lighthouse CLS = 0. (verify post-deploy)
