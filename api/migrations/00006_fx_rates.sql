-- +goose Up
-- EGP exchange rates, immutable history (project.md §2.4, §4).
CREATE TABLE fx_rates (
    id             BIGSERIAL PRIMARY KEY,
    market         TEXT NOT NULL,                 -- official | parallel
    base           TEXT NOT NULL DEFAULT 'EGP',
    quote          TEXT NOT NULL,                 -- USD, EUR, SAR, AED, KWD, GBP
    value          DOUBLE PRECISION NOT NULL,     -- EGP per 1 unit of quote
    source         TEXT NOT NULL,
    fetched_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    pending_review BOOLEAN NOT NULL DEFAULT FALSE  -- held outliers excluded from serving (§9.5)
);
CREATE INDEX fx_rates_lookup_idx ON fx_rates (market, quote, fetched_at DESC);

-- +goose Down
DROP TABLE fx_rates;
