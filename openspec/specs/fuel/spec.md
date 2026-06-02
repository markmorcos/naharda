# fuel Specification

## Purpose
TBD - created by archiving change add-safe-endpoints. Update Purpose after archive.
## Requirements
### Requirement: The API SHALL serve official pump prices, manually maintained
`GET /v1/fuel` MUST return prices for gasoline 80/92/95 and diesel, each with an `effective_from`
date and attribution to the Ministry of Petroleum/EGPC. Values are maintained via the manual
override mechanism (no scraper in v1) and cacheable for a day.

#### Scenario: Get fuel prices
- **WHEN** a client calls `GET /v1/fuel`
- **THEN** the response lists each product with its price, `effective_from`, EGPC attribution, and `max-age=86400`

#### Scenario: Announced price change
- **WHEN** the government announces new prices and an operator records them via manual override
- **THEN** subsequent responses reflect the new values with the new `effective_from`

