# public-api Specification

## Purpose
TBD - created by archiving change add-dashboard. Update Purpose after archive.
## Requirements
### Requirement: The API SHALL expose public aggregate stats
`GET /v1/stats` MUST return aggregate, PII-free metrics (requests served, uptime, data-point and
source counts, per-family last-updated timestamps, signups count) derived from `usage_log`, using
the standard response envelope and a `300s` cache policy.

#### Scenario: Fetch public stats
- **WHEN** any client calls `GET /v1/stats`
- **THEN** the response contains only aggregate, non-personal metrics in the standard envelope and
  carries `Cache-Control: public, max-age=300`

### Requirement: The API SHALL accept single-opt-in email capture
`POST /v1/signups` MUST store an email to the `signups` table under explicit consent (single
opt-in; no confirmation email in v1), protected by a honeypot field and the existing IP rate-limit.
It MUST NOT send any email in v1.

#### Scenario: Visitor submits their email with consent
- **WHEN** a visitor submits `POST /v1/signups` with a valid email and the consent flag set
- **THEN** the email is persisted to `signups` and no confirmation email is sent

#### Scenario: Bot trips the honeypot
- **WHEN** a request arrives with the honeypot field populated
- **THEN** the request is rejected and nothing is written

### Requirement: The API SHALL publish a machine-readable OpenAPI document
`GET /v1/openapi.json` MUST return a valid OpenAPI 3.1 document describing every public endpoint, the
standard `{data, meta}` envelope, the error shape, parameters, and example responses, and it MUST
stay in sync with the actual routes.

#### Scenario: Fetch the spec
- **WHEN** a developer calls `GET /v1/openapi.json`
- **THEN** they receive a valid OpenAPI 3.1 document covering all public endpoints and the envelope

#### Scenario: Spec matches the routes
- **WHEN** the route set changes
- **THEN** a test fails unless the OpenAPI document is updated to match (no silent drift)

