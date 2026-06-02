# fx

Parallel-market EGP rate. Added by `add-sensitive-sources` (gated).

## ADDED Requirements

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
