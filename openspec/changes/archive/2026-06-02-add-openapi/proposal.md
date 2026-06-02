# add-openapi

> A machine-readable OpenAPI 3 spec + an interactive reference, for a dev-first API. Cites
> `project.md` §1 (dev-first data API), §9.1 (the envelope/error contract to describe).

## Why

The product is a **dev-first API** (§1). Hand-written docs (`/docs`) help humans, but developers
want a machine-readable contract to generate clients, validate, and explore. An OpenAPI spec is the
standard. Deferred from `add-status-and-docs`.

## What changes

- **`GET /v1/openapi.json`** — an OpenAPI 3.1 document describing every public endpoint, the
  `{data, meta}` envelope, the error shape, params, and example responses.
- **Interactive reference** — render it on `/docs` (or `/docs/reference`) with a lightweight viewer
  (Scalar or Stoplight Elements) — "try it" against the live API.
- Keep the hand-written `/docs` overview; the OpenAPI view is the reference companion.

## Scope

In: the openapi.json handler (served from the API, long-cached), the rendered reference page, links.

## Non-goals

- Code-first generation framework magic — a hand-maintained spec (or a small generator) is fine for
  ~12 stable endpoints (§2.1 boring/durable).
- Auth flows in the spec beyond noting the v2 `Authorization: Bearer` (no-op in v1).

## Acceptance criteria

- [ ] `GET /v1/openapi.json` returns a valid OpenAPI 3.1 doc covering all public endpoints + the envelope.
- [ ] An interactive reference renders it and can call the live API.
- [ ] The spec stays in sync with the endpoints (a test asserts every route is documented).

## Dependencies

Builds on `public-api` (the endpoints) + `docs` (the reference page).
