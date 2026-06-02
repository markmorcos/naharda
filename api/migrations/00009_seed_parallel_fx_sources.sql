-- +goose Up
-- Seed the parallel-market FX sources (register-parallel-fx-sources, §16 #1).
-- These are non-canonical 🟡 sources scraped at low frequency; the 8% outlier
-- threshold is wider than official's 5% because the parallel market is noisier
-- (§9.5). The ingest reads outlier_threshold per source, so these are tunable
-- without a redeploy.
INSERT INTO sources (name, url, family, canonical, outlier_threshold)
VALUES
    ('egcurrency.com',      'https://egcurrency.com/en/currency/usd-to-egp/blackmarket', 'fx', FALSE, 8.0),
    ('blackmarketlive.org', 'https://en.blackmarketlive.org/egp/usd/',                   'fx', FALSE, 8.0),
    ('sarfegp.com',         'https://sarfegp.com/en/us-dollar-to-egp-black-market/',     'fx', FALSE, 8.0)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM sources WHERE family = 'fx'
  AND name IN ('egcurrency.com', 'blackmarketlive.org', 'sarfegp.com');
