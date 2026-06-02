# ingest-runtime

## ADDED Requirements

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
