// Package config loads runtime configuration from environment variables only
// (no config files in the image — project.md §6).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the fully-resolved runtime configuration.
type Config struct {
	Mode              string   // api | ingest | all
	Port              string   // HTTP listen port
	DatabaseURL       string   // Postgres DSN
	AlertWebhookURL   string   // optional data-quality alert webhook (§9.5)
	RatePerMinute     int      // IP rate limit per minute (§9.3)
	RatePerDay        int      // IP rate limit per day (§9.3)
	CORSOrigins       []string // allowed CORS origins ("*" = all)
	SensitiveEnabled  bool     // 🟡 parallel FX / retail gold gate (§8, §16 #1) — default false
	StreamMaxConns    int      // max concurrent SSE connections per instance (add-live-updates)
	TelegramBotToken  string   // optional Telegram alert bot token (add-telegram-alerts)
	TelegramChatID    string   // optional Telegram alert chat id
	FXInterval        string   // FX ingest cron spec (add-fx-cadence)
	GoldInterval      string   // gold ingest cron spec
	SensitiveInterval string   // 🟡 sensitive ingest cron spec
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	c := Config{
		Mode:              getenv("MODE", "all"),
		Port:              getenv("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		AlertWebhookURL:   os.Getenv("ALERT_WEBHOOK_URL"),
		RatePerMinute:     getenvInt("RATE_PER_MINUTE", 60),
		RatePerDay:        getenvInt("RATE_PER_DAY", 1000),
		CORSOrigins:       splitCSV(getenv("CORS_ORIGINS", "*")),
		SensitiveEnabled:  getenvBool("SENSITIVE_SOURCES_ENABLED", false),
		StreamMaxConns:    getenvInt("STREAM_MAX_CONNS", 500),
		TelegramBotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:    os.Getenv("TELEGRAM_CHAT_ID"),
		FXInterval:        getenv("FX_INTERVAL", "@every 10m"),
		GoldInterval:      getenv("GOLD_INTERVAL", "@every 15m"),
		SensitiveInterval: getenv("SENSITIVE_INTERVAL", "@every 30m"),
	}
	switch c.Mode {
	case "api", "ingest", "all":
	default:
		return c, fmt.Errorf("invalid MODE %q (want api|ingest|all)", c.Mode)
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
