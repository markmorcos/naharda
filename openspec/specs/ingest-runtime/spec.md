# ingest-runtime Specification

## Purpose
TBD - created by archiving change add-bootstrap. Update Purpose after archive.
## Requirements
### Requirement: Behavior SHALL be selected by MODE
The binary MUST run as `api`, `ingest`, or `all` based on the `MODE` env var, defaulting to `all`,
structured so `api` and `ingest` can split into separate Deployments without code change.

#### Scenario: API-only mode
- **WHEN** the binary starts with `MODE=api`
- **THEN** HTTP handlers run and no ingest cron jobs are scheduled

### Requirement: Every datum SHALL carry provenance
Stored values MUST record their source and fetch time, and past data MUST be immutable (corrections
are new rows).

#### Scenario: Ingested value stored
- **WHEN** an ingester writes a value
- **THEN** the row records `source` and `fetched_at`, and prior rows are never mutated

### Requirement: Manual overrides SHALL take precedence within their window
A value present in `manual_override` for a family/key MUST be served in preference to ingested
values while within its `effective_from`/`effective_to` window.

#### Scenario: Active override
- **WHEN** an active override exists for a family and the endpoint is queried
- **THEN** the override value is served instead of the ingested value

### Requirement: Outliers SHALL be held, not auto-published
A new value beyond its source's `outlier_threshold` from the trailing-1h average MUST be written
with a `pending_review` flag and alerted; the API MUST keep serving the last-good value and flag
staleness in `meta`.

#### Scenario: Suspicious spike
- **WHEN** an ingested value exceeds the outlier threshold
- **THEN** it is flagged `pending_review`, an alert fires, and the endpoint still serves the
  last-good value with a `meta` staleness flag

### Requirement: Ingest SHALL emit a change notification on new data
When an ingester writes a new (non-pending) datum, it MUST emit a Postgres `NOTIFY` carrying the
affected family, so API instances can broadcast it without coupling ingest to the API process
(preserving the api/ingest split §5).

#### Scenario: New FX value ingested
- **WHEN** the FX ingester stores a new official rate
- **THEN** a NOTIFY is emitted on the updates channel with the family identifier

#### Scenario: Separate api/ingest deployments
- **WHEN** ingest and api run as separate Deployments
- **THEN** the API still receives the notification via Postgres, not shared memory

### Requirement: Ingest cadences SHALL be configurable and frequent enough for the guard
Ingest intervals MUST be configurable via environment variables, and the FX interval MUST default
to a frequency that lets the outlier guard accumulate a baseline (≥3 samples/hour), so held-outlier
detection and alerting are functional.

#### Scenario: FX guard has a baseline
- **WHEN** the API has been ingesting FX at the default cadence for the baseline window
- **THEN** a sufficiently-deviating new FX value is held as pending_review and an alert is emitted

#### Scenario: Cadence override
- **WHEN** `FX_INTERVAL` is set to a valid cron spec
- **THEN** the FX ingester uses it; an invalid value falls back to the default without crashing

### Requirement: usage_log SHALL be partitioned by day and pruned to 90 days
A scheduled maintenance job MUST ensure daily `usage_log` partitions exist ahead of time and MUST
drop partitions older than 90 days, enforcing the §9.4 retention without losing in-window rows.

#### Scenario: Partitions are maintained
- **WHEN** the maintenance job runs
- **THEN** today's and the next several days' partitions exist and rows land in the correct dated partition

#### Scenario: Old data is pruned
- **WHEN** a partition's date range is entirely older than 90 days
- **THEN** that partition is detached and dropped, while in-window data is retained

