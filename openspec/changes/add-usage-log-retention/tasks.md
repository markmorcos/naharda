# Tasks — add-usage-log-retention

## Slice 1 — Partition maintenance
- [ ] Job: ensure dated partitions for [today .. today+7] (IF NOT EXISTS); keep a DEFAULT catch-all.
- [ ] Job: DETACH + DROP partitions whose upper bound is older than 90 days.
- [ ] Register on the scheduler (daily, ingest/all mode); idempotent.

## Slice 2 — Backfill + verify
- [ ] Migration note: dated partitions start at deploy; existing default rows age out within 90d.
- [ ] Verify: inserts land in the right dated partition; >90d partitions drop; `/v1/stats` correct.
