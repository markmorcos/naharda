# ingest-runtime

## ADDED Requirements

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
