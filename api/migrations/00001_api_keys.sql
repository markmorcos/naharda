-- +goose Up
-- API keys (dormant until v2; schema-forward per project.md §2.7, §10).
CREATE TABLE api_keys (
    id           SERIAL PRIMARY KEY,
    key_hash     TEXT NOT NULL UNIQUE,       -- never plaintext (§9.6)
    tier         TEXT NOT NULL DEFAULT 'free',
    minute_quota INTEGER,
    daily_quota  INTEGER,
    revoked      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE api_keys;
