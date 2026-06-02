# Tasks — add-bootstrap

First change; no dependencies. Keep each slice a small PR.

## Slice 1 — Binary, config, modes
- [x] `go.mod`; `cmd/server/main.go` reading `MODE` (`api|ingest|all`, default `all`).
- [x] `internal/config` — env-only config (DB DSN, port, rate limits, `ALERT_WEBHOOK_URL`).
- [x] `internal/domain` — core types + `cities.go` (the 13 canonical cities, lat/lon).

## Slice 2 — Database & migrations
- [x] `pgxpool` store init; `goose` wired.
- [x] Migrations: `api_keys`, `usage_log` (partitioned, 90-day), `signups`, `manual_override`,
      `sources` (+ `canonical`, `outlier_threshold`, `disagreement_threshold`).

## Slice 3 — HTTP middleware chain
- [x] recover · request-id · slog-logging · CORS.
- [x] Response envelope + error helpers (English-only errors).
- [x] Cache-headers helper (per §9.2 policy); ETag + `If-None-Match` → `304`.
- [x] Rate-limit (IP token bucket 60/min + 1000/day; reads `Bearer`, no-op without key).
- [x] Auth middleware stub (parses `Bearer`, no-op in v1).

## Slice 4 — Health & runtime
- [x] `/healthz` (process) + `/readyz` (process + DB ping); wire to probes.
- [x] `scheduler` cron registration skeleton (no jobs yet).
- [x] Per-request `usage_log` write + structured logging.

## Slice 5 — Data-quality skeleton
- [x] `manual_override` precedence helper (override wins within `effective_from/to`).
- [x] Outlier-guard interface + `pending_review` write path + serve-last-good behavior.
- [x] Alerting: slog WARN + ntfy/webhook ping via `ALERT_WEBHOOK_URL`.

## Slice 6 — Deploy pipeline
- [x] `api/Dockerfile` (multi-stage → distroless/static, amd64+arm64, non-root, read-only FS).
- [x] `api/deployment.yaml` (chart 0.6.0; resources + probes + grace per §11).
- [x] `.github/workflows/deploy-api.yaml`; add `deploy-naharda-api` dispatch in `markmorcos/infrastructure`.
- [x] `api/README.md` with run + env docs.
