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

-- Default partition; a maintenance job creates daily partitions and prunes >90d.
CREATE TABLE usage_log_default PARTITION OF usage_log DEFAULT;

CREATE INDEX usage_log_created_at_idx ON usage_log (created_at);

-- +goose Down
DROP TABLE usage_log;
