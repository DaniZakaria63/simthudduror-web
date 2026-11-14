package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg, err := LoadFrom(func(string) string { return "" })
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	if cfg.Addr != ":8080" || cfg.DBPath != "./data/maulid.db" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.Environment != EnvironmentProduction || !cfg.CookieSecure {
		t.Fatalf("production should default to secure cookies: %+v", cfg)
	}
	if cfg.SessionTTL.Hours() != 168 {
		t.Fatalf("unexpected session ttl: %v", cfg.SessionTTL)
	}
}

func TestLoadFromValues(t *testing.T) {
	env := map[string]string{
		"MAULID_ADDR":              "127.0.0.1:9000",
		"MAULID_DB_PATH":           "/tmp/maulid.db",
		"MAULID_ENV":               "development",
		"MAULID_SESSION_TTL":       "24h",
		"MAULID_RATE_LIMIT_MAX":    "3",
		"MAULID_RATE_LIMIT_WINDOW": "2m",
		"MAULID_TRUSTED_ORIGINS":   "https://a.example, https://b.example/",
		"MAULID_COOKIE_SECURE":     "true",
	}
	cfg, err := LoadFrom(func(key string) string { return env[key] })
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	if cfg.Addr != "127.0.0.1:9000" || cfg.DBPath != "/tmp/maulid.db" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.SessionTTL.Hours() != 24 || cfg.RateLimitMax != 3 || cfg.RateLimitFor.Minutes() != 2 {
		t.Fatalf("unexpected durations: %+v", cfg)
	}
	if !cfg.CookieSecure {
		t.Fatal("explicit cookie secure override was ignored")
	}
	if len(cfg.TrustedOrigins) != 2 || cfg.TrustedOrigins[1] != "https://b.example" {
		t.Fatalf("unexpected trusted origins: %#v", cfg.TrustedOrigins)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	cases := map[string]map[string]string{
		"environment": {"MAULID_ENV": "staging"},
		"session ttl": {"MAULID_SESSION_TTL": "1m"},
		"rate max":    {"MAULID_RATE_LIMIT_MAX": "0"},
		"rate window": {"MAULID_RATE_LIMIT_WINDOW": "0s"},
		"cookie flag": {"MAULID_COOKIE_SECURE": "maybe"},
		"prod http":   {"MAULID_ENV": "production", "MAULID_COOKIE_SECURE": "false"},
	}
	for name, env := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := LoadFrom(func(key string) string { return env[key] }); err == nil {
				t.Fatalf("expected an error for %s", name)
			}
		})
	}
}
