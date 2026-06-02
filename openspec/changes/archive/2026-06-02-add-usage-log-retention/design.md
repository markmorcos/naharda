# Design — add-usage-log-retention

## Maintenance job
A scheduled job (cron, `ingest`/`all` mode), runs daily:
```
  ensure partitions:   for d in [today .. today+7]:
                          CREATE TABLE IF NOT EXISTS usage_log_YYYYMMDD
                            PARTITION OF usage_log FOR VALUES FROM (d) TO (d+1)
  prune:               for each attached partition with upper bound < now()-90d:
                          DETACH + DROP
```
Idempotent (IF NOT EXISTS); cheap (DDL only). Registered like the other ingest jobs.

## One-time backfill
The existing `usage_log_default` holds all current rows. Migration: create dated partitions covering
the existing range and move rows out of DEFAULT into them (or, simplest: keep DEFAULT as a catch-all
and start dated partitions from deploy day — older rows age out of the 90-day window naturally). The
catch-all-keep approach is simplest and loses nothing; documented in the migration.

## Notes
- Inserts are unchanged (router by `created_at`); the writer already targets `usage_log`.
- `/v1/stats` count(*) spans partitions but recent windows hit only recent ones.

## Decisions
1. **Keep a DEFAULT catch-all** alongside dated partitions (belt-and-suspenders: a row with an
   unexpected timestamp never errors the insert), prune dated partitions >90d.
