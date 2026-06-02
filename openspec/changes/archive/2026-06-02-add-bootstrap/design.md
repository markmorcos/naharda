# Design — add-bootstrap

## Binary & modes (§5)
Single Go binary; `MODE` env selects behavior: `api` (HTTP, read-only), `ingest` (cron goroutines),
`all` (both — v1 default). Code is structured so `api` and `ingest` can split into separate
Deployments later **without code change**.

## Stack (§6, locked)
Go 1.23+ · `net/http` + `go-chi/chi/v5` · `jackc/pgx/v5/pgxpool` · `pressly/goose/v3` ·
`robfig/cron/v3` · `log/slog` (JSON). Config via env only; no config files in the image.

## Middleware chain (order matters)
```
  recover → request-id → slog-logging → cors → cache-headers → etag/304
          → rate-limit (IP; key-aware, no-op w/o key) → auth (reads Bearer, no-op) → handler
```
- **Envelope** (§9.1): `{ data, meta:{ tier, cached_at, freshness_seconds, sources[], attribution } }`.
- **Errors**: `{ error:{ code, message, retry_after_seconds } }`, English-only in v1.
- **ETag** on every response; honor `If-None-Match` → `304`.
- **Rate limit** (§9.3): in-memory token bucket, 60/min + 1000/day per IP; middleware reads optional
  `Authorization: Bearer <key>` from day one; per-key limits come from `api_keys` (dormant until v2).

## Schema decision (schema-forward §2.7, but vertical slices §2.8)
Bootstrap migrates **cross-cutting** tables only; **per-family data tables ship with their
endpoints** (so each later change is one end-to-end slice).

| Table | Purpose |
|---|---|
| `api_keys` | tier, quotas, `key_hash` (never plaintext), revoked — dormant until v2 |
| `usage_log` | daily partitions, 90-day retention; feeds `/v1/stats` + future billing |
| `signups` | email capture (consent, created_at) — endpoint added in `add-dashboard` |
| `manual_override` | family, key, value, `effective_from/to`, author — hand-set a value when an ingester breaks |
| `sources` | name, url, family, `canonical bool`, `outlier_threshold`, `disagreement_threshold` |

## Data-quality skeleton (§9.5 — resolves explore gap #3)
The *table* + *interfaces* land here; the *logic* is applied in `add-fx-official-and-gold-world`.
- **Held value** semantics: a value > its source's `outlier_threshold` off the trailing-1h average
  is written with a `pending_review` flag; the API keeps serving the **last-good** value and flags
  staleness in `meta` (fail-soft §2.6).
- **Thresholds**: per-source columns on `sources`, seeded with constitution defaults (5% outlier,
  2% FX-disagreement) — FX can run tighter than gold.
- **Canonical source**: a `canonical` flag per family; on disagreement, prefer it and flag in `meta`.
- **Alert channel (v1)**: `slog` WARN + a single ntfy/webhook ping (`ALERT_WEBHOOK_URL` env). No
  Slack/PagerDuty in v1.

## Health & observability (§9.4)
`/healthz` = process up; `/readyz` = process up + DB ping. slog JSON per request:
`endpoint, ip_hash, key_hash|null, status, bytes`. `usage_log` written from day one.

## Deploy (§11)
Multi-stage Dockerfile → `gcr.io/distroless/static-debian12` (amd64+arm64), non-root, read-only FS.
`api/deployment.yaml` (chart 0.6.0): requests `50m`/`64Mi`, limits `500m`/`256Mi`; startup+readiness
on `/readyz`, liveness on `/healthz`; `terminationGracePeriodSeconds: 15`.
`deploy-api.yaml` fires `repository_dispatch: deploy-naharda-api` → the infrastructure
`deploy-app` workflow (add the dispatch type there). Secrets via `secretKeyRef` only (§9.6).
