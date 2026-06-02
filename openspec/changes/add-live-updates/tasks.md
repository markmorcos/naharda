# Tasks — add-live-updates

Vertical slices. Transport assumed **SSE** per design.md (swap only the handler for WebSocket if chosen).

## Slice 0 — Decide transport (quick gate)
- [ ] Confirm SSE vs WebSocket vs poll-only (design.md "Open decisions"); record the call.

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
