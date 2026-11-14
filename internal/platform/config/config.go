// Package config loads the runtime configuration for the maulid.id monolith.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// EnvironmentDevelopment relaxes the Secure cookie flag so the app can
	// be used over plain HTTP on localhost.
	EnvironmentDevelopment = "development"
	// EnvironmentProduction is the default; cookies are Secure and HSTS is sent.
	EnvironmentProduction = "production"
)

// Config is the complete configuration read from the process environment.
type Config struct {
	Addr           string
	DBPath         string
	Environment    string
	SessionTTL     time.Duration
	CookieSecure   bool
	TrustedOrigins []string
	RateLimitMax   int
	RateLimitFor   time.Duration
}

// Load reads configuration from the environment. It fails rather than silently
// falling back when a supplied value is invalid.
func Load() (Config, error) {
	return LoadFrom(os.Getenv)
}

// LoadFrom supports tests without mutating the process environment.
func LoadFrom(getenv func(string) string) (Config, error) {
	cfg := Config{
		Addr:           valueOr(getenv, "MAULID_ADDR", ":8080"),
		DBPath:         valueOr(getenv, "MAULID_DB_PATH", "./data/maulid.db"),
		Environment:    valueOr(getenv, "MAULID_ENV", EnvironmentProduction),
		SessionTTL:     7 * 24 * time.Hour,
		RateLimitMax:   12,
		RateLimitFor:   5 * time.Minute,
		TrustedOrigins: splitList(getenv("MAULID_TRUSTED_ORIGINS")),
	}

	if raw := strings.TrimSpace(getenv("MAULID_SESSION_TTL")); raw != "" {
		ttl, err := time.ParseDuration(raw)
		if err != nil || ttl < 5*time.Minute || ttl > 30*24*time.Hour {
			return Config{}, fmt.Errorf("MAULID_SESSION_TTL must be a duration between 5m and 720h")
		}
		cfg.SessionTTL = ttl
	}
	if raw := strings.TrimSpace(getenv("MAULID_RATE_LIMIT_MAX")); raw != "" {
		max, err := strconv.Atoi(raw)
		if err != nil || max < 1 || max > 1000 {
			return Config{}, fmt.Errorf("MAULID_RATE_LIMIT_MAX must be an integer between 1 and 1000")
		}
		cfg.RateLimitMax = max
	}
	if raw := strings.TrimSpace(getenv("MAULID_RATE_LIMIT_WINDOW")); raw != "" {
		window, err := time.ParseDuration(raw)
		if err != nil || window < time.Second || window > time.Hour {
			return Config{}, fmt.Errorf("MAULID_RATE_LIMIT_WINDOW must be a duration between 1s and 1h")
		}
		cfg.RateLimitFor = window
	}
	if cfg.Environment != EnvironmentDevelopment && cfg.Environment != EnvironmentProduction {
		return Config{}, errors.New("MAULID_ENV must be development or production")
	}
	if strings.TrimSpace(cfg.Addr) == "" {
		return Config{}, errors.New("MAULID_ADDR must not be empty")
	}
	if strings.TrimSpace(cfg.DBPath) == "" {
		return Config{}, errors.New("MAULID_DB_PATH must not be empty")
	}

	// Secure cookies are the production default. Development is the only
	// environment that deliberately allows plain-HTTP localhost cookies.
	cfg.CookieSecure = cfg.Environment != EnvironmentDevelopment
	if raw := strings.TrimSpace(getenv("MAULID_COOKIE_SECURE")); raw != "" {
		secure, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, errors.New("MAULID_COOKIE_SECURE must be true or false")
		}
		if cfg.Environment == EnvironmentProduction && !secure {
			return Config{}, errors.New("MAULID_COOKIE_SECURE cannot be disabled in production")
		}
		cfg.CookieSecure = secure
	}

	return cfg, nil
}

func valueOr(getenv func(string) string, key, fallback string) string {
	if value := strings.TrimSpace(getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, strings.TrimSuffix(trimmed, "/"))
		}
	}
	return out
}
