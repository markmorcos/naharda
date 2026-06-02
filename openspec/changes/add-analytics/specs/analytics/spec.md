# analytics

Cookieless usage analytics for the dashboard.

## ADDED Requirements

### Requirement: The dashboard SHALL collect cookieless usage analytics
The dashboard MUST record page views, unique visitors (without cookies), approximate geography, and
named custom events for key interactions, via a self-hosted analytics service. It MUST NOT set
cookies, MUST NOT require a consent banner, and MUST NOT store raw IP addresses.

#### Scenario: A page is viewed
- **WHEN** a visitor loads a dashboard page
- **THEN** a pageview is recorded (cookieless) with country-level geo and no raw IP stored

#### Scenario: A tracked button is clicked
- **WHEN** a visitor clicks a tracked control (e.g. subscribe, a details link, copy-curl)
- **THEN** the corresponding named event is recorded
