# data-quality Specification

## Purpose
TBD - created by archiving change add-fx-official-and-gold-world. Update Purpose after archive.
## Requirements
### Requirement: Outlier values SHALL be held and the last-good served
An ingested numeric value beyond its source's outlier threshold from the trailing-1h average MUST be
flagged `pending_review`, MUST trigger an alert, and MUST NOT be auto-published; the endpoint MUST
serve the last-good value and flag staleness in `meta`.

#### Scenario: Spike held
- **WHEN** an FX or gold value arrives more than the outlier threshold off the recent average
- **THEN** it is held `pending_review`, an alert fires, and the endpoint serves the last-good value
  with a `meta` staleness flag

### Requirement: Operators SHALL be able to override a broken ingester
An active `manual_override` for a family MUST take precedence over ingested values within its
effective window.

#### Scenario: Override a broken feed
- **WHEN** an operator records a manual override for a family whose ingester is failing
- **THEN** the endpoint serves the override value until its window ends

