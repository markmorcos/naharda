# Design — add-telegram-alerts

Small, contained change in `internal/quality`. No schema/API change.

## Message formatting (fixes the lossy `fields`)

```
  format(msg, fields) → "msg — k1=v1 · k2=v2 …"   (keys sorted for stable output)
  e.g. "fx official outlier held — quote=USD · trailing_avg=52.22 · value=999.00"
```
This single `text` is what every channel receives, so Slack/Telegram/ntfy/generic all show the values.

## Delivery (send to whatever is configured)

```
  Alert(ctx, msg, fields):
    text = format(msg, fields)
    log.Warn(...)                                  # always
    if TELEGRAM_BOT_TOKEN && TELEGRAM_CHAT_ID:      # native Telegram
        POST https://api.telegram.org/bot<token>/sendMessage   {"chat_id": <id>, "text": text}
    if ALERT_WEBHOOK_URL:                           # generic (Slack incoming webhook, ntfy, …)
        POST <url>   {"text": text}
    # neither configured → log only (unchanged default)
```
Both can be on at once. Best-effort + short timeout; failures are logged, never block ingest.

## Config (env, from the `naharda-api` secret — §9.6)

| Env | Notes |
|---|---|
| `TELEGRAM_BOT_TOKEN` | from @BotFather; secret |
| `TELEGRAM_CHAT_ID` | your user id, or a negative group id |
| `ALERT_WEBHOOK_URL` | unchanged; optional generic webhook |

`Alerter` gains `tgToken`/`tgChatID` fields; `main` passes them from config. The token is never
logged.

## Decisions
1. **Native Telegram over the URL-query hack** — `{chat_id,text}` in the body is reliable;
   `?chat_id=…` + JSON body relies on Telegram merging query+body, which isn't guaranteed.
2. **One `text` string, not structured payloads** — keeps the Alerter channel-agnostic; richer
   per-channel formatting (Slack blocks, Markdown) is a future nicety, not v1.

## Operator inputs (what you provide — not in code)
- `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` added to the `naharda-api` k8s secret. The repo never
  contains them. (How-to in the apply summary.)
