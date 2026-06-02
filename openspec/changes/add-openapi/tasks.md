# Tasks — add-openapi

## Slice 1 — The spec
- [ ] Author the OpenAPI 3.1 doc: all public endpoints, the envelope + error components, examples.
- [ ] `GET /v1/openapi.json` handler (embedded doc; long Cache-Control + ETag).
- [ ] Drift test: every chi route is in the doc and vice-versa.

## Slice 2 — Interactive reference
- [ ] Render the spec on `/docs/reference` (Scalar or Stoplight Elements) with "try it".
- [ ] Link it from `/docs` and the homepage API section.
