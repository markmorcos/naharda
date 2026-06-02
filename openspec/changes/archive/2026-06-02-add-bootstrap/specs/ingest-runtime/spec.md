# ingest-runtime

The ingest/runtime substrate: run modes, provenance, manual override, and the data-quality guard.
Added by `add-bootstrap` (logic applied by later ingest changes).

## ADDED Requirements

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
