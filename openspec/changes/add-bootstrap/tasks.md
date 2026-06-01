# Tasks — add-bootstrap

First change; no dependencies. Keep each slice a small PR.

## Slice 1 — Binary, config, modes
- [ ] `go.mod`; `cmd/server/main.go` reading `MODE` (`api|ingest|all`, default `all`).
- [ ] `internal/config` — env-only config (DB DSN, port, rate limits, `ALERT_WEBHOOK_URL`).
- [ ] `internal/domain` — core types + `cities.go` (the 13 canonical cities, lat/lon).

## Slice 2 — Database & migrations
- [ ] `pgxpool` store init; `goose` wired.
- [ ] Migrations: `api_keys`, `usage_log` (partitioned, 90-day), `signups`, `manual_override`,
      `sources` (+ `canonical`, `outlier_threshold`, `disagreement_threshold`).

## Slice 3 — HTTP middleware chain
- [ ] recover · request-id · slog-logging · CORS.
- [ ] Response envelope + error helpers (English-only errors).
- [ ] Cache-headers helper (per §9.2 policy); ETag + `If-None-Match` → `304`.
- [ ] Rate-limit (IP token bucket 60/min + 1000/day; reads `Bearer`, no-op without key).
- [ ] Auth middleware stub (parses `Bearer`, no-op in v1).

## Slice 4 — Health & runtime
- [ ] `/healthz` (process) + `/readyz` (process + DB ping); wire to probes.
- [ ] `scheduler` cron registration skeleton (no jobs yet).
- [ ] Per-request `usage_log` write + structured logging.

## Slice 5 — Data-quality skeleton
- [ ] `manual_override` precedence helper (override wins within `effective_from/to`).
- [ ] Outlier-guard interface + `pending_review` write path + serve-last-good behavior.
- [ ] Alerting: slog WARN + ntfy/webhook ping via `ALERT_WEBHOOK_URL`.

## Slice 6 — Deploy pipeline
- [ ] `api/Dockerfile` (multi-stage → distroless/static, amd64+arm64, non-root, read-only FS).
- [ ] `api/deployment.yaml` (chart 0.6.0; resources + probes + grace per §11).
- [ ] `.github/workflows/deploy-api.yaml`; add `deploy-naharda-api` dispatch in `markmorcos/infrastructure`.
- [ ] `api/README.md` with run + env docs.
