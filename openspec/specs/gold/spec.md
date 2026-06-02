# gold Specification

## Purpose
TBD - created by archiving change add-fx-official-and-gold-world. Update Purpose after archive.
## Requirements
### Requirement: The API SHALL serve the world-derived gold price
`GET /v1/gold` MUST return a `world_derived` stream computed as `spot × FX × karat` for karats
18, 21, and 24, tagged distinctly from `egypt_retail`. The two streams MUST NOT be merged. Current
responses cache for 600s; history is immutable.

#### Scenario: Get current gold
- **WHEN** a client calls `GET /v1/gold`
- **THEN** the `world_derived` stream returns values for 18/21/24k with provenance, and `egypt_retail`
  is present (possibly empty) and never merged with `world_derived`

#### Scenario: Gold history is immutable
- **WHEN** a client requests a past gold value
- **THEN** the value is unchanged from first capture and cached for a year

### Requirement: Egypt-retail gold SHALL be a separate, flag-gated stream
When `SENSITIVE_SOURCES_ENABLED` is on, `GET /v1/gold` MUST populate an `egypt_retail` stream
(including the local masna3eya premium), kept distinct from `world_derived` and never merged. When
the flag is off, `egypt_retail` MUST be present but empty.

#### Scenario: Flag on
- **WHEN** the flag is enabled and a client calls `GET /v1/gold`
- **THEN** `egypt_retail` returns its own values separate from `world_derived`, never merged

#### Scenario: Flag off
- **WHEN** the flag is disabled
- **THEN** `egypt_retail` is present but empty and nothing is scraped

