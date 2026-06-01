-- +goose Up
-- EGP gold prices, immutable history. Streams never merge (project.md §4).
CREATE TABLE gold_prices (
    id             BIGSERIAL PRIMARY KEY,
    stream         TEXT NOT NULL,                 -- world_derived | egypt_retail
    karat          INTEGER NOT NULL,              -- 18 | 21 | 24
    value_egp      DOUBLE PRECISION NOT NULL,     -- EGP per gram
    source         TEXT NOT NULL,
    fetched_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    pending_review BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX gold_prices_lookup_idx ON gold_prices (stream, karat, fetched_at DESC);

-- +goose Down
DROP TABLE gold_prices;
