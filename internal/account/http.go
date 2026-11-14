package account

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/httpx"
	"maulid.id/internal/platform/ratelimit"
)

// Handler exposes the account feature over HTTP.
type Handler struct {
	service        *Service
	limiter        *ratelimit.Limiter
	logger         *slog.Logger
	cookieSecure   bool
	trustedOrigins []string
}

// NewHandler constructs the account HTTP adapter.
func NewHandler(service *Service, limiter *ratelimit.Limiter, logger *slog.Logger, cookieSecure bool, trustedOrigins []string) *Handler {
	return &Handler{
		service: service, limiter: limiter, logger: logger,
		cookieSecure: cookieSecure, trustedOrigins: trustedOrigins,
	}
}

// RegisterRoutes mounts public and session-protected account endpoints.
func (h *Handler) RegisterRoutes(mux *http.ServeMux, require, optional func(http.Handler) http.Handler) {
	mux.HandleFunc("POST /api/v1/auth/register", h.register)
	mux.HandleFunc("POST /api/v1/auth/login", h.login)
	mux.Handle("POST /api/v1/auth/logout", require(http.HandlerFunc(h.logout)))
	mux.Handle("GET /api/v1/auth/session", optional(http.HandlerFunc(h.session)))
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type sessionResponse struct {
	Account   View      `json:"account"`
	CSRFToken string    `json:"csrf_token,omitempty"`
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if !h.guard(w, r, "register") {
		return
	}
	var req registerRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.HandleProblem(w, err)
		return
	}
	result, err := h.service.Register(r.Context(), req.Email, req.Name, req.Password)
	if err != nil {
		h.writeServiceError(w, "register", err)
		return
	}
	h.writeSession(w, r, http.StatusCreated, result)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if !h.guard(w, r, "login") {
		return
	}
	var req loginRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.HandleProblem(w, err)
		return
	}
	result, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.writeServiceError(w, "login", err)
		return
	}
	h.limiter.Reset("login:" + httpx.ClientIP(r))
	h.writeSession(w, r, http.StatusOK, result)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	principal, ok := authn.FromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "Masuk terlebih dahulu untuk melanjutkan.")
		return
	}
	if err := h.service.Logout(r.Context(), principal.SessionCredential); err != nil {
		h.logger.Error("logout failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
		return
	}
	h.clearCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) session(w http.ResponseWriter, r *http.Request) {
	principal, ok := authn.FromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "Belum ada sesi aktif.")
		return
	}
	httpx.JSON(w, http.StatusOK, sessionResponse{
		Account: View{
			ID: principal.AccountID, Email: principal.Email, Name: principal.Name,
			Role: Role(principal.Role), ValidationStatus: ValidationStatus(principal.ValidationStatus),
			CreatedAt: principal.AccountCreatedAt,
		},
		CSRFToken: principal.SessionCSRFToken,
		ExpiresAt: principal.SessionExpiresAt,
	})
}

func (h *Handler) guard(w http.ResponseWriter, r *http.Request, scope string) bool {
	if !httpx.SameOrigin(r, h.trustedOrigins, h.cookieSecure) {
		httpx.Error(w, http.StatusForbidden, "origin_invalid", "Permintaan berasal dari origin yang tidak diizinkan.")
		return false
	}
	allowed, retry := h.limiter.Allow(scope + ":" + httpx.ClientIP(r))
	if !allowed {
		w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
		httpx.Error(w, http.StatusTooManyRequests, "rate_limited", "Terlalu banyak percobaan. Tunggu sebentar lalu coba lagi.")
		return false
	}
	return true
}

func (h *Handler) writeSession(w http.ResponseWriter, r *http.Request, status int, result AuthResult) {
	http.SetCookie(w, &http.Cookie{
		Name:     authn.CookieName,
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(result.ExpiresAt).Seconds()),
	})
	response := sessionResponse{
		Account:   result.Account.ToView(),
		CSRFToken: result.CSRFToken,
		ExpiresAt: result.ExpiresAt,
	}
	if wantsMobileToken(r) {
		response.Token = result.Token
	}
	httpx.JSON(w, status, response)
}

func (h *Handler) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     authn.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (h *Handler) writeServiceError(w http.ResponseWriter, operation string, err error) {
	var input *InputError
	switch {
	case errors.As(err, &input):
		httpx.Error(w, http.StatusBadRequest, "invalid_input", input.Message)
	case errors.Is(err, ErrEmailTaken):
		httpx.Error(w, http.StatusConflict, "email_unavailable", "Email tidak dapat digunakan.")
	case errors.Is(err, ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, "invalid_credentials", "Email atau kata sandi salah.")
	default:
		h.logger.Error(operation+" failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
	}
}

func wantsMobileToken(r *http.Request) bool {
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Maulid-Client")), "mobile")
}
