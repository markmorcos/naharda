package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	for _, k := range []string{"MODE", "PORT", "RATE_PER_MINUTE", "RATE_PER_DAY", "SENSITIVE_SOURCES_ENABLED", "CORS_ORIGINS"} {
		t.Setenv(k, "")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Mode != "all" || c.Port != "8080" || c.RatePerMinute != 60 || c.RatePerDay != 1000 {
		t.Errorf("unexpected defaults: %+v", c)
	}
	if c.SensitiveEnabled {
		t.Error("SensitiveEnabled should default to false")
	}
}

func TestLoadInvalidMode(t *testing.T) {
	t.Setenv("MODE", "bogus")
	if _, err := Load(); err == nil {
		t.Error("expected error for invalid MODE")
	}
}

func TestSensitiveFlagOn(t *testing.T) {
	t.Setenv("SENSITIVE_SOURCES_ENABLED", "true")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !c.SensitiveEnabled {
		t.Error("expected SensitiveEnabled=true")
	}
}
