# status Specification

## Purpose
TBD - created by archiving change add-status-and-docs. Update Purpose after archive.
## Requirements
### Requirement: A public status page SHALL show live component health and freshness
A `/status` page MUST display live per-component health (API, database, and each data family) derived
from server-side checks, plus per-family freshness and counters from `/v1/stats`, and MUST link to
the UptimeKuma history for uptime/incidents. It MUST be reachable at `naharda.com/status` (`status.naharda.com` an optional alias).

#### Scenario: Everything healthy
- **WHEN** the API, database and data families are healthy and fresh
- **THEN** `/status` shows each component green, current freshness per family, and a link to UptimeKuma

#### Scenario: A component degraded
- **WHEN** a family is stale beyond its cache window or a check fails
- **THEN** that component is shown degraded/down while the rest remain green

### Requirement: Status checks SHALL be cached to protect the API
The status page's server-side checks MUST be cached briefly so a burst of status views does not
hammer the API; the page itself MUST be a cacheable read.

#### Scenario: Burst of status views
- **WHEN** many clients load `/status` within a short window
- **THEN** the underlying health checks are served from a short-lived cache rather than re-run per view

