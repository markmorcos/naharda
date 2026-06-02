# Project: Naharda

> **Naharda** (Egyptian Arabic النهاردة, "today" — romanized without the article, à la Talabat/Vezeeta).
> Domain: **naharda.com**. A single monorepo: **`markmorcos/naharda`**.
> This document is the **project constitution**: the durable decisions, principles, and
> standards that every feature spec must adhere to. It is intentionally stable. Things that
> change per feature (endpoint contracts, DB migrations, scrapers) live in **feature specs**,
> not here. See [§13 Spec-Driven Development](#13-spec-driven-development).

| Meta | Value |
|---|---|
| Status | Pre-implementation (planning complete) |
| Current phase | v1 — Public API + Dashboard |
| Repo | `markmorcos/naharda` (monorepo) |
| Domain | `naharda.com` (primary, `.com`) |
| Doc version | 1.1.0 |
| Last updated | 2026-06-01 |
| Owner | Mark Morcos |
| Jurisdiction | Berlin, Germany (entity + hosting) |

---

## 1. Mission & North Star

A reliable, fairly-priced, **dev-first data API** answering _"what's happening in Egypt right
now"_ — FX (official + parallel), gold, fuel, prayer times, weather, air quality, and calendar
— with a free public dashboard as the marketing front door and **B2B data licensing** as the
long-term revenue engine.

**Why this can win:**
- No incumbent does it well — existing options are fragile scrapers, paywalled terminals, or
  "free APIs" that vanish in months.
- Cache-friendly data on commodity hardware = near-zero marginal cost.
- Powers other people's apps → built-in network effect.
- _"What's the dollar today"_ is asked by ~100M Egyptians weekly; the **parallel-market FX**
  field alone is a wedge no one serves cleanly.
- The brand itself answers the question: **Naharda** = "today."

**One-line success test per phase:**
- **v1** — devs integrate the free tier and the dashboard ranks for Egyptian-currency queries.
- **v2** — a meaningful subset of free users convert to self-serve paid plans.
- **v3** — a handful of B2B contracts make this a real, low-cost, high-margin business.

---

## 2. Principles

These are the tie-breakers. When a spec decision is ambiguous, resolve in favor of these.

### Engineering
1. **Boring, durable tech.** Prefer stdlib and proven libraries over frameworks and magic.
2. **One binary, modes not microservices.** Split only when a real constraint forces it.
3. **Cache at the edge, compute rarely.** Every response is cacheable; the origin should mostly
   serve cache misses and ingest writes.
4. **Immutable history.** Past data is never mutated; corrections are new rows with provenance.
5. **Provenance always.** Every datum carries its source and fetch time. No anonymous numbers.
6. **Fail soft.** A broken source degrades one field, never the whole response.
7. **Schema-forward.** Build v2/v3 hooks into the schema from day one to avoid migrations later.
8. **Smaller working slice > bigger broken one.** Ship vertical slices end-to-end.

### Product
9. **Free tier is the marketing budget.** It exists to build dependency and an audience.
10. **Honesty about uncertainty.** Parallel rates are ranges with N sources, never a fake-precise
    single number.
11. **Attribution is non-negotiable** — both to respect upstream sources and to build trust.
12. **The dashboard sells the API.** It is a marketing surface first, a product second.

---

## 3. Non-Goals (project-wide)

The project will **not**, at any phase unless explicitly re-scoped:
- Become a general "world data" API — Egypt/MENA focus is the moat.
- Offer financial advice, trading signals, or predictions.
- Guarantee real-time (sub-minute) data — this is a "fresh enough" API, not a trading feed.
- Store personal data beyond what billing/accounts require.
- Incorporate or host in Egypt (jurisdictional risk — see [§12](#12-legal--compliance-posture)).

Phase-specific exclusions live in each phase section.

---

## 4. Domain Model

The durable data families. Concrete fields, units, and endpoint shapes are defined in feature
specs; this section fixes the **vocabulary and invariants**.

| Family | Description | Key invariants |
|---|---|---|
| **FX** | EGP exchange rates | Two markets: `official`, `parallel`. Parallel is always a `{min, avg, max, n, sources[]}` aggregate, never a single value. |
| **Gold** | EGP gold prices | Two streams: `world_derived` (computed from spot × FX × karat) and `egypt_retail` (scraped). Karats: 18, 21, 24. Streams diverge by local premium ("masna3eya"); never merge them. |
| **Fuel** | Pump prices | Products: gasoline 80/92/95, diesel. Changes are rare; tracked by `effective_from` date. |
| **Prayer times** | Daily salah times by city | Past dates immutable. Method is explicit (`method` code). |
| **Calendar** | Hijri ↔ Gregorian | Computed locally, no external dependency. Coptic calendar is a deferred add-on. |
| **Weather** | Current conditions by city | From a free, commercial-use-OK source. |
| **Air quality** | AQI by city | Same source family as weather. |

**Cities (v1 canonical list, hardcoded lat/lon):** Cairo, Giza, Alexandria, Hurghada,
Sharm El Sheikh, Aswan, Luxor, Mansoura, Tanta, Asyut, Port Said, Suez, Ismailia.

**Data-source posture** (sources themselves are pinned in feature specs):
- 🟢 Safe (free, commercial-use OK, stable): official FX cross-check, prayer times, weather, AQI,
  local calendar compute, world-derived gold.
- 🟡 Sensitive (ToS / scraping risk): **parallel FX**, Egypt-retail gold. These require explicit
  human sign-off on the exact sources before any scraper is written, and must degrade gracefully.

---

## 5. Architecture

```
[ Cloudflare ] → [ k3s Ingress ] → [ naharda-api pod (Go) ] → [ Postgres ]
                                          ↑
                            [ cron goroutines polling external sources ]
```

- **Single Go binary** (the `api/` service), behavior selected by `MODE` env: `api | ingest | all`.
  - `api` — HTTP handlers, read-only from Postgres, sets cache headers.
  - `ingest` — cron-scheduled goroutines (one per source) writing to Postgres.
  - `all` — both in one process. **v1 deploys `all`.**
- Designed so `api` and `ingest` can split into separate Deployments later **without code change**.
- Cloudflare absorbs read load via HTTP caching; origin sees cache misses + ingest writes only.
- The `web/` dashboard is a separate deployable in the same monorepo; it consumes the public API.

---

## 6. Tech Stack (locked)

| Layer | Choice | Notes |
|---|---|---|
| Language (api) | **Go 1.23+** | Light memory, fast cold start, single binary |
| HTTP | stdlib `net/http` + `github.com/go-chi/chi/v5` | No framework magic |
| DB | **Postgres** (existing in-cluster instance) | Reuse; do not deploy a new one |
| DB driver | `github.com/jackc/pgx/v5/pgxpool` | |
| Migrations | `github.com/pressly/goose/v3` | |
| Scheduler | `github.com/robfig/cron/v3` | |
| Logging | stdlib `log/slog` (JSON handler) | k8s-friendly |
| HTML scraping | `github.com/PuerkitoBio/goquery` | For 🟡 sources |
| Hijri calendar | `github.com/hablullah/go-hijri` | Local compute, no external call |
| Config | Environment variables only | No config files in the image |
| Container | Multi-stage build → `gcr.io/distroless/static-debian12` | amd64 + arm64 |
| Payments (v2) | **Stripe**, EUR-priced | EU jurisdiction; avoid Egyptian gateways |
| Dashboard (`web/`) | Astro or SvelteKit | Static-first, SEO-friendly |
| CDN/edge | Cloudflare (free tier for v1) | TLS, caching, DDoS, WAF |

---

## 7. Repository & Project Structure

**Single monorepo: `markmorcos/naharda`.** Mirrors the eventlane split (`backend/` + `frontend/`),
but with a shared `specs/` for spec-driven development and the constitution at the root.

```
naharda/                          # markmorcos/naharda (monorepo)
├── project.md                    # this document (the constitution)
├── README.md                     # monorepo overview + links to api/ and web/
├── docker-compose.yml            # local dev: api + postgres (+ web later)
├── openspec/                     # OpenSpec: changes + durable capability specs (see §13)
│   ├── changes/<name>/           # a proposed change (verb-led name, unnumbered)
│   │   ├── proposal.md           # what & why
│   │   ├── design.md             # how (architecture, files, migrations, sources)
│   │   └── tasks.md              # ordered, checkable tasks
│   ├── specs/<capability>/       # durable truth per capability (fx, gold, …), fed on archive
│   │   └── spec.md
│   └── config.yaml
├── api/                          # Go backend service
│   ├── cmd/server/main.go        # entrypoint; MODE=api|ingest|all
│   ├── internal/
│   │   ├── config/               # env-driven config
│   │   ├── domain/               # core types + cities.go
│   │   ├── store/                # pgx queries
│   │   ├── ingest/{fx,gold,fuel,prayer,weather}/
│   │   ├── http/
│   │   │   ├── handlers/         # one file per resource
│   │   │   ├── middleware/       # cache, ratelimit, cors, auth, logging
│   │   │   └── router.go
│   │   └── scheduler/            # cron registration
│   ├── migrations/               # goose SQL files
│   ├── Dockerfile
│   ├── deployment.yaml           # chart 0.6.0 values for the API service
│   ├── go.mod
│   └── README.md
└── web/                          # dashboard (Astro/SvelteKit) — added in a later phase
    ├── ...
    ├── Dockerfile
    └── deployment.yaml           # chart 0.6.0 values for the dashboard
```

**Monorepo conventions:**
- Each deployable subproject (`api/`, `web/`) owns its own `Dockerfile` and `deployment.yaml`.
- `openspec/` is repo-wide; a change touching only `api/` or only `web/` still lives in the
  shared `openspec/changes/` set (changes are global, not per-service).
- `project.md` (this file) is the single source of truth for both services.

---

## 8. Phase Roadmap

Each phase has an explicit **start trigger** and **definition of done**. Detailed scope lives in
the phase's feature specs.

### v1 — Public API + Dashboard
- **Trigger:** now.
- **Goal:** validate data quality, attract free users, build the audience.
- **In scope:** all seven data families; 13 cities; public read endpoints (no auth, IP
  rate-limited); aggressive HTTP caching; static dashboard (`web/`); email capture; public stats;
  `attribution` in every response.
- **Out of scope:** auth/keys/Stripe (schema-ready only), Arabic error i18n, webhooks/push,
  external historical backfill (collection starts day-of-deploy), Coptic calendar.
- **Done when:** safe-tier endpoints live end-to-end; official FX + world-derived gold ingest
  live; 🟡 sources scaffolded behind a feature flag pending human approval; dashboard live;
  email capture + `/v1/stats` working; README with integration examples; public status page.

### v2 — Self-Serve Monetization
- **Trigger (any one):** >1k dashboard DAU for 30 consecutive days; OR an unsolicited
  commercial-use inquiry; OR free-tier abuse (>100 rate-limit-block events/day).
- **Goal:** convert a subset of free users to paid.
- **In scope:** signup/login; API key issuance; per-key quotas (middleware already key-aware);
  Stripe Checkout + Customer Portal; updated ToS with **non-commercial restriction on free
  tier**; user usage dashboard; quota-warning emails (via Resend).
- **Done when:** Stripe end-to-end (Checkout + Portal + webhook tier sync); per-key quotas
  enforced; ToS published & counsel-reviewed; ≥3 paying customers within 60 days.

### v3 — B2B / Enterprise
- **Trigger (any one):** v2 reaches ~€1k/mo MRR; OR an inbound enterprise lead; OR a customer
  offers ≥10× Business tier for a custom SLA.
- **Goal:** land high-value contracts; treat freemium as lead-gen for these.
- **In scope:** custom contracts/SOWs; dedicated rate-limit pools; private endpoints / custom
  feeds; bulk historical export (CSV/Parquet via signed MinIO URLs); contractual SLA; Stripe
  Invoicing; infra hardening (EU VPS failover, Postgres replication, paid Cloudflare).
- **Done when:** 3 enterprise customers and ≥€5k/mo MRR sustained for 3 months.

---

## 9. Cross-Cutting Standards

These bind every endpoint and every spec.

### 9.1 API conventions
- All endpoints versioned under `/v1`. Breaking changes → `/v2`, never silent.
- JSON only. Responses follow the envelope:
  ```json
  {
    "data": { },
    "meta": {
      "tier": "free",
      "cached_at": "<RFC3339>",
      "freshness_seconds": 240,
      "sources": [ { "name": "...", "url": "...", "fetched_at": "<RFC3339>" } ],
      "attribution": "Data: <sources>. Free tier non-commercial."
    }
  }
  ```
- Errors follow:
  ```json
  { "error": { "code": "rate_limited", "message": "...", "retry_after_seconds": 30 } }
  ```
  Error messages are English-only in v1.
- Every response carries `ETag`; handlers honor `If-None-Match` → `304`.

### 9.2 Caching (origin sets `Cache-Control: public, max-age=N`)
| Data | max-age |
|---|---:|
| FX current | 300s |
| Gold current / weather / AQI | 600s |
| Fuel | 86400s |
| Calendar | 86400s |
| Prayer times — today | 3600s |
| Prayer times — past / FX & gold history | 31536000s (immutable) |
| `/v1/today` | 300s (strictest component) |

### 9.3 Rate limiting
- v1: IP-based, 60 req/min + 1000 req/day (in-memory token bucket).
- Middleware reads optional `Authorization: Bearer <key>` from day one (no-op until v2).
- Per-key limits come from the `api_keys` table (present in schema from v1).

### 9.4 Observability
- Structured JSON logs (`slog`); every request logs endpoint, ip_hash, key_hash|null, status,
  bytes.
- `usage_log` table (daily partitions, 90-day retention) feeds `/v1/stats` and future billing.
- Health: `/healthz` (process up) and `/readyz` (process up + DB reachable) — wired to k8s probes.
- Uptime monitored via existing UptimeKuma; public status page at `naharda.com/status` (`status.naharda.com` an optional Cloudflare alias).

### 9.5 Data quality
- Outlier guard: a new value >5% off the trailing 1h average is held + alerted, not auto-published.
- Source disagreement: if cross-checked sources differ beyond threshold (e.g. FX >2%), log,
  prefer the canonical source, and flag in `meta`.
- A manual-override table allows setting a value by hand when an ingester is broken.

### 9.6 Security
- Secrets via k8s `secretKeyRef` only; never in image or repo.
- API keys stored hashed (`key_hash`), never plaintext.
- Distroless runtime; non-root; read-only root FS where possible.
- Honest `User-Agent` + contact link (abuse@naharda.com) on all outbound scrapes.

---

## 10. Monetization hooks (built in v1, activated in v2)

To avoid refactors, v1 already includes:
- `api_keys` table (tier, quotas, revocation).
- Key-aware rate-limit middleware (no-op without a key).
- `usage_log` partitioned table.
- `meta.tier`, `meta.sources`, `meta.attribution` in every response.
- Email capture → `signups` table.
- Public `/v1/stats` endpoint.

Pricing (proposed, finalized in the v2 spec): Anonymous (free, IP-limited) → Free (signed up, 5k
rpd, non-commercial) → Hobby €9 → Startup €29 → Business €99 → Enterprise from €500. All EUR;
Stripe Tax handles VAT/MOSS. **The non-commercial restriction on free/Hobby tiers is the moat.**

---

## 11. Deployment & Infrastructure

- **k3s homelab:** `m720q` (amd64) + `pi` (arm64). Cloudflare in front.
- **Helm:** reuse the published `ghcr.io/markmorcos/infrastructure` chart **v0.6.0** (supports
  `resources`, `startupProbe`/`livenessProbe`/`readinessProbe`, `terminationGracePeriodSeconds`).
  **Do not author a bespoke chart.** Each subproject ships a `deployment.yaml` values file.
- **Namespace:** `naharda`. Deployment names follow the chart's `{namespace}-{name}` pattern →
  `naharda-api-deployment`, `naharda-web-deployment`.
- **Images:** `ghcr.io/markmorcos/naharda-api`, `ghcr.io/markmorcos/naharda-web`.
- **Hosts:** `api.naharda.com` (API), `naharda.com` (dashboard apex; docs at `/docs`, status at
  `/status`). `docs.naharda.com` / `status.naharda.com` are optional Cloudflare aliases, not separate deployables.
- **Pipeline:** `repository_dispatch` → the `deploy-app` workflow in `markmorcos/infrastructure`
  builds multi-arch and runs `helm upgrade --version <chart-ver> -f <subproject>/deployment.yaml`.
  Add dispatch types **`deploy-naharda-api`** and **`deploy-naharda-web`** to that workflow, with a
  per-service deploy workflow in this repo's `.github/workflows/` (mirror eventlane's
  `deploy-api.yaml` / `deploy-*-web.yaml`).
- **Resource profile (API starting point):** requests `50m`/`64Mi`, limits `500m`/`256Mi`; probes
  on `/readyz` (startup/readiness) and `/healthz` (liveness); `terminationGracePeriodSeconds: 15`.
- **Backups:** daily `pg_dump` → MinIO, 30-day retention (v1); offsite weekly (v2); PITR via
  managed Postgres (v3).
- **Reliability ladder:** best-effort (v1) → target 99.5% (v2) → contractual 99.5–99.9% with EU
  VPS failover + Postgres replication (v3).

---

## 12. Legal & Compliance Posture

- **Jurisdiction:** Berlin/Germany for entity and primary hosting. Do **not** incorporate or host
  in Egypt (historical state pressure on parallel-rate publication; mitigated jurisdictionally).
- **GDPR:** email collection requires a privacy policy + data-deletion flow (standard SaaS).
- **VAT/MOSS:** handled by Stripe Tax (enable at v2).
- **Source ToS:** 🟡 sources (parallel FX, retail gold) handled with low frequency, honest UA,
  contact link, per-response attribution, and willingness to remove on request. Exact sources
  require human sign-off per the relevant feature spec.
- **Commercial redistribution:** safe sources (CBE, Aladhan, Open-Meteo) permit it under
  attribution; negotiate per-source for enterprise redistribution where needed.

---

## 13. Spec-Driven Development

This file is the **constitution**. Concrete work happens as **OpenSpec changes** under
`openspec/changes/` (shared across the whole monorepo), which accumulate durable truth into
**capability specs** under `openspec/specs/`. Workflow:

1. **Create a change:** `openspec/changes/<verb-led-name>/` (e.g. `add-bootstrap` — named, not
   numbered; roadmap order lives in §8, not in the folder name).
2. **`proposal.md`** — the _what & why_: problem, user value, scope, non-goals, acceptance criteria.
   Must cite which `project.md` principles/standards it relies on; must not contradict them
   (if it must, amend `project.md` first via the Decision Log).
3. **`design.md`** — the _how_: architecture, files touched (note `api/` vs `web/`), data model
   changes, migrations, external sources (with ToS notes), risks.
4. **`tasks.md`** — ordered, individually checkable tasks; each maps to a small PR.
5. **Implement** against `tasks.md`; keep PRs vertical (one slice end-to-end).
6. **Archive** the change (`openspec archive`) to fold its requirements into the relevant
   `openspec/specs/<capability>/spec.md` — the durable, living truth.
7. **Amend `project.md`** only when a durable decision changes; record it in the Decision Log.

**Rules of thumb:**
- If it's stable across features → it belongs here. If it churns per feature → it belongs in a change/spec.
- Changes & specs reference this doc by section; this doc never references them (stays evergreen).
- Capability specs are the durable layer; changes are ephemeral and archived once shipped.
- Ambiguity inside a change is resolved by [§2 Principles](#2-principles); if Principles don't
  settle it, mark `TODO(ask)` and ask the human.

Suggested first changes (named, built in dependency order — see §8):
- `add-bootstrap` — `api/` scaffold, config, Postgres+migrations, all middleware, health/readiness,
  deploy pipeline, and the full v1 schema incl. dormant hooks (`api_keys`, `usage_log`, `signups`,
  `manual_override`).
- `add-safe-endpoints` — calendar, prayer-times, weather, AQI, cities, and **fuel** (🟢, manual-entry).
- `add-fx-official-and-gold-world` — official FX + world-derived gold + data-quality machinery.
- `add-sensitive-sources` — parallel FX + retail gold (🟡, gated on source approval, flag-off).
- `add-dashboard` — `web/` static-first dashboard (hybrid SSR) + email capture + `/v1/stats`.

---

## 14. Glossary

| Term | Meaning |
|---|---|
| **Naharda** | النهاردة, Egyptian Arabic for "today" — the brand. |
| **Official rate** | CBE-published EGP exchange rate. |
| **Parallel rate** | Black/grey-market EGP rate; published only as an aggregated range. |
| **Stream (gold)** | `world_derived` (computed) vs `egypt_retail` (scraped) gold pricing. |
| **masna3eya** | Egyptian gold retail premium/workmanship fee — why streams diverge. |
| **Safe / Sensitive source** | 🟢 free + commercial-OK + stable vs 🟡 ToS/scraping-risk. |
| **Tier** | Access level: anonymous, free, hobby, startup, business, enterprise. |
| **Mode** | Binary run mode: `api`, `ingest`, or `all`. |

---

## 15. Decision Log

Append-only. Each entry: date, decision, rationale, and what it supersedes.

| Date | Decision | Rationale |
|---|---|---|
| 2026-06-01 | Stack = Go + Postgres (not Spring/Kotlin) | Lighter footprint for a cache-heavy read API; deliberate divergence from the JVM apps. |
| 2026-06-01 | Reuse infrastructure Helm chart v0.6.0 | Probes/resources already supported; no bespoke chart. |
| 2026-06-01 | Parallel FX as `{min,avg,max,n,sources}` range | Honesty about uncertainty; avoids fake precision and single-source liability. |
| 2026-06-01 | Monetization hooks in schema from v1 | Avoids migrations when v2 activates billing. |
| 2026-06-01 | Free/Hobby tiers non-commercial | The licensing moat; forces commercial users to paid tiers. |
| 2026-06-01 | Berlin jurisdiction for entity + hosting | Avoid Egyptian regulatory risk around parallel-rate publication. |
| 2026-06-01 | Brand = **Naharda** (romanized w/o article) | `naharda.com` available; "today" is congruent with the product; follows Talabat/Vezeeta romanization norm. Supersedes working name "MENA Data API". |
| 2026-06-01 | Primary TLD = **`.com`** (`naharda.com`) | Long-term-safe and B2B-trusted; `.io` only as a redirect given Chagos/Mauritius uncertainty over that TLD. |
| 2026-06-01 | **Single monorepo `markmorcos/naharda`** (`api/` + `web/` + `specs/`) | Mirrors eventlane's monorepo split; one source of truth for the constitution and global spec sequence; per-service `deployment.yaml` + deploy workflow. |
| 2026-06-01 | Adopt **OpenSpec** (`openspec/changes/` + `openspec/specs/`), **unnumbered** verb-led change names | Use the configured tooling; gain a durable capability-spec layer §13 lacked; roadmap order lives in §8, not folder numbers. Supersedes the `specs/NNNN-slug/` convention previously in §7/§13. |
| 2026-06-02 | Docs + status served at `naharda.com/docs` and `/status` (paths), not dedicated subdomains | One `web/` deployable, no extra ingress/cert; SEO-clean canonical; `docs.`/`status.` subdomains optional Cloudflare aliases. Supersedes §9.4/§11 "reachable at the subdomain". |

---

## 16. Open Questions (resolve before / during early specs)

1. **Parallel FX sources** — shortlist 2–3, get explicit human sign-off (ToS-sensitive).
2. **Stripe entity** — GbR / UG / personal (deferrable to v2).
3. **Public-build vs stealth** — build-in-public on X/IndieHackers, or quiet until v2?
4. **Dashboard framework** — Astro vs SvelteKit for `web/` (decide in `0005-dashboard`).
