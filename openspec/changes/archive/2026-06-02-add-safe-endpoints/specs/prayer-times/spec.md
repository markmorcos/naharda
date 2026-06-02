# prayer-times

## ADDED Requirements

### Requirement: The API SHALL serve daily prayer times per city with an explicit method
`GET /v1/prayer-times/{city}` MUST return the day's salah times with an explicit `method` code.
Today is cacheable for an hour; past dates are immutable and cacheable for a year.

#### Scenario: Get today's prayer times
- **WHEN** a client requests prayer times for a valid city today
- **THEN** the response includes all salah times and the `method` code, with `max-age=3600`

#### Scenario: Past dates are immutable
- **WHEN** a client requests prayer times for a past date
- **THEN** the values are unchanged from when first computed and cached for a year
