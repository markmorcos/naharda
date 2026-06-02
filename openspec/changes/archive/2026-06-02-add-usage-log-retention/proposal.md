# add-usage-log-retention

> Make `usage_log` actually partition by day and prune to 90 days. Cites `project.md` §9.4
> (usage_log daily partitions, 90-day retention), §2.7 (schema-forward).

## Why

`usage_log` was created `PARTITION BY RANGE (created_at)` with only a DEFAULT partition — so every
row lands in one partition and **nothing is pruned**. §9.4 specifies daily partitions + 90-day
retention. Today the table grows unbounded and `/v1/stats`' `count(*)` slows over time. This wires
the maintenance the schema was built for.

## What changes

- **A maintenance job** (cron, runs in `ingest`/`all`): create the next few days' partitions ahead
  of time, and `DETACH`/`DROP` partitions older than 90 days.
- Migrate existing rows out of the DEFAULT partition into dated partitions (one-time).
- `/v1/stats` count stays fast (only recent partitions scanned for recent windows).

## Scope

In: the partition-management job + cron registration, the one-time backfill of the default partition.

## Non-goals

- Changing what `usage_log` stores or the `/v1/stats` shape.
- External log shipping / a data warehouse (out of scope).

## Acceptance criteria

- [ ] Daily partitions exist for today + the next N days; rows land in the correct dated partition.
- [ ] Partitions older than 90 days are dropped automatically.
- [ ] No data loss for in-window rows; `/v1/stats` still correct.

## Dependencies

Builds on `ingest-runtime` (the scheduler) + the `usage_log` schema (add-bootstrap).
