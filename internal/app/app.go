// Package app composes the single maulid.id service: SQLite persistence, the
// feature HTTP adapters, health checks and the embedded React build.
package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/config"
	"maulid.id/internal/platform/database"
	"maulid.id/internal/platform/httpx"
	"maulid.id/internal/platform/ratelimit"
	"maulid.id/internal/post"
	"maulid.id/internal/staff"
	"maulid.id/internal/webui"
)

// App owns the composed service and its durable store.
type App struct {
	Handler  http.Handler
	DB       *sql.DB
	Accounts *account.Service
	Posts    *post.Service
	Staff    *staff.Service
}

type validationReader struct {
	repo account.Repository
}

func (v validationReader) ValidationStatus(ctx context.Context, accountID string) (string, error) {
	status, err := v.repo.ValidationStatus(ctx, accountID)
	return string(status), err
}

// New opens the configured database and builds the complete HTTP application.
func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	if logger == nil {
		logger = slog.Default()
	}
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo, cfg.SessionTTL)
	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo, validationReader{repo: accountRepo})
	staffService := staff.NewService(staff.NewRepository(db))

	require := authn.Require(accountService)
	optional := authn.Optional(accountService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	account.NewHandler(accountService, ratelimit.New(cfg.RateLimitMax, cfg.RateLimitFor), logger, cfg.CookieSecure, cfg.TrustedOrigins).
		RegisterRoutes(mux, require, optional)
	post.NewHandler(postService, logger).RegisterRoutes(mux, require, optional)
	staff.NewHandler(staffService, logger).RegisterRoutes(mux, require)

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, "not_found", "Endpoint tidak ditemukan.")
	})
	mux.Handle("/", webui.Handler(webui.Embedded()))

	handler := httpx.Recoverer(logger)(
		httpx.AccessLog(logger)(
			httpx.SecurityHeaders(cfg.CookieSecure)(mux),
		),
	)
	return &App{
		Handler: handler, DB: db,
		Accounts: accountService, Posts: postService, Staff: staffService,
	}, nil
}

// Close releases the database handle.
func (a *App) Close() error {
	if a == nil || a.DB == nil {
		return nil
	}
	return a.DB.Close()
}
