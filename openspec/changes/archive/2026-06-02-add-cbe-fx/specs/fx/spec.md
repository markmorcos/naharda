# fx

## ADDED Requirements

### Requirement: The official FX rate SHALL be sourced from the Central Bank of Egypt
The `official` market MUST be populated from the Central Bank of Egypt (the §4 canonical source),
attributed to CBE with a fetch time. A secondary reference source MAY be used as a cross-check; on
divergence beyond the threshold the CBE value is served and the disagreement is flagged in `meta`.
If CBE is unreachable, the last-good value is served (fail-soft).

#### Scenario: CBE is the served official rate
- **WHEN** a client calls `GET /v1/fx`
- **THEN** `official` reflects the CBE-published rate with CBE attribution and `fetched_at`

#### Scenario: Cross-check disagreement
- **WHEN** the CBE value and the reference value differ beyond the threshold
- **THEN** the CBE (canonical) value is served and `meta` records the disagreement
