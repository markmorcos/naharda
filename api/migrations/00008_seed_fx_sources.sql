-- +goose Up
-- Seed the FX source registry (add-cbe-fx): the Central Bank of Egypt is the
-- canonical official source; exchangerate-api is kept as a non-canonical
-- cross-check / fallback. Idempotent on (name, family).
INSERT INTO sources (name, url, family, canonical)
VALUES
    ('Central Bank of Egypt', 'https://www.cbe.org.eg/en/economic-research/statistics/exchange-rates', 'fx', TRUE),
    ('exchangerate-api.com', 'https://www.exchangerate-api.com', 'fx', FALSE)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM sources WHERE family = 'fx' AND name IN ('Central Bank of Egypt', 'exchangerate-api.com');
