# add-bootstrap

> The foundation everything depends on: the single Go binary, Postgres + migrations, the full
> middleware chain, health/readiness, the deploy pipeline, and the schema-forward dormant hooks.
> Cites `project.md` §2 (Principles 1,2,3,7), §5 (Architecture), §6 (Tech Stack), §7 (Repo),
> §9 (Standards), §10 (Hooks), §11 (Deployment).

## Why

Per §8, v1 starts now and every endpoint rides shared infrastructure. This change builds the
**vertical skeleton** — one request can flow in, through every middleware, and back out with the
standard envelope — plus the deploy pipeline, so subsequent changes are pure feature slices.

## What changes

- **`api/` scaffold** — `cmd/server/main.go` selecting behavior by `MODE` (`api|ingest|all`; v1
  deploys `all`); env-driven `config`; `internal/{domain,store,http,scheduler}`.
- **Middleware chain** — request-id, slog logging, CORS, cache-headers, ETag/`304`, IP rate-limit
  (key-aware, no-op without a key), auth (reads `Authorization: Bearer`, no-op in v1).
- **Schema (goose)** — cross-cutting tables only: `api_keys`, `usage_log` (daily partitions, 90-day
  retention), `signups`, `manual_override`, `sources` (with per-source thresholds + `canonical`).
- **Health** — `/healthz` (process), `/readyz` (process + DB), wired to k8s probes.
- **Deploy** — `api/Dockerfile` (distroless/static, amd64+arm64), `api/deployment.yaml` (chart
  0.6.0), `.github/workflows/deploy-api.yaml`, and `deploy-naharda-api` dispatch in
  `markmorcos/infrastructure`.

## Scope

In: binary + modes, config, DB pool + migrations of cross-cutting tables, full middleware chain,
the response envelope + error format, ETag/304, health/readiness, scheduler skeleton, the
data-quality skeleton (`manual_override` + `sources` thresholds + guard interface + alert channel),
api deploy pipeline.

## Non-goals

- Any per-family data endpoints or ingesters (those are later changes).
- Active auth / key issuance / Stripe (v2 — schema-ready only here).
- Per-family data tables (created in their own changes, vertical-slice §2.8).

## Acceptance criteria

- [ ] A request to a stub endpoint returns the standard envelope with `ETag`; `If-None-Match` → `304`.
- [ ] `/healthz` and `/readyz` respond correctly; `/readyz` fails when DB is unreachable.
- [ ] IP rate-limit enforces 60/min + 1000/day; presenting a `Bearer` key is accepted but no-op.
- [ ] `goose up` creates all cross-cutting tables; `usage_log` is partitioned.
- [ ] `helm upgrade` via `deploy-naharda-api` rolls out the pod; probes pass.

## Dependencies

None — this is the first change. Unblocks all others.
