# Tasks — add-usage-log-retention

## Slice 1 — Partition maintenance
- [x] Job: ensure dated partitions for [today .. today+7] (IF NOT EXISTS); keep a DEFAULT catch-all.
- [x] Job: DETACH + DROP partitions whose upper bound is older than 90 days.
- [x] Register on the scheduler (daily, ingest/all mode); idempotent.

## Slice 2 — Backfill + verify
- [x] Migration note: dated partitions start at deploy; existing default rows age out within 90d.
- [x] Verify: inserts land in the right dated partition; >90d partitions drop; `/v1/stats` correct.
