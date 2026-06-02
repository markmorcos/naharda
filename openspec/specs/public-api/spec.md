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

