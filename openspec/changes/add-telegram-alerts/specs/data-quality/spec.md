# data-quality

## ADDED Requirements

### Requirement: Alerts SHALL be self-contained and support Telegram
Data-quality alerts MUST render the relevant field values in the message body itself (not only a
summary), so any channel that shows plain text conveys the full detail. The service MUST support
delivery to Telegram (via a bot token + chat id) and to a generic webhook, sending to whichever are
configured, and MUST fall back to logs only when none are configured. Credentials MUST come from
secrets and never be logged.

#### Scenario: Held outlier alerts to Telegram
- **WHEN** an outlier is held and `TELEGRAM_BOT_TOKEN` + `TELEGRAM_CHAT_ID` are configured
- **THEN** a Telegram message is sent containing the alert summary and the field values (e.g. quote,
  value, trailing average)

#### Scenario: Plain-text channel shows full detail
- **WHEN** an alert is delivered to a channel that renders only a text field (e.g. Slack)
- **THEN** the message still includes the field values, not just the summary

#### Scenario: No channel configured
- **WHEN** neither Telegram nor a webhook is configured
- **THEN** the alert is logged (WARN) and no external request is made
