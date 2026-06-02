# weather

Current weather and air quality by city (Open-Meteo). Added by `add-safe-endpoints`.

## ADDED Requirements

### Requirement: The API SHALL serve current weather per city
`GET /v1/weather/{city}` MUST return current conditions sourced from Open-Meteo with attribution,
cacheable for 10 minutes.

#### Scenario: Get current weather
- **WHEN** a client requests weather for a valid city
- **THEN** the response returns current conditions, `meta.attribution` naming Open-Meteo, `max-age=600`

### Requirement: The API SHALL serve current air quality per city
`GET /v1/aqi/{city}` MUST return current AQI sourced from the Open-Meteo air-quality family,
cacheable for 10 minutes.

#### Scenario: Get current AQI
- **WHEN** a client requests AQI for a valid city
- **THEN** the response returns the current AQI with attribution and `max-age=600`
