# naharda-api

The single Go binary behind `api.naharda.com`. See the project constitution at
[`../project.md`](../project.md) and the specs under `../openspec/`.

## Run modes

Behavior is selected by `MODE`:

| MODE | Behavior |
|---|---|
| `api` | HTTP handlers only (read-only) |
| `ingest` | Cron ingest jobs only |
| `all` | Both (the v1 default) |

## Configuration (env only)

| Var | Default | Notes |
|---|---|---|
| `MODE` | `all` | `api` \| `ingest` \| `all` |
| `PORT` | `8080` | HTTP listen port |
| `DATABASE_URL` | — | Postgres DSN (required for ingest/readiness) |
| `ALERT_WEBHOOK_URL` | — | Optional data-quality alert webhook (§9.5) |
| `RATE_PER_MINUTE` | `60` | IP rate limit / minute (§9.3) |
| `RATE_PER_DAY` | `1000` | IP rate limit / day (§9.3) |
| `CORS_ORIGINS` | `*` | Comma-separated allowed origins |

## Develop

```bash
go run ./cmd/server          # MODE=all on :8080
curl -s localhost:8080/healthz
curl -s localhost:8080/readyz
```

Migrations (goose, embedded under `migrations/`) run automatically on boot when
`DATABASE_URL` is set.

## Layout

```
cmd/server        entrypoint (MODE switch)
internal/config   env config
internal/domain   core types + the 13 canonical cities
internal/store    pgx pool + goose migrations + usage logging
internal/quality  data-quality skeleton (outlier guard + alerting)
internal/scheduler cron wrapper
internal/http     router + middleware + handlers + respond helpers
migrations        goose SQL (cross-cutting tables; per-family tables ship with their endpoints)
```
