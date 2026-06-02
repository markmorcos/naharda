# api-core Specification

## Purpose
TBD - created by archiving change add-bootstrap. Update Purpose after archive.
## Requirements
### Requirement: Responses SHALL use the standard envelope and error format
Every successful `/v1` response MUST use the envelope `{ data, meta }` where `meta` includes
`tier`, `cached_at`, `freshness_seconds`, `sources[]`, and `attribution`. Errors MUST use
`{ error: { code, message, retry_after_seconds } }` with English-only messages in v1.

#### Scenario: Successful response shape
- **WHEN** a client calls any `/v1` endpoint successfully
- **THEN** the body contains `data` and a `meta` object carrying `tier`, `sources[]`, and `attribution`

#### Scenario: Error response shape
- **WHEN** a request is rejected
- **THEN** the body is `{ error: { code, message, retry_after_seconds } }`

### Requirement: Every response SHALL carry an ETag and honor conditional requests
Handlers MUST set an `ETag` and a `Cache-Control` header per the project caching policy, and MUST
return `304 Not Modified` when `If-None-Match` matches.

#### Scenario: Conditional GET
- **WHEN** a client repeats a request with `If-None-Match` equal to the prior `ETag`
- **THEN** the server returns `304` with no body

### Requirement: The service SHALL expose health and readiness probes
`GET /healthz` MUST report process liveness; `GET /readyz` MUST report readiness including database
reachability.

#### Scenario: Database unreachable
- **WHEN** Postgres is unreachable and `/readyz` is called
- **THEN** the endpoint returns a non-200 status while `/healthz` still returns 200

### Requirement: Requests SHALL be rate-limited by IP and key-aware
The service MUST enforce 60 requests/minute and 1000/day per IP via an in-memory token bucket, and
MUST read an optional `Authorization: Bearer <key>` with no behavioral effect in v1.

#### Scenario: IP exceeds the limit
- **WHEN** an IP exceeds its minute or day budget
- **THEN** the server returns a `rate_limited` error with `retry_after_seconds`

#### Scenario: Bearer key present in v1
- **WHEN** a request includes `Authorization: Bearer <key>`
- **THEN** the request is processed normally with no key-specific quota applied

