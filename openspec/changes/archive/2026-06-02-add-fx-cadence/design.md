# Design — add-fx-cadence

## Config
```
  FX_INTERVAL    default "@every 10m"   (was 1h)
  GOLD_INTERVAL  default "@every 15m"
  SENSITIVE_INTERVAL default "@every 30m"
```
Env strings parse as robfig/cron specs; `main` passes them to `sch.Register`. Invalid spec → log +
fall back to the default (don't crash).

## Why 10m
The guard requires `n >= 3` non-stale samples in the trailing hour. At 10m that's ~6 samples — the
baseline forms within ~30m of boot, so a genuine >5% move is caught and alerted. CBE itself moves
slowly, so most 10m polls return the same value (cheap; cache window 300s unaffected) — the point is
the *baseline*, not freshness.

## Interplay
- Pairs with `add-cbe-fx` (CBE is the source; cadence is how often we poll it).
- Pairs with `add-telegram-alerts` (now FX holds will actually page you).
- Origin load stays trivial — `/v1/fx` is still `max-age=300` edge-cached; this only changes how
  often the *ingester* writes, not how often clients hit origin.

## Decisions
1. **Env-driven cadences** (not hardcoded) so this is tunable without a code change, per §6.
