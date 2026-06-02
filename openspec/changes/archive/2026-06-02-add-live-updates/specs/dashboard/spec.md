# dashboard

## ADDED Requirements

### Requirement: The dashboard SHALL update live via progressive enhancement
The dashboard MUST subscribe to the live stream and update displayed numbers in place on a pushed
event, using the SSR-rendered value as the initial state. It MUST cause no layout shift, honor
`prefers-reduced-motion`, and remain fully correct if the stream is unavailable.

#### Scenario: Live update arrives
- **WHEN** an update event arrives for a displayed value
- **THEN** the number updates in place (no reload, no layout shift) and the gold dot pulses unless reduced-motion is set

#### Scenario: Stream unavailable
- **WHEN** the stream cannot connect
- **THEN** the page still shows the correct SSR values and the client retries with backoff
