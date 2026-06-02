# Tasks — add-live-updates

Vertical slices. **Transport decided: SSE** (one-way push, plain HTTP, built-in reconnect, lightest
fit with §2.3/§5). WebSocket/poll-only rejected for the live feature (see design.md).

## Slice 0 — Decide transport (quick gate)
- [x] Transport chosen: **SSE** — WebSocket = full-duplex we'd never use; polling = the lighter
      free-tier fallback, kept in reserve, but the live feature ships SSE.

## Slice 1 — API stream + broadcaster
- [ ] In-process hub: register/unregister clients + broadcast; connection cap (`STREAM_MAX_CONNS`).
- [ ] `GET /v1/stream` SSE handler: initial snapshot (reuse the `/v1/today` data), `: heartbeat` ~25s, `no-store`.
- [ ] Graceful drain on shutdown; verify client auto-reconnect.

## Slice 2 — Ingest NOTIFY → broadcast
- [ ] Ingest emits `NOTIFY naharda_updates` (family payload) on each new datum.
- [ ] API holds one `LISTEN` connection (pgx `WaitForNotification`) → fan-out to the hub.

## Slice 3 — Dashboard live client
- [ ] `EventSource` island: update numbers in place; gold-dot pulse; honor `prefers-reduced-motion`.
- [ ] SSR value as initial state; fail-soft + backoff if the stream drops.

## Slice 4 — Gating + ops
- [ ] Route the stream through the key-aware middleware (tier-gating-ready, no-op in v1).
- [ ] Verify `/v1/*` REST unchanged (ETag/Cache-Control); load-check the conn cap; document
      Cloudflare/nginx notes (idle timeout, buffering off for SSE).
