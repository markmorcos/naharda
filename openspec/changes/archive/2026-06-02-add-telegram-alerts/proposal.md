# add-telegram-alerts

> Make data-quality alerts **self-contained** (include the values, not just a summary) and add
> **native Telegram** delivery, while keeping the generic webhook path. Cites `project.md` §9.4
> (alerting), §9.5 (held-outlier alerts), §9.6 (secrets).

## Why

The `Alerter` currently POSTs `{ "text": "<summary>", "fields": {…values…} }`. Telegram and Slack
render only `text` and **drop `fields`** — so a held-outlier alert arrives as "fx official outlier
held" with no `quote` / `value` / `trailing_avg`, which is the part you actually need to act on.
And Telegram (the natural homelab pager) isn't directly supported — its API needs a `chat_id`, which
the generic `{text}` webhook can't supply cleanly.

## What changes

- **Rich, channel-agnostic messages:** fold `fields` into a single human-readable `text`
  (`"fx official outlier held — quote=USD · value=999.00 · trailing_avg=52.22"`). Fixes the
  detail-loss for *every* channel (Telegram, Slack, ntfy, generic).
- **Native Telegram:** when `TELEGRAM_BOT_TOKEN` + `TELEGRAM_CHAT_ID` are set, send a proper
  `{ chat_id, text }` to `sendMessage` — no reliance on the URL-query trick.
- **Both/either:** still honor `ALERT_WEBHOOK_URL`; if both Telegram and the webhook are configured,
  send to both; if neither, log only (unchanged default).

## Scope

In: `internal/quality` Alerter formatting + Telegram sender; two new env vars; the `naharda-api`
secret keys; wiring in `main`.

## Non-goals

- Two-way bot commands / interactive controls (this is outbound alerting only).
- Alert routing rules, severity levels, dedup/rate-limiting (could come later).
- Changing *what* triggers an alert (still the §9.5 held-outlier path).

## Acceptance criteria

- [ ] A held outlier produces a Telegram message containing the field values (quote/value/avg).
- [ ] Channels that render only `text` (Slack) now show the full detail.
- [ ] With no alert env set, behavior is unchanged (slog WARN only).
- [ ] Token/chat-id come from the secret via `secretKeyRef` (never logged — §9.6).

## Dependencies

Builds on the shipped `data-quality` capability (the held-outlier alert path). No schema/API change.
