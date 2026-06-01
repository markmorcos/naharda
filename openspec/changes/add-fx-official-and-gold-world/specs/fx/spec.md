# fx

EGP exchange rates. Official market added by `add-fx-official-and-gold-world`.

## ADDED Requirements

### Requirement: The API SHALL serve the official EGP exchange rate
`GET /v1/fx` MUST return an `official` market populated from the Central Bank of Egypt with
provenance, and MUST always include the `parallel` aggregate shape `{ min, avg, max, n, sources[] }`
(empty until sensitive sources are enabled). Current responses cache for 300s; history is immutable.

#### Scenario: Get current FX
- **WHEN** a client calls `GET /v1/fx`
- **THEN** `official` is populated with `source` and `fetched_at`, and `parallel` is present (possibly empty)

#### Scenario: FX history is immutable
- **WHEN** a client requests a past FX value
- **THEN** the value is unchanged from first capture and cached for a year

### Requirement: Cross-source disagreement SHALL prefer the canonical source
When cross-checked FX sources differ beyond the disagreement threshold, the API MUST prefer the
source flagged `canonical` and flag the disagreement in `meta`.

#### Scenario: Sources disagree
- **WHEN** two FX sources differ by more than the threshold
- **THEN** the canonical source's value is served and `meta` records the disagreement
