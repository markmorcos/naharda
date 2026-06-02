# Design — add-live-updates

## The core tension (read this first)
The API is stateless + edge-cached (§2.3, §5): Cloudflare serves most reads, the origin computes
rarely, one light pod, replicas:1. A live connection inverts all of that — long-lived, stateful,
uncacheable, memory-resident, CDN-bypassing. So the design goal is to **add the live layer as a
thin, optional side-channel that never touches the cacheable REST path**, and pick the *lightest*
transport. Also: data changes at most every ~5 min (FX 300s, gold 600s), so "live" means "within
seconds of the ingest tick," not sub-second — which materially weakens the case for full-duplex
sockets.

## Transport decision (the ask was "sockets" — here's the honest comparison)

| | WebSocket | **SSE / EventSource** | Smart polling |
|---|---|---|---|
| direction | full-duplex | server→client (all we need) | client pull |
| protocol | `ws://` upgrade | plain HTTP (`text/event-stream`) | plain HTTP |
| nginx + Cloudflare | works, needs WS proxying + idle-timeout tuning | works as-is | works (edge-cached!) |
| reconnect | hand-rolled | **built-in** (auto + `Last-Event-ID`) | n/a |
| state / memory | per-conn | per-conn | none |
| fits §2.3 / §5 | weakest | good | best |
| effort | highest | medium | lowest |

**Recommendation: SSE.** We only need server→client push; SSE is one-way by design, plain HTTP
(proxies cleanly, no special Cloudflare config), has built-in reconnect + `Last-Event-ID`, and is
far lighter than a WS stack. WebSocket is full-duplex we'd never use. **Decision left to the human:**
if you specifically want WS (future bidirectional plans), the broadcaster below is
transport-agnostic — only the handler differs. And honestly, for 5-min-cadence data, **option C —
the dashboard polling the already-edge-cached `/v1/today` every ~30–60s — gives ~the same UX for
zero new infra**, and is worth considering as the actual v1.

## Mechanism — Postgres LISTEN/NOTIFY (keeps api/ingest decoupled, §5)

```
  ingest writes fx_rates/gold_prices  ──►  NOTIFY naharda_updates '{"family":"fx",...}'
                                                   │  (Postgres pub/sub)
  API instance holds ONE LISTEN conn ◄────────────┘
       └─ in-memory hub fans out to its connected stream clients
             └─ client updates the number in place
```

- Ingest `NOTIFY`s a small JSON payload (family + optionally the new value) on each write. Works
  even when api and ingest are **separate Deployments** (§5) — they share Postgres, not memory.
- Each API instance runs **one** dedicated `LISTEN` connection (pgx `conn.WaitForNotification`) and
  an in-process **hub** broadcasting to its clients. With replicas:1 this is trivial; with N
  replicas each LISTENs and serves its own clients — no sticky sessions needed for one-way SSE.
- Payload stays tiny; the client can use the pushed value directly, or treat it as a nudge and
  re-fetch `/v1/today` (edge-cached) — cheapest on the origin.

## Endpoint shape (SSE)

```
  GET /v1/stream            Accept: text/event-stream
    → event: snapshot       data: { ...current /v1/today... }     (immediately on connect)
    → event: update         data: { family:"fx", official:{usd:…}, at:"…" }   (on each change)
    → : heartbeat           (a comment every ~25s — keeps Cloudflare's ~100s idle from killing it)
  Cache-Control: no-store   (never cache a stream)
```

## Dashboard integration (progressive enhancement)
SSR renders the initial numbers (works with JS off / SEO intact). A tiny island opens
`new EventSource('/v1/stream')` and, on `update`, swaps the relevant number in place — tabular
figures → no layout shift, the gold dot pulses, `prefers-reduced-motion` skips the animation. If the
stream fails, the page is still correct (it had the SSR value); the client backs off and retries.
No markup/SEO impact.

## Cross-cutting
- **Connection cap** per instance (`STREAM_MAX_CONNS`) → reject over the cap (503) to protect the
  256Mi pod. **Heartbeat** every ~25s. **Graceful drain**: on shutdown send a final event + close so
  clients reconnect to the new pod.
- **Additive & non-breaking**: `/v1/*` REST is untouched and stays cacheable; the stream is a
  separate `no-store` route. This is the hard requirement — §2.3 must survive.
- **Tier-gating hook (§10)**: route the stream through the existing key-aware middleware (no-op in
  v1) so it *can* later be paid-tier-only (free = cached REST + polling, paid = live stream); v1
  leaves it open and simple.

## Schema
No new tables. Reuses `fx_rates`/`gold_prices` (read for snapshots) and Postgres NOTIFY (no storage).

## Decisions
1. **Transport: SSE** ✅ — one-way server→client push over plain HTTP, built-in reconnect, proxies
   cleanly through nginx+Cloudflare, far lighter than WS. WebSocket (full-duplex we'd never use) and
   poll-only (kept as the free-tier fallback) rejected for the live feature.

## Open decisions (for the human)
2. Push the value inline vs a "changed" nudge (client re-fetches the cached endpoint).
3. Gate behind a tier now, or leave open in v1.
