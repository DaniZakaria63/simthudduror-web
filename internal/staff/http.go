package staff

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/httpx"
)

// Handler exposes the protected backoffice API.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs the staff HTTP adapter.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes mounts staff-only endpoints. The require middleware proves a
// session; the service independently enforces the staff role.
func (h *Handler) RegisterRoutes(mux *http.ServeMux, require func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/staff/accounts", require(http.HandlerFunc(h.list)))
	mux.Handle("PATCH /api/v1/staff/accounts/{id}/validation", require(http.HandlerFunc(h.setValidation)))
	mux.Handle("GET /api/v1/staff/audit", require(http.HandlerFunc(h.audit)))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	principal, _ := authn.FromContext(r.Context())
	filter := ListFilter{
		ValidationStatus: account.ValidationStatus(strings.TrimSpace(r.URL.Query().Get("validation_status"))),
		Query:            strings.TrimSpace(r.URL.Query().Get("q")),
	}
	accounts, err := h.service.List(r.Context(), principal.Role, filter)
	if err != nil {
		h.writeError(w, "list accounts", err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"accounts": accounts})
}

type validationRequest struct {
	ValidationStatus string `json:"validation_status"`
}

func (h *Handler) setValidation(w http.ResponseWriter, r *http.Request) {
	principal, _ := authn.FromContext(r.Context())
	var req validationRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.HandleProblem(w, err)
		return
	}
	updated, audit, err := h.service.SetValidation(r.Context(), principal.Role, principal.AccountID,
		strings.TrimSpace(r.PathValue("id")), req.ValidationStatus)
	if err != nil {
		h.writeError(w, "set validation", err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"account": updated, "audit": audit})
}

func (h *Handler) audit(w http.ResponseWriter, r *http.Request) {
	principal, _ := authn.FromContext(r.Context())
	limit := 0
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 100 {
			httpx.Error(w, http.StatusBadRequest, "invalid_limit", "Batas audit harus antara 1 dan 100.")
			return
		}
		limit = parsed
	}
	entries, err := h.service.Audit(r.Context(), principal.Role, limit)
	if err != nil {
		h.writeError(w, "list audit", err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"audit": entries})
}

func (h *Handler) writeError(w http.ResponseWriter, operation string, err error) {
	var input *InputError
	switch {
	case errors.As(err, &input):
		httpx.Error(w, http.StatusBadRequest, "invalid_input", input.Message)
	case errors.Is(err, ErrForbidden):
		httpx.Error(w, http.StatusForbidden, "staff_required", "Halaman ini hanya untuk staf.")
	case errors.Is(err, ErrAccountNotFound):
		httpx.Error(w, http.StatusNotFound, "account_not_found", "Akun tidak ditemukan.")
	default:
		h.logger.Error(operation+" failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
	}
}
