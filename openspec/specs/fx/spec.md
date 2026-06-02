# fx Specification

## Purpose
TBD - created by archiving change add-fx-official-and-gold-world. Update Purpose after archive.
## Requirements
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

### Requirement: Parallel FX SHALL be published only as an aggregate range, behind a flag
When `SENSITIVE_SOURCES_ENABLED` is on, `GET /v1/fx` MUST populate `parallel` as
`{ min, avg, max, n, sources[] }` aggregated from multiple approved sources — never a single value.
When the flag is off, `parallel` MUST be present but empty and nothing is scraped.

#### Scenario: Flag off
- **WHEN** the sensitive-sources flag is disabled and a client calls `GET /v1/fx`
- **THEN** `parallel` is present but empty and no scraping occurs

#### Scenario: Flag on
- **WHEN** the flag is enabled (after source sign-off)
- **THEN** `parallel` returns `{ min, avg, max, n, sources[] }` with `n` ≥ 2 and never a single number

### Requirement: A failed parallel scrape SHALL degrade only its field
A failure fetching a 🟡 source MUST NOT break the response; `parallel` degrades (empty/stale with a
`meta` flag) while `official` and the rest remain intact.

#### Scenario: Source unavailable
- **WHEN** a parallel source fails to fetch
- **THEN** `official` still returns normally and `parallel` carries a `meta` staleness flag

