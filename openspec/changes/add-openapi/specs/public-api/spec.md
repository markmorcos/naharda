# public-api

## ADDED Requirements

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
