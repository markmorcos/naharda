// Command server is the single Naharda API binary. Behavior is selected by MODE:
// api (HTTP only), ingest (cron only), or all (both — the v1 default).
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/markmorcos/naharda/api/internal/config"
	httpapi "github.com/markmorcos/naharda/api/internal/http"
	"github.com/markmorcos/naharda/api/internal/scheduler"
	"github.com/markmorcos/naharda/api/internal/store"
	"github.com/markmorcos/naharda/api/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Migrations are idempotent; run them on boot when a database is configured.
	if cfg.DatabaseURL != "" {
		if err := store.Migrate(cfg.DatabaseURL, migrations.FS); err != nil {
			logger.Error("migrations failed", "err", err)
			os.Exit(1)
		}
		logger.Info("migrations applied")
	}

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("store init failed", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	runAPI := cfg.Mode == "api" || cfg.Mode == "all"
	runIngest := cfg.Mode == "ingest" || cfg.Mode == "all"

	var srv *http.Server
	if runAPI {
		srv = &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           httpapi.NewRouter(cfg, st, logger),
			ReadHeaderTimeout: 10 * time.Second,
		}
		go func() {
			logger.Info("http server listening", "addr", srv.Addr, "mode", cfg.Mode)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("http server error", "err", err)
				os.Exit(1)
			}
		}()
	}

	var sch *scheduler.Scheduler
	if runIngest {
		sch = scheduler.New(logger)
		// No jobs are registered in bootstrap; ingest changes add them here.
		sch.Start()
		logger.Info("scheduler started", "mode", cfg.Mode)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if srv != nil {
		_ = srv.Shutdown(shutdownCtx)
	}
	if sch != nil {
		sch.Stop()
	}
	logger.Info("shutdown complete")
}
