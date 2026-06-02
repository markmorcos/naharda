# add-sensitive-sources

> The 🟡 wedge: parallel-market FX (as an honest aggregate range) and Egypt-retail gold. ToS/
> scraping-sensitive — **gated behind a feature flag, dark until human source sign-off.**
> Cites §4 (Domain Model), §2 (source posture, honesty), §9.5, §12 (Legal), Open Question #1.

## Why

The single highest-value query — the parallel/black-market dollar — is the moat no one serves
cleanly (§1). This change builds it the *honest* way (a range with N sources, never a fake-precise
single number) and the *careful* way (low frequency, honest UA, removable on request).

## What changes

- **`/v1/fx`** — populates the `parallel` aggregate `{ min, avg, max, n, sources[] }` from 2–3
  approved sources (scraped via `goquery`). Never a single value (§4, Decision Log).
- **`/v1/gold`** — populates the `egypt_retail` stream (includes the "masna3eya" premium); kept
  **separate** from `world_derived`.
- **Feature flag** — `SENSITIVE_SOURCES_ENABLED` (default **false**); endpoints expose the fields
  empty until the flag is turned on *after* human sign-off on the exact sources.
- **Compliance scaffolding** — honest `User-Agent` + `abuse@naharda.com`, low-frequency polling,
  per-response attribution, remove-on-request posture.

## Scope

In: parallel-FX aggregation, retail-gold scraping, the feature flag + gating, the compliance
posture, graceful degradation.

## Non-goals

- Publishing any 🟡 value before human sign-off (the flag stays off — `TODO(ask)` on sources, §16 #1).
- Merging retail with world-derived gold (forbidden — §4).
- High-frequency / real-time scraping (§3 — "fresh enough", not a feed).

## Acceptance criteria

- [ ] With the flag **off**, `parallel`/`egypt_retail` fields are present but empty; nothing is scraped.
- [ ] With the flag **on**, `parallel` returns `{ min, avg, max, n, sources[] }` — never a single number.
- [ ] A failed scrape degrades only its field; the rest of the response is intact (fail-soft §2.6).
- [ ] Outbound requests carry an honest UA + contact link; polling is low-frequency.

## Dependencies

After `add-fx-official-and-gold-world` (the `fx`/`gold` shapes + data-quality guard). **Plus a
human gate**: explicit sign-off on the exact sources before the flag is enabled (§12, §16 #1).
