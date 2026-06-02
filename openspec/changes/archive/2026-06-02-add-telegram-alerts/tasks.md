# Tasks — add-telegram-alerts

Contained change in `internal/quality` + config. No schema/API change.

## Slice 1 — Rich messages + Telegram
- [x] `Alerter`: format `msg` + sorted `fields` into one human-readable `text`.
- [x] Add Telegram sender: `POST .../bot<token>/sendMessage` with `{chat_id, text}` when
      `TELEGRAM_BOT_TOKEN` + `TELEGRAM_CHAT_ID` are set.
- [x] Keep the generic `ALERT_WEBHOOK_URL` path (now sending the rich `text`); send to both if both set.
- [x] Best-effort + short timeout; failures logged, never block ingest; token never logged.

## Slice 2 — Config + wiring + ops
- [x] Config: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` envs; `main` passes them to `NewAlerter`.
- [x] `api/deployment.yaml`: add the two secret keys via `secretKeyRef` (optional: true).
- [x] Unit test: `format()` produces stable text incl. values; Alerter no-ops with nothing configured.
- [x] Document the operator setup (BotFather token + chat id → `naharda-api` secret).
