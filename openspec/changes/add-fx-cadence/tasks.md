# Tasks — add-fx-cadence

## Slice 1 — Configurable cadences
- [ ] Config: `FX_INTERVAL` (default `@every 10m`), `GOLD_INTERVAL` (`@every 15m`),
      `SENSITIVE_INTERVAL` (`@every 30m`); invalid → log + default.
- [ ] `main` registers ingest jobs using the configured specs.

## Slice 2 — Verify the guard is live
- [ ] After ≥3 in-window samples, an injected >5% FX value is held `pending_review` + alerts fire.
- [ ] `/v1/fx` still `max-age=300`; origin load unchanged (cache-fronted).
