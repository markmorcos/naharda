# Tasks — add-analytics

## Slice 1 — Deploy Umami
- [ ] `umami` database (on m720q Postgres or a small dedicated DB) + secret (APP_SECRET, DSN).
- [ ] Umami `deployment.yaml` (chart) + deploy workflow/dispatch; Cloudflare host.
- [ ] Create the website entry in Umami; capture the website-id.

## Slice 2 — Instrument the dashboard
- [ ] Cookieless tracking `<script>` in `Base.astro` (defer, website-id).
- [ ] `data-umami-event` on: subscribe, per-family details, copy-curl, city-select, theme-toggle, docs/status nav.
- [ ] Verify: pageviews + unique visitors + countries + the button events appear; no cookies set.

## Slice 3 — Docs
- [ ] Reference cookieless analytics in `/privacy` (no raw IPs, no cookies).
