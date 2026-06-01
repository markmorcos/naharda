package store

import (
	"context"
	"time"
)

// InsertSignup stores an email capture (idempotent on email). §10, §12.
func (s *Store) InsertSignup(ctx context.Context, email string, consent bool, locale string) error {
	if s == nil || s.Pool == nil {
		return errNoDB
	}
	var loc any
	if locale != "" {
		loc = locale
	}
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO signups (email, consent, locale) VALUES ($1,$2,$3)
		   ON CONFLICT (email) DO NOTHING`, email, consent, loc)
	return err
}

// Stats is the public aggregate exposed at /v1/stats (PII-free).
type Stats struct {
	RequestsServedTotal int64      `json:"requests_served_total"`
	SignupsCount        int64      `json:"signups_count"`
	FXDataPoints        int64      `json:"fx_data_points"`
	GoldDataPoints      int64      `json:"gold_data_points"`
	LastFXAt            *time.Time `json:"last_fx_at"`
	LastGoldAt          *time.Time `json:"last_gold_at"`
}

// GetStats returns aggregate counters (no personal data).
func (s *Store) GetStats(ctx context.Context) (Stats, error) {
	if s == nil || s.Pool == nil {
		return Stats{}, errNoDB
	}
	var st Stats
	err := s.Pool.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM usage_log),
			(SELECT count(*) FROM signups),
			(SELECT count(*) FROM fx_rates),
			(SELECT count(*) FROM gold_prices),
			(SELECT max(fetched_at) FROM fx_rates),
			(SELECT max(fetched_at) FROM gold_prices)
	`).Scan(&st.RequestsServedTotal, &st.SignupsCount, &st.FXDataPoints,
		&st.GoldDataPoints, &st.LastFXAt, &st.LastGoldAt)
	return st, err
}
