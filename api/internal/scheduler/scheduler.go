// Package scheduler wraps robfig/cron for ingest jobs (project.md §5, §6).
package scheduler

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

// Scheduler runs registered cron jobs.
type Scheduler struct {
	c   *cron.Cron
	log *slog.Logger
}

// New creates a Scheduler.
func New(log *slog.Logger) *Scheduler {
	return &Scheduler{c: cron.New(), log: log}
}

// Register adds a named job on the given cron spec.
func (s *Scheduler) Register(spec, name string, job func()) error {
	_, err := s.c.AddFunc(spec, func() {
		s.log.Info("ingest job start", "job", name)
		job()
		s.log.Info("ingest job done", "job", name)
	})
	return err
}

// Start begins the cron loop.
func (s *Scheduler) Start() { s.c.Start() }

// Stop halts the cron loop and waits for running jobs to finish.
func (s *Scheduler) Stop() {
	ctx := s.c.Stop()
	<-ctx.Done()
}
