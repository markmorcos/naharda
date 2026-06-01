-- +goose Up
-- Manual overrides: hand-set a value when an ingester is broken (§9.5).
CREATE TABLE manual_override (
    id             SERIAL PRIMARY KEY,
    family         TEXT NOT NULL,           -- fx | gold | fuel | ...
    key            TEXT NOT NULL,           -- e.g. usd, gold-21k, gasoline-92
    value          NUMERIC NOT NULL,
    effective_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    effective_to   TIMESTAMPTZ,             -- NULL = open-ended
    author         TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX manual_override_lookup_idx ON manual_override (family, key, effective_from);

-- +goose Down
DROP TABLE manual_override;
