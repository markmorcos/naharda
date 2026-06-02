# add-fx-cadence

> Make ingest cadences configurable and bump FX so the data-quality guard (and alerts) actually
> work. Cites `project.md` §9.5 (held-outlier guard needs a baseline), §6 (env-only config).

## Why

The FX ingester runs `@every 1h`. The outlier guard needs **≥3 samples in the trailing hour** to
form a baseline — unreachable at 1h, so the FX guard (and its Telegram/webhook alerts) **never
fire**. During an EGP devaluation — the exact event worth catching — there's no protection or
alert. Bumping FX to ~10m makes the guard live, and making cadences env-driven keeps it tunable.

## What changes

- FX official cadence default `1h → 10m` (≈6 samples/hour → the guard activates).
- Ingest cadences become **env-configurable** (`FX_INTERVAL`, `GOLD_INTERVAL`, …) per §6.
- Confirm gold (`15m`) is fine; sensitive (`30m`) unchanged.
- No data-shape change; just more frequent ingest + a live guard.

## Scope

In: config for the intervals, wiring them into the scheduler, the FX default change.

## Non-goals

- Sub-minute polling (§3 — "fresh enough", not a trading feed).
- Changing the thresholds (that's data-quality tuning, separate).

## Acceptance criteria

- [ ] FX ingests ~every 10m by default; `FX_INTERVAL`/`GOLD_INTERVAL` override via env.
- [ ] With ≥3 in-window samples, an injected >5% FX move is held + alerted (the guard is live).
- [ ] Caching unchanged (FX still served `max-age=300`); origin load stays low (still cache-fronted).

## Dependencies

Builds on `ingest-runtime` (scheduler) + `data-quality` (the guard) + `add-telegram-alerts` (alerts).
