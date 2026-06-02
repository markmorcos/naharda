# add-cbe-fx

> Wire the Central Bank of Egypt as the canonical `official` FX source, replacing the
> exchangerate-api stand-in. Cites `project.md` §4 (official = CBE-published), §2.5/§2.11
> (provenance, attribution), §9.5 (cross-source disagreement), §12 (CBE is 🟢/redistributable).

## Why

`add-fx-official-and-gold-world` shipped with a documented stand-in: the `official` market is fed by
exchangerate-api.com (a market reference rate), not the Central Bank of Egypt. §4 defines the
official rate as **CBE-published**. This makes it true.

## What changes

- **CBE source**: a `goquery` scraper of the CBE published rates (cbe.org.eg), feeding
  `market=official` with CBE as the named source/attribution.
- **Canonical + cross-check**: seed `sources` with CBE `canonical=true`; keep exchangerate-api as a
  secondary cross-check so the §9.5 disagreement guard can flag divergence (prefer CBE).
- Gold's world-derived calc keeps using the (now-CBE) USD/EGP automatically.

## Scope

In: the CBE scraper, the sources seed (CBE canonical), attribution update, disagreement wiring.

## Non-goals

- Parallel/black-market FX (that's `add-sensitive-sources`, human-gated).
- Changing the `/v1/fx` shape (same envelope; only the source/values change).

## Acceptance criteria

- [ ] `/v1/fx.official` is sourced from CBE with CBE attribution + `fetched_at`.
- [ ] CBE is the canonical source; on >2% divergence from the cross-check, CBE is served and `meta` flags it.
- [ ] Honest UA + contact on the CBE fetch (§9.6); fail-soft if CBE is unreachable (serve last-good).

## Dependencies

Builds on `fx` + `data-quality` (the canonical-source + disagreement machinery already exist).
