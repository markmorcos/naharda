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

### Requirement: Parallel FX SHALL aggregate at least two independent approved sources
When `SENSITIVE_SOURCES_ENABLED` is on, `parallel` MUST be computed from at least two
**independent** approved publishers (no source that merely re-serves another), so `n` ≥ 2 and the
`{min, avg, max}` spread reflects genuine cross-source disagreement. Each contributing source MUST
be named in `meta` with its fetch time.

#### Scenario: Aggregate from independent sources
- **WHEN** the flag is on and at least two parallel sources return a quote
- **THEN** `parallel` returns `{ min, avg, max, n, sources[] }` with `n` ≥ 2 and `meta` names each source

#### Scenario: Correlated source excluded
- **WHEN** a candidate source only re-publishes another registered source's value
- **THEN** it is NOT registered, so it cannot inflate `n` or skew the spread

### Requirement: Each parallel source SHALL be polite and fail-soft
Every parallel scraper MUST send the honest `User-Agent` + contact link, poll at low frequency,
and on a failed or out-of-band fetch contribute nothing rather than a bogus number — one source
failing degrades only `n`, never the response or `official` (§2.6, §9.5).

#### Scenario: One source fails
- **WHEN** a single parallel source errors or returns an out-of-band value (beyond its threshold)
- **THEN** that quote is excluded/held for review and `parallel` is served from the remaining sources

#### Scenario: All sources fail
- **WHEN** every parallel source fails to fetch
- **THEN** `official` still returns normally and `parallel` is empty/stale with a `meta` flag

### Requirement: Parallel sources SHALL live in the source registry with a tunable threshold
Each approved parallel source MUST be recorded in the `sources` registry (`family='fx'`,
`canonical=false`) with an `outlier_threshold`, and the parallel ingest MUST use that per-source
threshold (defaulting to 8% when absent) rather than a hardcoded value — so a source's tolerance is
tunable without a redeploy, symmetric with official sources.

#### Scenario: Per-source threshold drives the outlier guard
- **WHEN** a registered parallel source has an `outlier_threshold` set in the registry
- **THEN** the ingest holds/admits its quotes against that threshold, not a hardcoded constant

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

