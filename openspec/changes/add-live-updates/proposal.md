# add-live-updates

> Push real-time data updates (FX, gold, the `/v1/today` snapshot) to connected dashboard clients
> so the hero number and cards update **in place, without a page refresh** — while preserving the
> cache-first, stateless-read architecture. Cites `project.md` §2.3 (cache at edge, compute rarely),
> §5 (architecture; api/ingest split), §9 (standards), §10 (monetization).

## Why

The brand is "today / right now" (§1). The dashboard already renders fresh numbers per request
(edge-SSR), but a user sitting on the page sees stale values until they reload. A live push makes
the page feel alive the moment new data lands — the gold-dot heartbeat firing on a *real* update.

The catch, stated up front: the entire API is built on **stateless, edge-cacheable reads** (§2.3,
§5). A streaming connection is the opposite — stateful and uncacheable. This change must add the
live layer **without** compromising the cacheable REST endpoints, and should prefer the lightest
transport that does the job (see design.md — the recommendation is **SSE**, not WebSockets, despite
the original ask).

## What changes

- **A streaming endpoint** (e.g. `GET /v1/stream`) that, on connect, sends the current snapshot and
  thereafter pushes an event whenever ingest writes new data. Transport: SSE (recommended) or
  WebSocket — see the design.md decision.
- **Ingest emits a change signal** via Postgres `LISTEN/NOTIFY` on each new datum, so the API layer
  can broadcast without coupling to the ingest goroutines (preserves the §5 api/ingest split).
- **Dashboard live client** — progressive enhancement: the SSR-rendered value is the initial state;
  the stream updates numbers in place (tabular figures → no layout shift; gold-dot pulse on update;
  `prefers-reduced-motion` honored).

## Scope

In: the stream endpoint + broadcaster, the ingest NOTIFY hook, the dashboard client, connection
caps + heartbeat + graceful drain, and a tier-gating hook (no-op in v1, consistent with §10).

## Non-goals

- Client→server messages, presence, rooms (this is one-way server→client).
- Sub-second latency or replacing the REST endpoints (the stream is **additive**; `/v1/*` stays
  exactly as-is and cacheable).
- Multi-region fan-out / external message bus (Postgres NOTIFY is enough at this scale).
- Activating paid-tier enforcement (hook only; billing is the v2 track).

## Acceptance criteria

- [ ] A client connecting to the stream receives an initial snapshot, then an event within seconds
      of each ingest write.
- [ ] The dashboard updates the hero number + cards in place on a pushed event — no reload, no
      layout shift; pulse respects `prefers-reduced-motion`.
- [ ] All `/v1/*` REST endpoints remain unchanged and edge-cacheable (same ETag/Cache-Control).
- [ ] Connections are capped, heart-beated (survives Cloudflare's ~100s idle timeout), and drained
      cleanly on shutdown; a dropped connection auto-reconnects.
- [ ] Works behind nginx ingress + Cloudflare with no special config (SSE) — or documents the WS
      requirements if WebSocket is chosen.

## Dependencies

Builds on shipped capabilities: `api-core` (envelope/middleware), `ingest-runtime` (write path +
NOTIFY), `fx`/`gold` (data), `dashboard` (client). No new tables.
