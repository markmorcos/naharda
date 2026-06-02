# docs

Public API documentation surface.

## ADDED Requirements

### Requirement: The API SHALL be documented with a public, crawlable docs surface
A `/docs` section MUST document every public endpoint with its path, parameters, a working `curl`
example and a sample JSON response, plus the shared conventions: the `{ data, meta }` envelope, the
error format, ETag/`304`, the Cache-Control policy, the rate limits, and attribution/free-tier terms.
It MUST be reachable at `docs.naharda.com` and carry proper title/description/canonical for discovery.

#### Scenario: A developer reads an endpoint reference
- **WHEN** a developer opens the docs for a resource (e.g. FX)
- **THEN** they see the path, params, a copyable `curl`, and a sample response, plus a link to the
  shared conventions (envelope, errors, rate limits, caching)

#### Scenario: The live stream is documented
- **WHEN** a developer looks up real-time updates
- **THEN** `GET /v1/stream` (SSE) is documented with an `EventSource` example

### Requirement: Docs SHALL not change the API contract
The docs surface MUST be read-only and MUST NOT alter any `/v1/*` endpoint.

#### Scenario: Docs deployed
- **WHEN** the docs surface ships
- **THEN** every `/v1/*` endpoint behaves exactly as before
