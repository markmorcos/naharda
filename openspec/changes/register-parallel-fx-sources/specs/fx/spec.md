# fx

Concrete parallel-market sources for the (already-gated) `parallel` aggregate.
Added by `register-parallel-fx-sources` — closes the §16 #1 source sign-off.

## ADDED Requirements

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
- **WHEN** a single parallel source errors or returns an out-of-band value (>8% off trailing avg)
- **THEN** that quote is excluded/held for review and `parallel` is served from the remaining sources

#### Scenario: All sources fail
- **WHEN** every parallel source fails to fetch
- **THEN** `official` still returns normally and `parallel` is empty/stale with a `meta` flag
