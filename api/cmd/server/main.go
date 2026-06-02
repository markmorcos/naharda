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
	fxingest "github.com/markmorcos/naharda/api/internal/ingest/fx"
	goldingest "github.com/markmorcos/naharda/api/internal/ingest/gold"
	"github.com/markmorcos/naharda/api/internal/ingest/sensitive"
	"github.com/markmorcos/naharda/api/internal/quality"
	"github.com/markmorcos/naharda/api/internal/scheduler"
	"github.com/markmorcos/naharda/api/internal/store"
	"github.com/markmorcos/naharda/api/internal/stream"
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

	// Live-update hub + a runtime context cancelled on shutdown (add-live-updates).
	hub := stream.NewHub(cfg.StreamMaxConns)
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()

	var srv *http.Server
	if runAPI {
		srv = &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           httpapi.NewRouter(cfg, st, logger, hub),
			ReadHeaderTimeout: 10 * time.Second,
		}
		go func() {
			logger.Info("http server listening", "addr", srv.Addr, "mode", cfg.Mode)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("http server error", "err", err)
				os.Exit(1)
			}
		}()
		// Bridge Postgres NOTIFY → SSE clients.
		go stream.Listen(runCtx, st, hub, logger)
	}

	var sch *scheduler.Scheduler
	if runIngest {
		sch = scheduler.New(logger)
		alerter := quality.NewAlerter(cfg.AlertWebhookURL, logger)
		runFX := func() { fxingest.Run(context.Background(), st, alerter, logger) }
		runGold := func() { goldingest.Run(context.Background(), st, alerter, logger) }
		if err := sch.Register("@every 1h", "fx-official", runFX); err != nil {
			logger.Error("register fx job", "err", err)
		}
		if err := sch.Register("@every 15m", "gold-world", runGold); err != nil {
			logger.Error("register gold job", "err", err)
		}
		// 🟡 sensitive ingest — only when the flag is on AND sources are registered (§8, §16 #1).
		var runParallel, runRetail func()
		if cfg.SensitiveEnabled {
			runParallel = func() { sensitive.ParallelFXRun(context.Background(), st, alerter, logger) }
			runRetail = func() { sensitive.RetailGoldRun(context.Background(), st, alerter, logger) }
			if err := sch.Register("@every 30m", "fx-parallel", runParallel); err != nil {
				logger.Error("register parallel job", "err", err)
			}
			if err := sch.Register("@every 30m", "gold-retail", runRetail); err != nil {
				logger.Error("register retail job", "err", err)
			}
			logger.Warn("sensitive sources ENABLED (🟡 parallel FX + retail gold)")
		} else {
			logger.Info("sensitive sources disabled (flag off)")
		}

		sch.Start()
		logger.Info("scheduler started", "mode", cfg.Mode)
		// Prime data on startup: FX first, then gold (which depends on USD/EGP).
		go func() {
			runFX()
			runGold()
			if runParallel != nil {
				runParallel()
				runRetail()
			}
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	runCancel()    // stop the stream listener
	hub.Shutdown() // drain SSE clients so they reconnect to the next pod
	if srv != nil {
		_ = srv.Shutdown(shutdownCtx)
	}
	if sch != nil {
		sch.Stop()
	}
	logger.Info("shutdown complete")
}
