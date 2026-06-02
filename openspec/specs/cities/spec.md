# cities Specification

## Purpose
TBD - created by archiving change add-safe-endpoints. Update Purpose after archive.
## Requirements
### Requirement: The API SHALL serve the canonical city list
`GET /v1/cities` MUST return the 13 canonical cities with their coordinates from a hardcoded source,
cacheable long-term.

#### Scenario: List cities
- **WHEN** a client calls `GET /v1/cities`
- **THEN** the response contains the 13 cities with lat/lon and a long `Cache-Control` max-age

