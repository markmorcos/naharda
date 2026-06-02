# Design — add-dashboard

The durable decisions from the explore session, organized as four interlocking threads:
**Stack · SEO · Typography · Brand.** Guiding tension resolved throughout:

```
   PERFECT SEO  ←→  SUPER LIGHT  ←→  MODERN / ALIVE
   (numbers IN the HTML, fresh)   (minimal JS)   (motion = data freshness, never decoration)
```

---

## 1. Rendering model & tech stack (resolves Open Q#4)

**Framework: Astro.** Ships 0 KB JS by default; islands opt-in; HTML-first ("boring durable",
§2.1); native View Transitions; content collections drive programmatic pages. Svelte is available
as an *island* if any single widget needs it — but the baseline is zero-framework.

**Rendering: hybrid, one renderer, one source of truth.**

```
  DATA pages (home, intent pages)   → SSR + edge cache + stale-while-revalidate
  STATIC pages (about/privacy/docs) → prerendered (Astro `prerender = true`)

  Cache-Control: public, max-age=300, stale-while-revalidate=86400
    hit fresh  → CF serves cached HTML instantly
    hit stale  → CF serves stale instantly, revalidates in background via SSR
    true miss  → SSR renders (first-ever request / full expiry only)
```

The fresh-number-for-SEO problem (Google must *see* `47.50` in the markup, yet it changes every
5 min) is solved by SSR injecting the number server-side, then edge-caching the result:

```
  request → Cloudflare (HTML cached 300s, SWR) ──miss──► k3s naharda-web (Astro SSR)
                                                              └─ fetch http://naharda-api.naharda.svc/v1/today
                                                                 (in-cluster; API's own 300s cache applies)
                                                              → inject numbers → return → CF caches
```

**Decision — hosting: (A) k3s container, NOT (B) Cloudflare Pages.** Keeps `web/` inside §7/§11
(distroless container, `deployment.yaml`, same pipeline as `api/`), needs no Decision Log
amendment, and the in-cluster API fetch never hits the public internet. SWR gives the
"instant static-feeling default + fresh on demand" behavior **without** a second SSG artifact
(which would be a second, drift-prone source of truth). Flip to (B) only if global TTFB is later
*measured* to matter — and then it's a deliberate Decision Log entry.

**Styling: design tokens (Open Props) + Astro scoped CSS.** The token set *is* the design system;
dark mode = a token swap; near-zero runtime; no framework lock-in. (Tailwind was the considered
alternative; rejected for token-purity + "boring durable".)

**Islands: vanilla / Web Components only.** Live number "tick" (rAF count-up), city selector
(a `<form>` that navigates — 0 JS), copy-curl, theme toggle, mobile nav. No framework unless one
island forces it.

**Motion budget: motion = freshness, never decoration.** Native View Transitions (page→page),
CSS keyframes/transitions, `prefers-reduced-motion` respected. **No Framer Motion / GSAP / Lottie.**
Nothing animates that isn't signaling the data changed. Motion must add ~0 JS and never cause CLS.

**Build/deploy:** multi-stage → `gcr.io/distroless/nodejs22-debian12` (SSR needs a JS runtime,
unlike the API's `distroless/static`); `web/deployment.yaml` (chart 0.6.0), lighter profile
(~25m/64Mi req → ~250m/128Mi limit, mostly idle behind the cache); `deploy-naharda-web` dispatch.

---

## 2. SEO architecture (the "reach people" engine)

Mental model: **a search-answer machine that happens to have a homepage.** One real query → one
URL → answer-first HTML → ranks.

**Intent map, sized by where we can actually win:**

```
  TIER 1 — the wedge, no good incumbent. Concentrate fire.
  ├ /usd-to-egp           official + PARALLEL range  ← highest-value page
  ├ /black-market-dollar  parallel-specific intent
  ├ /eur-to-egp /sar-to-egp /aed-to-egp /kwd-to-egp /gbp-to-egp   (Gulf = remittances)
  └ /gold-price-egypt + /gold-21k /gold-18k /gold-24k             (21k = the default)

  TIER 2 — high volume, decent incumbents; cheap recurring traffic.
  └ /prayer-times/{13 cities}

  TIER 3 — DO NOT fight for SERP (show on dashboard only).
  └ /fuel-prices-egypt · /weather/{city} · /hijri-date-today
```

Pages are **generated from data via content collections**, not hand-written, × {en, ar-EG}.

**Anti-thin-content:** each intent page carries genuine value (your data is the value) — answer-first
sentence, official-vs-parallel explainer, 30-day history sparkline (immutable history competitors
can't fake), FAQ, "get via API" CTA. Substantive by construction, not a doorway page.

**Bilingual (EN UI + ar-EG SEO pages, Egyptian dialect):**

```
  /usd-to-egp  ↔  /ar/usd-to-egp   (Latin slug even on AR — ranks on title/H1/body, stays shareable)
  <html lang=en>      <html lang="ar-EG" dir="rtl">
  "Dollar to EGP"     H1: "الدولار بكام النهاردة؟"   ← dialect = how people actually TYPE/search
  hreflang: en · ar-EG (Egyptian-specific) · x-default, reciprocal on every page
  Western digits (47.50) for SERP/snippet safety.
```

**Four levers:**

```
  1. FRESHNESS (QDF)   FX is THE query-deserves-freshness case. Visible "updated 3 min ago" +
                       <time datetime> + JSON-LD dateModified = fetched_at. (Free — we already carry it.)
  2. STRUCTURED DATA   FAQPage (answer-box eligible) · Dataset (temporalCoverage+provider) ·
                       BreadcrumbList. Home: WebSite+SearchAction + Organization.
                       ⚠ NO fake Product/Offer schema.
  3. CORE WEB VITALS   a ranking factor → WHY "light" = SEO. LCP<2.5s (text hero, no image),
                       CLS<0.1 (tabular figures, size-adjusted fonts), INP<200ms (≈no JS). Target 100.
  4. HYGIENE           @astrojs/sitemap (per-lang, daily lastmod) · robots allow · self-canonical ·
                       STABLE <title> (number lives in H1+description, not the volatile title) ·
                       dynamic OG image w/ current number [v1.1] · strong internal link graph.
```

**Analytics: self-hosted Umami/Plausible in-cluster — NOT GA4.** No cookies → no banner (a banner
is itself a CLS/UX tax), GDPR-clean (§12, Berlin), ~2 KB. Search Console + Bing Webmaster for
ranking data.

---

## 3. Typography

Two heavy jobs: **numbers are the product** (tabular figures, no jitter) and **two scripts, one
identity** (Latin + Arabic that look like one brand).

**System: one superfamily — IBM Plex (OFL, self-hosted).**

```
  UI / headings   IBM Plex Sans
  Arabic pages    IBM Plex Sans Arabic   (same family → harmonized x-height/weight/rhythm)
  NUMBERS + code  IBM Plex Mono          (tabular by nature; doubles for curl snippets)
```

Coherence is free; "Egyptian" is carried by colour/voice, not type (= "subtly Egyptian, dev-first").
*(Cairo-for-Arabic was the considered warmer alternative; rejected to keep type as the calm
instrument. One config change away if Plex reads too corporate.)*

**Hero number = the brand's most important glyphs:** Plex Mono, `tabular-nums`, lining figures,
huge (`clamp(2.5rem, 8vw, 6rem)`); unit (`EGP`/`جنيه`) smaller + muted; the **delta** (▲/▼) carries
the only loud colour. Tabular figures are what stop the count-up + 5-min refresh from shifting
pixels — directly protecting the CLS budget.

**Scale:** ~1.25 modular; fluid `clamp()` (few media queries); body 16–18px.

**Bilingual/RTL mechanics:** CSS **logical properties** everywhere (`margin-inline`, `inset-inline`)
→ `dir="rtl"` flips the whole layout from one stylesheet. AR gets ~+1–2px + more line-height
(`:lang(ar)`). Test bidi number wrapping inside RTL text.

**Delivery (a CWV decision):** self-host (Google Fonts CDN is GDPR-illegal in DE), woff2,
**subset per script** (never ship Arabic glyphs to EN readers), 2–3 static weights (400/600[/700];
Mono 400/500), preload hero font, `font-display: swap`, **`size-adjust`/`ascent-override` on the
fallback @font-face** to make the swap cause **zero** CLS (automate with Fontaine). Net ~30–60 KB.

---

## 4. Brand & visual identity

Anchor: **the brand is "today" — a fresh number, honestly sourced, right now.** Both natural
colours are *semantically earned*.

**Palette: Nile base + Gold accent. Every colour means something.**

```
  Nile (deep blue-teal)  = Egypt (no pyramids) · trust · structure  → canvas, headings, links
  Gold (refined ochre)   = you LITERALLY track gold · value · "today's number"  → ONE accent (sparingly)
  Green ▲ / Red ▼         = up / down, SEMANTIC ONLY (never decorative) · colourblind-safe w/ arrows+sign
  Sand (warm neutral)    = light-mode paper (whisper of desert, not beige)
  Nile-ink               = dark-mode paper (gold glows on it)
```

**Tokens (direction; tune for AA/AAA when real):**

```
              LIGHT (default — mass/mobile/SEO-safe)     DARK (dev/data)
  --canvas    sand-50  #FBFAF6                            nile-900 #0B2A33
  --ink       #0E1A1F                                     paper    #F4F1EA
  --brand     nile-700 #114B5F                            nile-300 #7FB2BE
  --accent    gold-600 #B8860B                            gold-500 #D4A017
  --up/--down #1B873F / #C4452F                           (brightened on ink)
  --official  nile (solid)        --parallel  gold (dashed/range)   ← encodes §4 honesty invariants
```

Default: respect `prefers-color-scheme`, **first-paint light** (consumer mobile in daylight); dark
is a first-class toggle. `--official` solid vs `--parallel` gold-dashed makes the honesty-about-
uncertainty (§4, §2.10) *visible*; same trick for the two gold streams (world-derived vs retail).

**Mark: wordmark + the gold dot.**

```
        naharda•      the dot = today's data point = a period ("today.") = the sun/now
               ▔       favicon = the dot on Nile · it PULSES on each data refresh = the live heartbeat
```

One glyph is logo + favicon + animation principle: **the site feels alive because the data is
alive.** Latin `naharda•` leads (EN-first); `النهاردة` is the Arabic-context mark (the name is
already the dialect word — the most authentic asset).

**Voice:** Honest · Fresh · Precise + **Plain (EN) / Familiar (ar-EG dialect)**.
EN: "47.50. Updated 3 minutes ago. Here's where it came from." No hype, show the curl, state
uncertainty as a range. ar-EG: "الدولار النهاردة بكام؟ بـ 47.50 جنيه رسمي، وفي السوق الموازي بين 51 و 53."
A knowledgeable friend, not a sterile API. The honesty *is* the marketing.

**Texture: let the data be the decoration.** No pyramids/camels/hieroglyph clip-art. Visuals =
sparklines + the big number (also lightest → CWV). Optional low-opacity Nile-contour or Islamic-
geometric divider. One stroked line-icon set (Lucide-style), self-hosted SVG sprite.

---

## 5. Backend touchpoints (handlers added here; tables ship in `add-bootstrap`)

- `GET /v1/stats` — public aggregate metrics: `requests_served_total`, `uptime`,
  `data_points_count`, `sources_count`, `last_updated` per family, `signups_count`. Aggregate only,
  zero PII. Reads `usage_log`. **Add `/v1/stats: 300s` to the §9.2 cache table.**
- `POST /v1/signups` — email capture → `signups`. **Single opt-in** + explicit consent checkbox +
  privacy policy + deletion flow (§12). Honeypot field + the existing IP rate-limit. **No email sent
  in v1** (Resend is v2). Double opt-in deferred to v2.

## 6. Open / deferred

- Arabic URL slugs vs Latin — chose Latin (shareable; ranks on content). Revisitable.
- Dynamic OG image with live number → v1.1 (satori).
- `docs.naharda.com` → deferred; thin in-app docs page only in v1.
- Cairo-vs-Plex for Arabic — Plex chosen; one config swap if it reads too corporate.
