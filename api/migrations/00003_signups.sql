-- +goose Up
-- Email capture (§10). Endpoint added in add-dashboard.
CREATE TABLE signups (
    id         SERIAL PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    consent    BOOLEAN NOT NULL DEFAULT FALSE,  -- single opt-in (GDPR §12)
    locale     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE signups;
