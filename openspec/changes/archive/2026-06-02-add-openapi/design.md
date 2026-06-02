# Design — add-openapi

## The spec
A hand-maintained OpenAPI 3.1 document (small, ~12 endpoints, stable). Served by the API:
```
  GET /v1/openapi.json   →  the OpenAPI doc   (Cache-Control: long; ETag)
```
Source it from an embedded `openapi.json` (or build it from a Go struct). Components: the `Envelope`
(`data` + `meta{tier,cached_at,freshness_seconds,sources[],attribution}`), the `Error` shape, and a
schema per resource. Document `Authorization: Bearer` as present-but-no-op in v1, and that
`/v1/stream` is SSE (text/event-stream).

## Drift guard
A Go test enumerates the chi routes and asserts each has a path in the OpenAPI doc (and vice-versa)
— so the spec can't silently drift from the API.

## Interactive reference
Render the spec on `/docs/reference` (or a tab on `/docs`) with **Scalar** or **Stoplight Elements**
(a single script/component, no heavy build). It loads `/v1/openapi.json` and offers "try it" against
the live API (CORS already allows it).

## Decisions
1. **Hand-maintained spec + a drift test**, not full code-gen — simplest for a small, stable API
   (§2.1), and the test keeps it honest.
2. Scalar/Elements over Swagger-UI (lighter, nicer, fewer deps).
