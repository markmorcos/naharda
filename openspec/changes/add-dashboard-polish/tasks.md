# Tasks — add-dashboard-polish

## Slice 1 — Count-up
- [ ] Upgrade the live-update island: rAF count-up old→new (~400ms), tabular figures (no CLS).
- [ ] `prefers-reduced-motion` → instant swap. Apply on `/` and `/ar`.

## Slice 2 — Public stats strip
- [ ] Component fetching `/v1/stats` (SSR, cached): requests served · data points · last-updated.
- [ ] Place near the API section; localized EN/ar-EG.

## Slice 3 — 0-CLS fonts
- [ ] Metric-matched fallback @font-face (size-adjust/overrides) for Plex Sans/Arabic; preload hero font.
- [ ] Verify Lighthouse CLS = 0.
