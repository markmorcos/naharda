# live-updates Specification

## Purpose
TBD - created by archiving change add-live-updates. Update Purpose after archive.
## Requirements
### Requirement: The API SHALL expose an additive live-update stream
The service MUST provide a streaming endpoint that, on connect, sends the current snapshot and
thereafter emits an event within seconds of each new datum. The stream MUST be `no-store` and MUST
NOT alter or degrade the cacheable `/v1/*` REST endpoints.

#### Scenario: Client connects and receives updates
- **WHEN** a client opens the stream
- **THEN** it first receives a snapshot event, then an update event within seconds of the next ingest write

#### Scenario: REST stays cacheable
- **WHEN** the stream feature is deployed
- **THEN** every `/v1/*` REST response keeps its existing ETag and Cache-Control (unchanged, edge-cacheable)

### Requirement: Streams SHALL be bounded, kept alive, and drained
The service MUST cap concurrent stream connections per instance, send periodic heartbeats so
intermediaries (~100s idle) don't drop the connection, and close streams gracefully on shutdown so
clients reconnect.

#### Scenario: Over the connection cap
- **WHEN** connections exceed the configured cap on an instance
- **THEN** further stream requests are rejected (503) while REST continues to serve

#### Scenario: Graceful shutdown
- **WHEN** the instance shuts down
- **THEN** open streams are closed cleanly and clients auto-reconnect to another instance

### Requirement: The stream SHALL be tier-gating-ready
The stream MUST pass through the key-aware middleware so it can later be restricted to paid tiers
(§10), without enforcing any restriction in v1.

#### Scenario: Anonymous access in v1
- **WHEN** an anonymous client connects in v1
- **THEN** the stream is served and the middleware records the (absent) key without blocking

