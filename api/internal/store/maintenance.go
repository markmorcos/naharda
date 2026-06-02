package store

import (
	"context"
	"fmt"
	"time"
)

// retentionDays is the §9.4 usage_log retention window.
const retentionDays = 90

// MaintainUsageLog ensures daily usage_log partitions exist for the next week
// and drops partitions whose date range is entirely older than the retention
// window. It is idempotent (CREATE ... IF NOT EXISTS) and cheap (DDL only).
//
// A DEFAULT catch-all partition (usage_log_default) is kept alongside the dated
// partitions: rows for days that don't yet have a dated partition land there and
// age out naturally, and a row with an unexpected timestamp never errors the
// insert. Because the default may already hold rows for the current day, creating
// today's partition can fail the overlap check; such per-day errors are tolerated
// (future days always succeed, so rows route correctly once their day arrives).
func (s *Store) MaintainUsageLog(ctx context.Context) error {
	if s == nil || s.Pool == nil {
		return nil
	}
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Ensure partitions for [today .. today+7].
	for i := 0; i <= 7; i++ {
		day := today.AddDate(0, 0, i)
		next := day.AddDate(0, 0, 1)
		name := "usage_log_" + day.Format("20060102")
		ddl := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s PARTITION OF usage_log FOR VALUES FROM ('%s') TO ('%s')",
			name, day.Format("2006-01-02"), next.Format("2006-01-02"),
		)
		if _, err := s.Pool.Exec(ctx, ddl); err != nil {
			// Tolerated: the DEFAULT partition may hold rows overlapping this range
			// (typically only the current day on first run).
			continue
		}
	}

	// Prune partitions whose entire range is older than the retention window.
	cutoff := today.AddDate(0, 0, -retentionDays)
	rows, err := s.Pool.Query(ctx, `
		SELECT c.relname
		FROM pg_inherits i
		JOIN pg_class c ON c.oid = i.inhrelid
		JOIN pg_class p ON p.oid = i.inhparent
		WHERE p.relname = 'usage_log'
		  AND c.relname ~ '^usage_log_[0-9]{8}$'`)
	if err != nil {
		return err
	}
	var stale []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			return err
		}
		day, perr := time.Parse("20060102", name[len("usage_log_"):])
		if perr != nil {
			continue
		}
		// The partition's upper bound is day+1; drop only if that bound is < cutoff.
		if day.AddDate(0, 0, 1).Before(cutoff) {
			stale = append(stale, name)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, name := range stale {
		if _, err := s.Pool.Exec(ctx, "ALTER TABLE usage_log DETACH PARTITION "+name); err != nil {
			return fmt.Errorf("detach %s: %w", name, err)
		}
		if _, err := s.Pool.Exec(ctx, "DROP TABLE "+name); err != nil {
			return fmt.Errorf("drop %s: %w", name, err)
		}
	}
	return nil
}
