# data-quality Specification

## Purpose
TBD - created by archiving change add-fx-official-and-gold-world. Update Purpose after archive.
## Requirements
### Requirement: Outlier values SHALL be held and the last-good served
An ingested numeric value beyond its source's outlier threshold from the trailing-1h average MUST be
flagged `pending_review`, MUST trigger an alert, and MUST NOT be auto-published; the endpoint MUST
serve the last-good value and flag staleness in `meta`.

#### Scenario: Spike held
- **WHEN** an FX or gold value arrives more than the outlier threshold off the recent average
- **THEN** it is held `pending_review`, an alert fires, and the endpoint serves the last-good value
  with a `meta` staleness flag

### Requirement: Operators SHALL be able to override a broken ingester
An active `manual_override` for a family MUST take precedence over ingested values within its
effective window.

#### Scenario: Override a broken feed
- **WHEN** an operator records a manual override for a family whose ingester is failing
- **THEN** the endpoint serves the override value until its window ends

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

