# gold

EGP gold prices. World-derived stream added by `add-fx-official-and-gold-world`.

## ADDED Requirements

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
