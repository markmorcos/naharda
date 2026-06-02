# Design — add-analytics

## Why Umami (vs Plausible / Cloudflare Web Analytics)
- **Umami**: cookieless, self-host-friendly, generous, first-class **custom events** (we need button
  clicks), owns the data. ← chosen
- Plausible: similar, nicer funnels, slightly heavier self-host.
- Cloudflare Web Analytics: free + zero infra, but **no custom events** (pageviews only) — insufficient.

## Deployment
```
  umami (Node)  →  Postgres (a `umami` DB on m720q, or a small dedicated instance)
  k3s deployment + service + ingress (analytics.naharda.com or a path)  ·  Cloudflare in front
  env: DATABASE_URL (umami db), APP_SECRET (secret)
```
Mirrors the existing chart deploy pattern (`deployment.yaml` + dispatch). Tiny resource profile.

## Instrumentation
```
  Base.astro <head>:  <script defer src="https://<umami-host>/script.js"
                        data-website-id="<id>"></script>   (cookieless, ~2KB)
  buttons:  data-umami-event="subscribe" | "details-fx" | "copy-curl" | "city-select" | "theme"
```
The script is the only client addition; it does not block render and adds no cookies.

## Privacy
Umami hashes IP+UA+salt for unique-visitor counts and derives country/city geo — **no raw IP
stored**, no cookies → no consent banner (§3/§12). Reference it in the `/privacy` page
(add-privacy-gdpr).
