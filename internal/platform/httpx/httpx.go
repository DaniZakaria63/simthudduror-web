// Package httpx contains the small HTTP plumbing shared by feature adapters.
package httpx

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MaxBodyBytes bounds every JSON request body accepted by the API.
const MaxBodyBytes = 32 << 10

// Problem is a safe, client-facing API error.
type Problem struct {
	Status  int
	Code    string
	Message string
}

func (p *Problem) Error() string { return p.Code }

// JSON writes a JSON response with no-store caching.
func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

// Error writes the standard error envelope.
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

// DecodeJSON strictly decodes one JSON object: it requires a JSON content
// type, limits the body, rejects unknown fields and trailing data.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	contentType := r.Header.Get("Content-Type")
	if mediaType := strings.TrimSpace(strings.Split(contentType, ";")[0]); mediaType != "application/json" {
		return &Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Permintaan harus berformat JSON."}
	}

	body := http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	defer body.Close()

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return decodeProblem(err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return &Problem{Status: http.StatusBadRequest, Code: "invalid_json", Message: "Data permintaan tidak valid."}
	}
	return nil
}

func decodeProblem(err error) error {
	var maxBytes *http.MaxBytesError
	if errors.As(err, &maxBytes) {
		return &Problem{Status: http.StatusRequestEntityTooLarge, Code: "body_too_large", Message: "Data permintaan terlalu besar."}
	}
	return &Problem{Status: http.StatusBadRequest, Code: "invalid_json", Message: "Data permintaan tidak valid."}
}

// HandleProblem writes an error response for a Problem or an internal error.
func HandleProblem(w http.ResponseWriter, err error) {
	var problem *Problem
	if errors.As(err, &problem) {
		Error(w, problem.Status, problem.Code, problem.Message)
		return
	}
	Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
}

// SameOrigin rejects cross-site browser requests for cookie-authenticated
// writes. Requests without an Origin header (mobile clients, curl) are allowed
// because the credential is not ambient. Scheme and host are both checked; in a
// secure deployment only an https origin is accepted.
func SameOrigin(r *http.Request, trusted []string, secure bool) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return false
	}
	if secure && parsed.Scheme != "https" {
		return false
	}
	for _, allowed := range trusted {
		if strings.EqualFold(strings.TrimSuffix(allowed, "/"), strings.TrimSuffix(origin, "/")) {
			return true
		}
	}
	if !strings.EqualFold(parsed.Host, r.Host) {
		return false
	}
	if secure {
		return parsed.Scheme == "https"
	}
	return true
}

// ClientIP uses the direct peer address. X-Forwarded-For is deliberately not
// trusted by default.
func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RequestID derives a short non-secret request identifier for logs.
func RequestID(r *http.Request) string {
	sum := sha256.Sum256([]byte(r.RemoteAddr + "|" + r.Method + "|" + r.URL.Path + "|" + time.Now().UTC().Format(time.RFC3339Nano)))
	return hex.EncodeToString(sum[:6])
}

// SecurityHeaders applies the response hardening used by every route.
func SecurityHeaders(secureTransport bool) func(http.Handler) http.Handler {
	csp := "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			if secureTransport {
				h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/healthz") {
				h.Set("Cache-Control", "no-store")
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Recoverer turns a panic into a generic 500 without exposing internals.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("request panic", "request_id", RequestID(r), "path", r.URL.Path)
					Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan pada server.")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// AccessLog records method, path, status and duration. Bodies, headers and
// credentials are never logged.
func AccessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)
			logger.Info("request",
				"request_id", RequestID(r),
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
