-- +goose Up
-- Usage log, partitioned by day (90-day retention — project.md §9.4).
CREATE TABLE usage_log (
    id         BIGSERIAL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    endpoint   TEXT NOT NULL,
    ip_hash    TEXT,
    key_hash   TEXT,
    status     INTEGER,
    bytes      INTEGER,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- v1: a single DEFAULT partition holds all rows (the table is partition-ready
-- for schema-forwardness, but daily-partition creation and 90-day pruning are a
-- follow-up job, not yet implemented — see add-bootstrap design).
CREATE TABLE usage_log_default PARTITION OF usage_log DEFAULT;

CREATE INDEX usage_log_created_at_idx ON usage_log (created_at);

-- +goose Down
DROP TABLE usage_log;
