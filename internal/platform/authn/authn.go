// Package authn holds the transport-neutral authenticated principal and the
// middleware that loads it from trusted persistence on every request.
package authn

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"maulid.id/internal/platform/httpx"
)

// Kind identifies how the request proved its identity.
type Kind string

const (
	KindCookie Kind = "cookie"
	KindBearer Kind = "bearer"
)

// Principal is a request-scoped view of the authenticated account. It is
// refreshed from the database on each request; callers must not treat it as a
// long-lived claim.
type Principal struct {
	AccountID         string
	AccountCreatedAt  time.Time
	Email             string
	Name              string
	Role              string
	ValidationStatus  string
	SessionID         string
	SessionCSRFHash   string
	SessionCSRFToken  string
	SessionExpiresAt  time.Time
	SessionCredential string
	Kind              Kind
}

// IsStaff reports whether the principal holds the staff role.
func (p Principal) IsStaff() bool { return p.Role == RoleStaff }

// RoleStaff mirrors the canonical account role. It is duplicated here only to
// keep this transport package free of feature imports.
const RoleStaff = "staff"

// Authenticator resolves an opaque session credential. Implementations must
// return ErrInvalidCredential for unknown, expired or revoked sessions.
type Authenticator interface {
	Authenticate(ctx context.Context, credential string) (Principal, error)
}

// ErrInvalidCredential is returned for any unusable credential.
var ErrInvalidCredential = errors.New("invalid credential")

// CookieName is the browser session cookie.
const CookieName = "maulid_session"

type contextKey struct{}

// WithPrincipal stores a principal in a request context.
func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, contextKey{}, principal)
}

// FromContext reads the principal established by the middleware.
func FromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(contextKey{}).(Principal)
	return principal, ok
}

// Require rejects unauthenticated requests and enforces the CSRF header for
// cookie-authenticated unsafe methods. Bearer clients are not vulnerable to
// CSRF because their credential is explicit.
func Require(authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := resolve(r, authenticator)
			if !ok {
				httpx.Error(w, http.StatusUnauthorized, "unauthenticated", "Sesi tidak valid atau sudah berakhir. Silakan masuk kembali.")
				return
			}
			if !safeMethod(r.Method) && principal.Kind == KindCookie {
				if !validCSRF(principal, r.Header.Get("X-CSRF-Token")) {
					httpx.Error(w, http.StatusForbidden, "csrf_invalid", "Token keamanan tidak valid. Muat ulang halaman lalu coba lagi.")
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), principal)))
		})
	}
}

// Optional attaches a principal when a valid credential is present but never
// rejects or writes a response. An unusable credential is treated as anonymous.
func Optional(authenticator Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if principal, ok := resolve(r, authenticator); ok {
				r = r.WithContext(WithPrincipal(r.Context(), principal))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// resolve silently maps a request credential to a principal. It never writes a
// response, which keeps Optional transparent for anonymous requests.
func resolve(r *http.Request, authenticator Authenticator) (Principal, bool) {
	credential, kind := credentialFromRequest(r)
	if credential == "" {
		return Principal{}, false
	}
	principal, err := authenticator.Authenticate(r.Context(), credential)
	if err != nil {
		return Principal{}, false
	}
	principal.Kind = kind
	principal.SessionCredential = credential
	return principal, true
}

func credentialFromRequest(r *http.Request) (string, Kind) {
	authz := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(authz) > 7 && strings.EqualFold(authz[:7], "Bearer ") {
		return strings.TrimSpace(authz[7:]), KindBearer
	}
	if cookie, err := r.Cookie(CookieName); err == nil && cookie.Value != "" {
		return cookie.Value, KindCookie
	}
	return "", KindCookie
}

func safeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func validCSRF(principal Principal, token string) bool {
	if token == "" || principal.SessionCSRFHash == "" {
		return false
	}
	sum := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(principal.SessionCSRFHash)) == 1
}

// HashToken returns the hex SHA-256 used as the lookup key for session tokens.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
