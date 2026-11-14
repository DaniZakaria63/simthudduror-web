package post

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"maulid.id/internal/platform/authn"
	"maulid.id/internal/platform/httpx"
)

// Handler exposes the post feature over HTTP.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs the post HTTP adapter.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// RegisterRoutes mounts the public feed and authenticated write endpoints.
func (h *Handler) RegisterRoutes(mux *http.ServeMux, require, optional func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/posts", optional(http.HandlerFunc(h.list)))
	mux.Handle("POST /api/v1/posts", require(http.HandlerFunc(h.create)))
	mux.Handle("PUT /api/v1/posts/{id}/like", require(http.HandlerFunc(h.like)))
	mux.Handle("DELETE /api/v1/posts/{id}/like", require(http.HandlerFunc(h.unlike)))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	viewerID := ""
	if principal, ok := authn.FromContext(r.Context()); ok {
		viewerID = principal.AccountID
	}
	limit := 0
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 100 {
			httpx.Error(w, http.StatusBadRequest, "invalid_limit", "Batas kabar harus antara 1 dan 100.")
			return
		}
		limit = parsed
	}
	posts, err := h.service.List(r.Context(), r.URL.Query().Get("q"), limit, viewerID)
	if err != nil {
		h.writeError(w, "list", err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"posts": posts})
}

type createRequest struct {
	Body string `json:"body"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	principal, _ := authn.FromContext(r.Context())
	var req createRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.HandleProblem(w, err)
		return
	}
	created, err := h.service.Create(r.Context(), principal.AccountID, req.Body)
	if err != nil {
		h.writeError(w, "create", err)
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]any{"post": created})
}

func (h *Handler) like(w http.ResponseWriter, r *http.Request) {
	h.interact(w, r, true)
}

func (h *Handler) unlike(w http.ResponseWriter, r *http.Request) {
	h.interact(w, r, false)
}

func (h *Handler) interact(w http.ResponseWriter, r *http.Request, like bool) {
	principal, _ := authn.FromContext(r.Context())
	postID := strings.TrimSpace(r.PathValue("id"))
	var (
		updated Post
		err     error
	)
	if like {
		updated, err = h.service.Like(r.Context(), principal.AccountID, postID)
	} else {
		updated, err = h.service.Unlike(r.Context(), principal.AccountID, postID)
	}
	if err != nil {
		h.writeError(w, "interact", err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"post": updated})
}

func (h *Handler) writeError(w http.ResponseWriter, operation string, err error) {
	var input *InputError
	switch {
	case errors.As(err, &input):
		httpx.Error(w, http.StatusBadRequest, "invalid_input", input.Message)
	case errors.Is(err, ErrUnauthenticated):
		httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "Masuk terlebih dahulu untuk melanjutkan.")
	case errors.Is(err, ErrNotValidated):
		httpx.Error(w, http.StatusForbidden, "validation_required", "Akun perlu divalidasi staf sebelum dapat menerbitkan kabar.")
	case errors.Is(err, ErrPostNotFound):
		httpx.Error(w, http.StatusNotFound, "post_not_found", "Kabar tidak ditemukan.")
	default:
		h.logger.Error(operation+" failed", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
	}
}
