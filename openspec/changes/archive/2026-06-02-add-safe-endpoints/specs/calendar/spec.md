# calendar

## ADDED Requirements

### Requirement: The API SHALL convert Hijri ↔ Gregorian locally
`GET /v1/calendar` MUST compute the Hijri/Gregorian mapping locally with no external dependency,
cacheable for a day.

#### Scenario: Get today's date mapping
- **WHEN** a client calls `GET /v1/calendar`
- **THEN** the response returns the Hijri and Gregorian dates computed without any network call
