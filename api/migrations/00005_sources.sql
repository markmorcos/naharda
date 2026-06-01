-- +goose Up
-- Source registry with per-source data-quality thresholds (§9.5).
CREATE TABLE sources (
    id                     SERIAL PRIMARY KEY,
    name                   TEXT NOT NULL,
    url                    TEXT,
    family                 TEXT NOT NULL,
    canonical              BOOLEAN NOT NULL DEFAULT FALSE,   -- preferred on disagreement
    outlier_threshold      NUMERIC NOT NULL DEFAULT 5.0,     -- percent off trailing-1h avg
    disagreement_threshold NUMERIC NOT NULL DEFAULT 2.0,     -- percent cross-source spread
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sources_family_idx ON sources (family);

-- +goose Down
DROP TABLE sources;
