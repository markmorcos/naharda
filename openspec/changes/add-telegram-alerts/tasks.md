# Tasks — add-telegram-alerts

Contained change in `internal/quality` + config. No schema/API change.

## Slice 1 — Rich messages + Telegram
- [ ] `Alerter`: format `msg` + sorted `fields` into one human-readable `text`.
- [ ] Add Telegram sender: `POST .../bot<token>/sendMessage` with `{chat_id, text}` when
      `TELEGRAM_BOT_TOKEN` + `TELEGRAM_CHAT_ID` are set.
- [ ] Keep the generic `ALERT_WEBHOOK_URL` path (now sending the rich `text`); send to both if both set.
- [ ] Best-effort + short timeout; failures logged, never block ingest; token never logged.

## Slice 2 — Config + wiring + ops
- [ ] Config: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` envs; `main` passes them to `NewAlerter`.
- [ ] `api/deployment.yaml`: add the two secret keys via `secretKeyRef` (optional: true).
- [ ] Unit test: `format()` produces stable text incl. values; Alerter no-ops with nothing configured.
- [ ] Document the operator setup (BotFather token + chat id → `naharda-api` secret).
