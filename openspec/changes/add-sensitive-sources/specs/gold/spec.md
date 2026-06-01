# gold

Egypt-retail gold stream. Added by `add-sensitive-sources` (gated).

## ADDED Requirements

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
