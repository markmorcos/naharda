# Tasks — add-analytics

## Slice 1 — Deploy Umami
- [x] Umami `deployment.yaml` (chart values) + deploy workflow/dispatch; analytics.naharda.com host.
- [ ] Provision the `umami` database + `naharda-analytics` secret (APP_SECRET, DATABASE_URL). (operator)
- [ ] Create the website entry in Umami; set PUBLIC_UMAMI_SRC + PUBLIC_UMAMI_WEBSITE_ID. (operator)

## Slice 2 — Instrument the dashboard
- [x] Cookieless tracking `<script>` in `Base.astro` (defer, env-gated website-id; omitted until set).
- [x] `data-umami-event` on: subscribe, per-family details, city-select, theme-toggle, lang-toggle, docs/status/privacy nav.
- [ ] Verify pageviews + visitors + countries + events appear; no cookies set. (operator, post-deploy)

## Slice 3 — Docs
- [x] Reference cookieless analytics in `/privacy` (no raw IPs, no cookies).
