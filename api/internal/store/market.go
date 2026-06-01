package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// FXRate is one stored exchange-rate observation.
type FXRate struct {
	Market    string    `json:"market"`
	Quote     string    `json:"quote"`
	Value     float64   `json:"value"`
	Source    string    `json:"source"`
	FetchedAt time.Time `json:"fetched_at"`
}

// GoldPrice is one stored gold-price observation.
type GoldPrice struct {
	Stream    string    `json:"stream"`
	Karat     int       `json:"karat"`
	ValueEGP  float64   `json:"value_egp"`
	Source    string    `json:"source"`
	FetchedAt time.Time `json:"fetched_at"`
}

var errNoDB = errors.New("database not configured")

// InsertFXRate appends an immutable FX observation.
func (s *Store) InsertFXRate(ctx context.Context, market, quote string, value float64, source string, pending bool) error {
	if s == nil || s.Pool == nil {
		return errNoDB
	}
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO fx_rates (market, quote, value, source, pending_review) VALUES ($1,$2,$3,$4,$5)`,
		market, quote, value, source, pending)
	return err
}

// TrailingAvgFX returns the average and count of non-pending values within the window.
func (s *Store) TrailingAvgFX(ctx context.Context, market, quote string, within time.Duration) (float64, int, error) {
	if s == nil || s.Pool == nil {
		return 0, 0, errNoDB
	}
	var avg *float64
	var n int
	err := s.Pool.QueryRow(ctx,
		`SELECT avg(value), count(*) FROM fx_rates
		   WHERE market=$1 AND quote=$2 AND pending_review=false AND fetched_at >= $3`,
		market, quote, time.Now().Add(-within)).Scan(&avg, &n)
	if err != nil {
		return 0, 0, err
	}
	if avg == nil {
		return 0, 0, nil
	}
	return *avg, n, nil
}

// LatestFXRates returns the latest non-pending value per quote (last-good).
func (s *Store) LatestFXRates(ctx context.Context, market string) ([]FXRate, error) {
	if s == nil || s.Pool == nil {
		return nil, errNoDB
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT DISTINCT ON (quote) market, quote, value, source, fetched_at
		   FROM fx_rates WHERE market=$1 AND pending_review=false
		   ORDER BY quote, fetched_at DESC`, market)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFX(rows)
}

// LatestFXRate returns the latest non-pending value for a single quote.
func (s *Store) LatestFXRate(ctx context.Context, market, quote string) (float64, bool, error) {
	if s == nil || s.Pool == nil {
		return 0, false, errNoDB
	}
	var v float64
	err := s.Pool.QueryRow(ctx,
		`SELECT value FROM fx_rates WHERE market=$1 AND quote=$2 AND pending_review=false
		   ORDER BY fetched_at DESC LIMIT 1`, market, quote).Scan(&v)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return v, true, nil
}

// FXHistory returns immutable history for a quote, newest first.
func (s *Store) FXHistory(ctx context.Context, market, quote string, limit int) ([]FXRate, error) {
	if s == nil || s.Pool == nil {
		return nil, errNoDB
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT market, quote, value, source, fetched_at FROM fx_rates
		   WHERE market=$1 AND quote=$2 ORDER BY fetched_at DESC LIMIT $3`, market, quote, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFX(rows)
}

func scanFX(rows pgx.Rows) ([]FXRate, error) {
	out := []FXRate{}
	for rows.Next() {
		var r FXRate
		if err := rows.Scan(&r.Market, &r.Quote, &r.Value, &r.Source, &r.FetchedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// LatestParallelQuotes returns the latest non-pending parallel quote per source
// within the last day — the inputs to the {min,avg,max,n,sources} aggregate (§4).
func (s *Store) LatestParallelQuotes(ctx context.Context, quote string) ([]FXRate, error) {
	if s == nil || s.Pool == nil {
		return nil, errNoDB
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT DISTINCT ON (source) market, quote, value, source, fetched_at
		   FROM fx_rates
		   WHERE market='parallel' AND quote=$1 AND pending_review=false
		     AND fetched_at >= now() - interval '1 day'
		   ORDER BY source, fetched_at DESC`, quote)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFX(rows)
}

// LatestRetailGold returns the latest non-pending egypt_retail price per karat.
func (s *Store) LatestRetailGold(ctx context.Context) ([]GoldPrice, error) {
	return s.LatestGoldPrices(ctx, "egypt_retail")
}

// InsertGoldPrice appends an immutable gold observation.
func (s *Store) InsertGoldPrice(ctx context.Context, stream string, karat int, valueEGP float64, source string, pending bool) error {
	if s == nil || s.Pool == nil {
		return errNoDB
	}
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO gold_prices (stream, karat, value_egp, source, pending_review) VALUES ($1,$2,$3,$4,$5)`,
		stream, karat, valueEGP, source, pending)
	return err
}

// LatestGoldPrices returns the latest non-pending price per karat for a stream.
func (s *Store) LatestGoldPrices(ctx context.Context, stream string) ([]GoldPrice, error) {
	if s == nil || s.Pool == nil {
		return nil, errNoDB
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT DISTINCT ON (karat) stream, karat, value_egp, source, fetched_at
		   FROM gold_prices WHERE stream=$1 AND pending_review=false
		   ORDER BY karat, fetched_at DESC`, stream)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGold(rows)
}

// GoldHistory returns immutable history for a stream/karat, newest first.
func (s *Store) GoldHistory(ctx context.Context, stream string, karat, limit int) ([]GoldPrice, error) {
	if s == nil || s.Pool == nil {
		return nil, errNoDB
	}
	rows, err := s.Pool.Query(ctx,
		`SELECT stream, karat, value_egp, source, fetched_at FROM gold_prices
		   WHERE stream=$1 AND karat=$2 ORDER BY fetched_at DESC LIMIT $3`, stream, karat, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGold(rows)
}

func scanGold(rows pgx.Rows) ([]GoldPrice, error) {
	out := []GoldPrice{}
	for rows.Next() {
		var g GoldPrice
		if err := rows.Scan(&g.Stream, &g.Karat, &g.ValueEGP, &g.Source, &g.FetchedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}
