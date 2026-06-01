// Package store wraps the Postgres connection pool and migrations.
package store

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver for goose
	"github.com/pressly/goose/v3"
)

// Store holds the pgx connection pool.
type Store struct {
	Pool *pgxpool.Pool
}

// New creates a Store. An empty DSN yields a Store with no pool (readiness will
// fail, but the process still boots for health checks / local runs).
func New(ctx context.Context, dsn string) (*Store, error) {
	if dsn == "" {
		return &Store{}, nil
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Store{Pool: pool}, nil
}

// Ping verifies database reachability (used by /readyz).
func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.Pool == nil {
		return errors.New("database not configured")
	}
	return s.Pool.Ping(ctx)
}

// Close releases the pool.
func (s *Store) Close() {
	if s != nil && s.Pool != nil {
		s.Pool.Close()
	}
}

// LogUsage best-effort records a request to usage_log (§9.4). Errors are ignored.
func (s *Store) LogUsage(endpoint, ipHash, keyHash string, status, bytes int) {
	if s == nil || s.Pool == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var key any
	if keyHash != "" {
		key = keyHash
	}
	_, _ = s.Pool.Exec(ctx,
		`INSERT INTO usage_log (endpoint, ip_hash, key_hash, status, bytes) VALUES ($1,$2,$3,$4,$5)`,
		endpoint, ipHash, key, status, bytes)
}

// ActiveOverride returns a manually-set value for a family/key if one is in
// effect now (effective_from <= now < effective_to, or open-ended). This is the
// precedence hook used by ingest endpoints when a feed is broken (§9.5).
func (s *Store) ActiveOverride(ctx context.Context, family, key string) (float64, bool, error) {
	if s == nil || s.Pool == nil {
		return 0, false, errors.New("database not configured")
	}
	var value float64
	err := s.Pool.QueryRow(ctx,
		`SELECT value FROM manual_override
		   WHERE family = $1 AND key = $2
		     AND effective_from <= now()
		     AND (effective_to IS NULL OR effective_to > now())
		   ORDER BY effective_from DESC
		   LIMIT 1`, family, key).Scan(&value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return value, true, nil
}

// Migrate runs goose migrations from the embedded filesystem.
func Migrate(dsn string, fsys fs.FS) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	goose.SetBaseFS(fsys)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
